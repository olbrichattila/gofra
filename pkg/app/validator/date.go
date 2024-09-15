package validator

import (
	"fmt"
	"time"
)

func DateRule(val string, _ ...string) (string, bool) {
	layout := "2006-01-02"
	_, err := time.Parse(layout, val)
	if err == nil {
		return "", true
	}

	return fmt.Sprintf("%s not in an ISO date YYYY-MM-DD", val), false
}

func DateTimeRule(val string, _ ...string) (string, bool) {
	layout := "2006-01-02 15:04:05"
	_, err := time.Parse(layout, val)
	if err == nil {
		return "", true
	}

	return fmt.Sprintf("%s not in an ISO date time YYYY-MM-DD HH:MM:SS", val), false
}
