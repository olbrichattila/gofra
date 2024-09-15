package request

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// TODO add JsonBody, passing a struct, and receiving back unmarshalled
type Requester interface {
	SetRequest(*http.Request)
	SetRouteParameters(map[string]string)
	GetRequest() *http.Request
	Get() map[string][]string
	Post() map[string][]string
	GetOne(string, string) string
	PostOne(string, string) string
	All() map[string][]string
	AllFlat() map[string]string
	AllOne(string, string) string
	Body() string
	JSONBody() map[string]interface{}
	JSONToStruct(result interface{})
	Headers() map[string][]string
}

func New() Requester {
	return &Request{}
}

type Request struct {
	r           *http.Request
	routeParams map[string]string
}

func (r *Request) SetRequest(req *http.Request) {
	r.r = req
}

func (r *Request) SetRouteParameters(routeParams map[string]string) {
	r.routeParams = routeParams
}

func (r *Request) GetRouteParameters() map[string]string {
	if r.routeParams == nil {
		return make(map[string]string, 0)
	}
	return r.routeParams
}

func (r *Request) GetRequest() *http.Request {
	return r.r
}

func (r *Request) Get() map[string][]string {
	return r.r.URL.Query()
}

func (r *Request) Post() map[string][]string {
	r.r.ParseForm() // TODO should handle error
	return r.r.Form
}

func (r *Request) GetOne(par, def string) string {
	v := r.r.URL.Query().Get(par)
	if v == "" {
		rp := r.GetRouteParameters()
		if par, ok := rp[par]; ok {
			return par
		}

		return def
	}

	return v
}

func (r *Request) PostOne(par, def string) string {
	err := r.r.ParseForm()
	if err != nil {
		return def
	}

	paramValue := r.r.FormValue(par)
	if paramValue == "" {
		return def
	}

	return paramValue
}

func (r *Request) All() map[string][]string {
	get := r.Get()
	post := r.Post()

	var res = make(map[string][]string, 0)
	for k, v := range get {
		res[k] = v
	}

	for k, v := range post {
		res[k] = v
	}

	return res
}

func (r *Request) AllFlat() map[string]string {
	result := make(map[string]string)
	all := r.All()
	for key, values := range all {
		if len(values) == 1 {
			result[key] = values[0]
			continue
		}
		for i, value := range values {
			result[key+"["+strconv.Itoa(i)+"]"] = value
		}
	}

	return result
}

func (r *Request) AllOne(par, def string) string {
	getV := r.GetOne(par, "")
	if getV != "" {
		return getV
	}

	return r.PostOne(par, def)
}

func (r *Request) Headers() map[string][]string {
	return r.r.Header
}

func (r *Request) Body() string {
	body, _ := io.ReadAll(r.r.Body)

	return string(body)
}

func (r *Request) JSONBody() map[string]interface{} {
	var result map[string]interface{}
	body, _ := io.ReadAll(r.r.Body)
	json.Unmarshal(body, &result)

	return result
}

func (r *Request) JSONToStruct(result interface{}) {
	body, _ := io.ReadAll(r.r.Body)
	json.Unmarshal(body, &result)
}
