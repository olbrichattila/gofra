package validator

import (
	"strings"
)

const (
	pipePlaceholder = "__PIPE__PLACEHOLDER__"
)

type CustomValidatorFunc func(map[string]string) (ValidationErrors, bool)

type ValidationRule struct {
	Redirect   string
	Rules      map[string]string
	CustomRule CustomValidatorFunc
}

type RuleFunc func(string, ...string) (string, bool)
type ValidationErrors map[string][]string

func New() Validator {
	return &Validate{}
}

type Validator interface {
	Validate(map[string]string, map[string]string) (bool, ValidationErrors, map[string]string)
	SetRules(map[string]RuleFunc)
}

type Validate struct {
	rules map[string]RuleFunc
}

func (v *Validate) SetRules(rules map[string]RuleFunc) {
	if v.rules == nil {
		v.rules = make(map[string]RuleFunc)
	}

	for ruleName, fn := range rules {
		v.rules[ruleName] = fn
	}
}

func (v *Validate) Validate(fields map[string]string, rules map[string]string) (bool, ValidationErrors, map[string]string) {
	valid := make(map[string]string)
	validationErrors := make(ValidationErrors)
	for field, rule := range rules {
		fieldValue := fields[field]
		if messages, ok := v.parse(fieldValue, rule); ok {
			valid[field] = fieldValue
		} else {
			validationErrors[field] = messages
		}
	}

	return len(validationErrors) == 0, validationErrors, valid
}

func (v *Validate) parse(val, pattern string) ([]string, bool) {
	errorMessages := make([]string, 0)
	pattern = strings.ReplaceAll(pattern, `\|`, pipePlaceholder)
	patterns := strings.Split(pattern, "|")
	for _, rule := range patterns {
		rule = strings.ReplaceAll(rule, pipePlaceholder, "|")
		if message, ok := v.parseRule(val, rule); !ok {
			errorMessages = append(errorMessages, message)
		}
	}
	return errorMessages, len(errorMessages) == 0
}

func (v *Validate) parseRule(val, rule string) (string, bool) {
	rulePars := make([]string, 0)
	ruleParts := strings.Split(rule, ":")
	ruleName := ruleParts[0]
	if len(ruleParts) > 1 {
		rulePars = strings.Split(ruleParts[1], ",")
	}

	if ruleFn, ok := v.rules[ruleName]; ok {
		return ruleFn(val, rulePars...)
	}

	return "", true
}
