package validator

import (
	"fmt"
	"net/url"
)

func UrlRule(val string, _ ...string) (string, bool) {
	_, err := url.Parse(val)
	if err == nil {
		return "", true
	}

	return fmt.Sprintf("%s not in an URL", val), false
}
