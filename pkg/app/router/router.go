package router

import (
	"fmt"
	"net/url"
	"strings"
)

type ControllerAction struct {
	Path            string
	RequestType     []string
	Fn              interface{}
	Middlewares     []any
	ViewAutoLoads   []string
	ValidationRules string
	Controller      func() any
	ActionName      string
	IsStatic        bool
	StaticPath      string
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

func (r *Route) Match(route, requestUrl string) (bool, map[string]string) {
	params := make(map[string]string)
	routePars := strings.Split(route, "/")
	baseUrl := strings.Split(requestUrl, "?")[0]
	urlPars := strings.Split(baseUrl, "/")
	if len(routePars) != len(urlPars) && routePars[len(routePars)-1] != "**" {
		return false, nil
	}

	for i, part := range routePars {
		// Bind parameters
		if len(part) > 0 && part[0] == ':' {
			par, err := url.QueryUnescape(urlPars[i])
			if err != nil {
				par = urlPars[i]
			}
			params[part[1:]] = par
			continue
		}

		// Bind single star, one route par, use only once in the route
		if part == "*" {
			par, err := url.QueryUnescape(urlPars[i])
			if err != nil {
				par = urlPars[i]
			}
			params["*"] = par
			continue
		}

		// Bind double star, resolves all other part of the route
		if part == "**" {
			params["*"] = r.getRestOfUrl(urlPars, i)
			return true, params
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

func (*Route) getRestOfUrl(urlPars []string, index int) string {
	sb := &strings.Builder{}
	for i := index; i < len(urlPars); i++ {
		if i > index {
			sb.WriteRune('/')
		}
		unescaped, err := url.QueryUnescape(urlPars[i])
		if err != nil {
			unescaped = urlPars[i]
		}
		sb.WriteString(unescaped)
	}

	return sb.String()
}
