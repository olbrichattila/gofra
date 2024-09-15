package validator

import (
	"fmt"
	"strconv"
)

func BetweenRule(val string, pars ...string) (string, bool) {
	if len(pars) != 2 {
		return "between, requires 2 numeric parameter like between:5,10", false
	}

	a, err := strconv.Atoi(val)
	if err != nil {
		return "The first comparing must be a number", false
	}

	b, err := strconv.Atoi(pars[0])
	if err != nil {
		return "the first parameter must be a number", false
	}

	c, err := strconv.Atoi(pars[1])
	if err != nil {
		return "the second parameter must be a number", false
	}

	if a >= b && a <= c {
		return "", true
	}

	return fmt.Sprintf("%s is not between %s and %s", val, pars[0], pars[1]), false
}
