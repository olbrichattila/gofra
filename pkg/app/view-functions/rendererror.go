package internalviewfunction

import (
	"encoding/json"
	"strings"

	"github.com/olbrichattila/gofra/pkg/app/validator"
)

func RenderErrors(pars ...any) string {
	if len(pars) != 1 {
		return ""
	}

	jsonStr, ok := pars[0].(string)
	if !ok {
		return ""
	}

	if jsonStr == "" {
		return ""
	}

	var validatorErrors validator.ValidationErrors
	err := json.Unmarshal([]byte(jsonStr), &validatorErrors)
	if err != nil {
		return err.Error()
	}

	sw := &strings.Builder{}
	sw.WriteString("<ul>")
	for field, errorList := range validatorErrors {
		sw.WriteString("<li>")
		sw.WriteString(field)
		sw.WriteString("<ul>")
		for _, eMess := range errorList {
			sw.WriteString("<li>")
			sw.WriteString(eMess)
			sw.WriteString("</li>")
		}
		sw.WriteString("</ul>")
		sw.WriteString("</li>")
	}
	sw.WriteString("</ul>")

	return sw.String()
}

func RenderError(pars ...any) string {
	if len(pars) != 2 {
		return ""
	}

	field, ok := pars[0].(string)
	if !ok {
		return ""
	}

	jsonStr, ok := pars[1].(string)
	if !ok {
		return ""
	}

	if jsonStr == "" {
		return ""
	}

	var validatorErrors validator.ValidationErrors

	err := json.Unmarshal([]byte(jsonStr), &validatorErrors)
	if err != nil {
		return err.Error()
	}

	if errorList, ok := validatorErrors[field]; ok {
		return strings.Join(errorList, ", ")
	}

	return ""
}
