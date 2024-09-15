package internalcommand

import (
	"fmt"

	"github.com/olbrichattila/gofra/pkg/app/args"
	wizard "github.com/olbrichattila/gofra/pkg/app/wizards/class"
	commandcreator "github.com/olbrichattila/gofra/pkg/app/wizards/command"
)

func CreateJob(a args.CommandArger, c commandcreator.CommandCreator, cc wizard.ClassCreator) {
	cc.SetHelpHeader("Usage:\n  go run ./cmd/ artisan create:job <job-name> <optional-parameters>\n")
	cc.SetParameterInfos(getDefaultCreateDiMapping())
	cc.SetTemplates(map[string]string{
		"": "job.tpl",
	})
	cc.SetOutParameterInfos(map[string]wizard.ParameterInfo{})

	if _, ok := a.GetFlagByName("help", ""); ok {
		fmt.Println(cc.GetHelp())
		return
	}

	flags := a.GetAllFlags()
	template := cc.GetTemplate(flags)

	templateParams := cc.GetTemplateParams(flags)
	err := c.Create(template, "./app/jobs", templateParams)
	if err != nil {
		fmt.Printf("%s\nTry -help\n", err.Error())
		return
	}

	fmt.Printf("Please register your new job in:\n  app/config/job.go\n")
}
