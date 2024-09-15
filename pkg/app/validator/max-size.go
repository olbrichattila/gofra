package validator

import (
	"fmt"
	"strconv"
)

func MaxSizeRule(val string, pars ...string) (string, bool) {
	if len(pars) != 1 {
		return "requires 1 numeric parameter like size:1", false
	}

	a, err := strconv.Atoi(pars[0])
	if err != nil {
		return "The value comparing must be a number", false
	}

	if len(val) <= a {
		return "", true
	}

	return fmt.Sprintf("must be shorter then %d", a), false
}
