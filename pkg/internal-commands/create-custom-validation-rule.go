package internalcommand

import (
	"fmt"

	"github.com/olbrichattila/gofra/pkg/app/args"
	commandcreator "github.com/olbrichattila/gofra/pkg/app/wizards/command"
)

func CreateCustomValidationRule(a args.CommandArger, c commandcreator.CommandCreator) {
	err := c.Create("customrule.tpl", "./app/validator-configs", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Please register your new rule in in:\n app/config/validators.go\n")
}
