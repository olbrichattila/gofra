package internalcommand

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/olbrichattila/gofra/pkg/app/config"
)

func ListJobs(c config.Configer) {
	jobs := c.Jobs()
	reorderSlice := make([]string, len(jobs))

	for i, job := range jobs {
		reorderSlice[i] = runtime.FuncForPC(reflect.ValueOf(job.Fn).Pointer()).Name()
	}

	for _, jobName := range reorderSlice {
		fmt.Println(jobName)
	}
}
