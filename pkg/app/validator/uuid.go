package validator

import (
	"fmt"
	"regexp"
)

func UuidRule(val string, _ ...string) (string, bool) {
	re, err := regexp.Compile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)
	if err != nil {
		return "regex error validating email", false
	}
	if re.MatchString(val) {
		return "", true
	}

	return fmt.Sprintf("%s not in an email", val), false
}
