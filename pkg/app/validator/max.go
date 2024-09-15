package validator

import (
	"fmt"
	"strconv"
)

func MaxRule(val string, pars ...string) (string, bool) {
	if len(pars) != 1 {
		return "min, requires 1 numeric parameter like min:1", false
	}

	a, err := strconv.Atoi(val)
	if err != nil {
		return "The value comparing must be a number", false
	}

	b, err := strconv.Atoi(pars[0])
	if err != nil {
		return "min parameter must be a number min:1", false
	}

	if a >= b {
		return fmt.Sprintf("%s must be smaller then %s", val, pars[0]), false
	}

	return "", true
}
