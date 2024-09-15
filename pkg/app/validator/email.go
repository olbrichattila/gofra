package validator

import (
	"fmt"
	"regexp"
)

func EmailRule(val string, _ ...string) (string, bool) {
	re, err := regexp.Compile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if err != nil {
		return "regex error validating email", false
	}
	if re.MatchString(val) {
		return "", true
	}

	return fmt.Sprintf("%s not in an email", val), false
}
