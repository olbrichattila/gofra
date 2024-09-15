package internalcommand

import (
	"fmt"
	"sort"

	"github.com/olbrichattila/gofra/pkg/app/config"
	"github.com/olbrichattila/gofra/pkg/app/router"
)

func ListRoutes(c config.Configer) {
	routes := c.Routes()
	reorderSlice := make([]router.ControllerAction, len(routes))
	copy(reorderSlice, routes)

	sort.Slice(reorderSlice, func(i, j int) bool {
		return reorderSlice[i].Path < reorderSlice[j].Path
	})

	for _, route := range reorderSlice {
		fmt.Println(route.Path)
	}
}
