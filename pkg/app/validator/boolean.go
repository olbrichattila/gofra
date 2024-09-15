package validator

import (
	"fmt"
	"strings"
)

func BooleanRule(val string, _ ...string) (string, bool) {
	c := strings.ToLower(val)

	if c == "0" || c == "1" || c == "true" || c == "false" {
		return "", true
	}

	return fmt.Sprintf("%s not in a boolean 1,0,true,false", val), false
}
