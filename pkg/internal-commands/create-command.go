package internalcommand

import (
	"fmt"

	"github.com/olbrichattila/gofra/pkg/app/args"
	commandcreator "github.com/olbrichattila/gofra/pkg/app/wizards/command"
)

func CreateCommand(a args.CommandArger, c commandcreator.CommandCreator) {
	err := c.Create("command.tpl", "./app/commands", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Please register your new command in:\n  app/config/commands.go\n")
}
