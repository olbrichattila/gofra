package internalcommand

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/olbrichattila/gofra/pkg/app/config"
)

func ListMiddlewares(c config.Configer) {
	c.ViewConfig()
	middlewares := c.Middlewares()
	reorderSlice := make([]string, len(middlewares))

	for i, middleware := range middlewares {
		reorderSlice[i] = runtime.FuncForPC(reflect.ValueOf(middleware).Pointer()).Name()
	}

	for _, middlewareName := range reorderSlice {
		fmt.Println(middlewareName)
	}
}
