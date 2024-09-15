package event

import "github.com/olbrichattila/godi"

type Eventer interface {
	Construct(godi.Container)
	Subscribe(string, string, interface{})
	UnSubscribe(string, string)
	Flush()
	Dispatch(string, string)
}
