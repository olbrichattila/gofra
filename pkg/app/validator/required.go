package validator

import (
	"strings"
)

func RequiredRule(val string, _ ...string) (string, bool) {
	if len(strings.TrimSpace(val)) == 0 {
		return "is required", false
	}

	return "", true

}
