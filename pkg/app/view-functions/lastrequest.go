package internalviewfunction

import (
	"encoding/json"
)

func RenderRequest(pars ...any) string {
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

	var lastRequest map[string]string

	err := json.Unmarshal([]byte(jsonStr), &lastRequest)
	if err != nil {
		return err.Error()
	}

	if requestText, ok := lastRequest[field]; ok {
		return requestText
	}

	return ""
}
