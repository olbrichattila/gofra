package validator

import (
	"fmt"
	"strings"
)

func InRule(val string, pars ...string) (string, bool) {
	for _, elem := range pars {
		if elem == val {
			return "", true
		}
	}

	return fmt.Sprintf("%s not in %s", val, strings.Join(pars, ",")), false
}
