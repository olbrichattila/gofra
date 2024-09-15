package gofra

import (
	"fmt"
	"html/template"
	"os"

	"github.com/olbrichattila/godi"
	app "github.com/olbrichattila/gofra/pkg"
	commandexecutor "github.com/olbrichattila/gofra/pkg/app/command"
	"github.com/olbrichattila/gofra/pkg/app/config"
	"github.com/olbrichattila/gofra/pkg/app/cron"
	"github.com/olbrichattila/gofra/pkg/app/router"
	"github.com/olbrichattila/gofra/pkg/app/validator"
)

func Run(
	bootstrapFunc interface{},
	routes []router.ControllerAction,
	jobs []cron.Job,
	middlewares []interface{},
	appBindings []config.DiCallback,
	consoleCommands map[string]commandexecutor.CommandItem,
	appViewConfig template.FuncMap,
	templateAutoLoad map[string][]string,
	validationRules map[string]validator.ValidationRule,
	validationRuleFuncs map[string]validator.RuleFunc,
) {
	args := os.Args

	app := app.New(
		godi.New(),
		bootstrapFunc,
		routes,
		jobs,
		middlewares,
		appBindings,
		consoleCommands,
		appViewConfig,
		templateAutoLoad,
		validationRules,
		validationRuleFuncs,
	)
	if len(args) < 2 {
		app.Serve()
		return
	}

	switch args[1] {
	case "serve":
		app.Serve()
	case "artisan":
		app.Command()
	default:
		displayHelp()
		return
	}
}

func displayHelp() {
	fmt.Printf(
		`Usage:
Run HTTP server:	
     go run ./cmd serve
Run command:	 
     go run ./cmd artisan <command>
`,
	)
}
