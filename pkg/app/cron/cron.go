package cron

import (
	"time"

	"github.com/olbrichattila/godi"
)

func New() JobTimer {
	return &timer{}
}

type JobTimer interface {
	Init(godi godi.Container, jobs []Job)
}

type Job struct {
	Seconds int
	Fn      interface{}
}

type timer struct {
}

func (t *timer) Init(godi godi.Container, jobs []Job) {
	for _, job := range jobs {
		go func(j Job) {
			for {
				godi.Call(j.Fn)
				time.Sleep(time.Duration(j.Seconds) * time.Second)
			}
		}(job)
	}
}
