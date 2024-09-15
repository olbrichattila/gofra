package internalcommand

import (
	"fmt"

	"github.com/olbrichattila/gofra/pkg/app/args"
	commandcreator "github.com/olbrichattila/gofra/pkg/app/wizards/command"
)

func CreateEventConsumer(a args.CommandArger, c commandcreator.CommandCreator) {
	err := c.Create("eventconsumer.tpl", "./app/events", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(
		`It is possible to subscribe to the events any place in your application, for example in app/bootstrap.go
Example:
	e.Subscribe("topic", "e4", eventconsumer.ExampleConsumer)
	e.Dispatch("topic", "event1")`)
}
