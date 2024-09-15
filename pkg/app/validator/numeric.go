package validator

import (
	"fmt"
	"strconv"
)

func NumericRule(val string, _ ...string) (string, bool) {
	_, err := strconv.ParseFloat(val, 64)
	if err == nil {
		return "", true
	}

	return fmt.Sprintf("%s not in an number", val), false
}
