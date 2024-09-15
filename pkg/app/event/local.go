package event

import (
	"github.com/olbrichattila/godi"
)

func NewLocalEvent() Eventer {
	return &LocalEvent{
		eventMap: make(map[string]map[string]interface{}),
	}
}

type LocalEvent struct {
	di       godi.Container
	eventMap map[string]map[string]interface{}
}

func (l *LocalEvent) Construct(di godi.Container) {
	l.di = di
}

func (l *LocalEvent) Subscribe(topic, name string, event interface{}) {
	_, ok := l.eventMap[topic]
	if !ok {
		l.eventMap[topic] = make(map[string]interface{})
	}

	l.eventMap[topic][name] = event
}

func (l *LocalEvent) UnSubscribe(topic, name string) {
	if mapTopic, ok := l.eventMap[topic]; ok {
		if _, topicOk := mapTopic[name]; topicOk {
			delete(mapTopic, name)

			if len(mapTopic) == 0 {
				delete(l.eventMap, topic)
			}
		}
	}
}

func (l *LocalEvent) Flush() {
	l.eventMap = make(map[string]map[string]interface{})
}

func (l *LocalEvent) Dispatch(topic, payload string) {
	if events, ok := l.eventMap[topic]; ok {
		for _, e := range events {
			go func(e interface{}, payload string) {
				l.di.Call(e, payload)
			}(e, payload)
		}
	}
}
