package internalcommand

import (
	"fmt"

	"github.com/olbrichattila/gofra/pkg/app/args"
	commandcreator "github.com/olbrichattila/gofra/pkg/app/wizards/command"
)

func CreateCustomViewFunction(a args.CommandArger, c commandcreator.CommandCreator) {
	err := c.Create("view-function.tpl", "./app/view-functions", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Please register your new command in:\n  app/config/view.go\n")
}
