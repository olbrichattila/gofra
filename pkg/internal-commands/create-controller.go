package internalcommand

import (
	"fmt"

	"github.com/olbrichattila/gofra/pkg/app/args"
	wizard "github.com/olbrichattila/gofra/pkg/app/wizards/class"
	commandcreator "github.com/olbrichattila/gofra/pkg/app/wizards/command"
)

func CreateController(a args.CommandArger, c commandcreator.CommandCreator, cc wizard.ClassCreator) {
	cc.SetHelpHeader("Usage:\n  go run ./cmd/ artisan create:controller <controller-name> <optional-parameters>\n")
	cc.SetParameterInfos(getDefaultCreateDiMapping())
	cc.SetTemplates(map[string]string{
		"":     "controller-default.tpl",
		"api":  "controller-api.tpl",
		"crud": "controller-crud.tpl",
	})
	cc.SetOutParameterInfos(map[string]wizard.ParameterInfo{
		"string": {Name: "string", Alias: "\"\"", ImportPath: ""},
		"error":  {Name: "error", Alias: "nil", ImportPath: ""},
		"bool":   {Name: "bool", Alias: "false", ImportPath: ""},
	})

	if _, ok := a.GetFlagByName("help", ""); ok {
		fmt.Println(cc.GetHelp())
		return
	}

	flags := a.GetAllFlags()
	template := cc.GetTemplate(flags)

	templateParams := cc.GetTemplateParams(flags)
	err := c.Create(template, "./app/controllers", templateParams)
	if err != nil {
		fmt.Printf("%s\nTry -help\n", err.Error())
		return
	}

	fmt.Printf("Please register your new controller action(s) in:\n  app/config/routes.go\n")
}
