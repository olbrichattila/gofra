package commandexecutor

import (
	"fmt"
	"os"

	"github.com/olbrichattila/godi"
)

const (
	defaultCommandName = "list-commands"
)

func New() CommandExecutor {
	return &Cexecutor{}
}

type CommandItem struct {
	Fn   interface{}
	Desc string
}

type CommandExecutor interface {
	Execute(godi.Container, map[string]CommandItem) error
}

type Cexecutor struct {
}

func (*Cexecutor) Execute(di godi.Container, commands map[string]CommandItem) error {
	args := os.Args
	commandName := defaultCommandName
	if len(args) >= 3 {
		commandName = args[2]
	}

	if commandItem, ok := commands[commandName]; ok {
		_, err := di.Call(commandItem.Fn)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("command not defined")
}
