package internalconfig

import "github.com/olbrichattila/gofra/pkg/app/validator"

var ValidatorRules = map[string]validator.RuleFunc{
	"required": validator.RequiredRule,
	"min":      validator.MinRule,
	"max":      validator.MaxRule,
	"in":       validator.InRule,
	"regex":    validator.RegexRule,
	"between":  validator.BetweenRule,
	"size":     validator.SizeRule,
	"email":    validator.EmailRule,
	"url":      validator.UrlRule,
	"uuid":     validator.UuidRule,
	"numeric":  validator.NumericRule,
	"integer":  validator.IntegerRule,
	"date":     validator.DateRule,
	"dateTime": validator.DateTimeRule,
	"boolean":  validator.BooleanRule,
	"json":     validator.JSONRule,
	"minSize":  validator.MinSizeRule,
	"maxSize":  validator.MaxSizeRule,
}
