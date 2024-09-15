package args

import (
	"fmt"
	"os"
	"strings"
)

func New() CommandArger {
	return &CommandArg{}
}

type CommandArger interface {
	GetAll() []string
	Get(int) (string, error)
	GetAllFlags() map[string]string
	GetFlagByName(string, string) (string, bool)
}

type CommandArg struct {
}

func (*CommandArg) GetAll() []string {

	newArgs := make([]string, 0)
	if len(os.Args) < 3 {
		return newArgs
	}

	commandArgs := os.Args[3:]

	for _, arg := range commandArgs {
		if !strings.HasPrefix(arg, "-") {
			newArgs = append(newArgs, arg)
		}
	}

	return newArgs
}

func (c *CommandArg) Get(index int) (string, error) {
	args := c.GetAll()

	if len(args) <= index || index < 0 {
		return "", fmt.Errorf("parameter index out of bounds")
	}

	return args[index], nil
}

func (c *CommandArg) GetAllFlags() map[string]string {
	commandArgs := os.Args[3:]
	newArgs := make(map[string]string, 0)

	for _, arg := range commandArgs {
		if strings.HasPrefix(arg, "-") {
			parts := strings.Split(arg, "=")
			flagName := strings.TrimPrefix(parts[0], "-")
			if len(parts) == 1 {
				newArgs[flagName] = ""
			} else {
				newArgs[flagName] = parts[1]
			}
		}
	}

	return newArgs
}

func (c *CommandArg) GetFlagByName(key, defaultValue string) (string, bool) {
	allFlags := c.GetAllFlags()
	if flag, ok := allFlags[key]; ok {
		return flag, true
	}

	return defaultValue, false
}
