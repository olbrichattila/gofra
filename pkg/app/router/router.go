package router

import (
	"fmt"
	"net/url"
	"strings"
)

type ControllerAction struct {
	Path            string
	RequestType     string
	Fn              interface{}
	Middlewares     []any
	ViewAutoLoads   []string
	ValidationRules string
}

type Router interface {
	Match(string, string) (bool, map[string]string)
	Build(string, map[string]string) (string, error)
}

type Route struct {
}

func NewRouter() Router {
	return &Route{}
}

func (*Route) Match(route, requestUrl string) (bool, map[string]string) {
	params := make(map[string]string)
	routePars := strings.Split(route, "/")
	baseUrl := strings.Split(requestUrl, "?")[0]
	urlPars := strings.Split(baseUrl, "/")
	if len(routePars) != len(urlPars) {
		return false, nil
	}

	for i, part := range routePars {
		if len(part) > 0 && part[0] == ':' {
			par, err := url.QueryUnescape(urlPars[i])
			if err != nil {
				par = urlPars[i]
			}
			params[part[1:]] = par
			continue
		}

		if urlPars[i] != part {
			return false, nil
		}
	}

	return true, params
}

func (*Route) Build(route string, pars map[string]string) (string, error) {
	sb := &strings.Builder{}
	routeParts := strings.Split(route, "/")
	for i, part := range routeParts {
		if i > 0 {
			sb.WriteRune('/')
		}
		if len(part) > 0 && part[0] == ':' {
			key := part[1:]
			if par, ok := pars[key]; ok {
				sb.WriteString(url.QueryEscape(par))
				continue
			}

			return "", fmt.Errorf("missing parameter for " + key)
		}

		sb.WriteString(part)

	}

	return sb.String(), nil
}
