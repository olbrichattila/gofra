package validator

import (
	"fmt"
	"regexp"
	"strings"
)

func RegexRule(val string, pars ...string) (string, bool) {
	pattern := strings.Join(pars, ",")

	fmt.Println(pattern)

	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Sprintf("Invalid regular expression %s", pattern), false
	}

	if re.MatchString(val) {
		return "", true
	}

	return fmt.Sprintf("the value does not match expression %s", pattern), false
}
