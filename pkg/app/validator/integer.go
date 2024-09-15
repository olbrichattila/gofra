package validator

import (
	"fmt"
	"strconv"
)

func IntegerRule(val string, _ ...string) (string, bool) {
	_, err := strconv.Atoi(val)
	if err == nil {
		return "", true
	}

	return fmt.Sprintf("%s not in an number", val), false
}
