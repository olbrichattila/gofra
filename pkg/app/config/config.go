package config

import (
	"text/template"

	commandexecutor "github.com/olbrichattila/gofra/pkg/app/command"
	"github.com/olbrichattila/gofra/pkg/app/cron"
	"github.com/olbrichattila/gofra/pkg/app/router"

	"github.com/olbrichattila/godi"
)

type DiCallback func(godi.Container) (string, interface{}, error)

func New(
	routes []router.ControllerAction,
	jobs []cron.Job,
	middlewares []interface{},
	appBindings []DiCallback,
	internalBindings []DiCallback,
	appCommands map[string]commandexecutor.CommandItem,
	internalCommands map[string]commandexecutor.CommandItem,
	appViewConfig template.FuncMap,
	internalViewConfig template.FuncMap,
	templateAutoLoad map[string][]string,
) Configer {
	return &Conf{
		routes:             routes,
		jobs:               jobs,
		middlewares:        middlewares,
		appBindings:        appBindings,
		internalBindings:   internalBindings,
		appCommands:        appCommands,
		internalCommands:   internalCommands,
		appViewConfig:      appViewConfig,
		internalViewConfig: internalViewConfig,
		templateAutoLoad:   templateAutoLoad,
	}
}

type Configer interface {
	Routes() []router.ControllerAction
	DiBindings() []DiCallback
	ConsoleCommands() map[string]commandexecutor.CommandItem
	Jobs() []cron.Job
	Middlewares() []interface{}
	ViewConfig() template.FuncMap
	GetTemplateAutoLoads() map[string][]string
}

type Conf struct {
	routes             []router.ControllerAction
	jobs               []cron.Job
	middlewares        []interface{}
	appBindings        []DiCallback
	internalBindings   []DiCallback
	appCommands        map[string]commandexecutor.CommandItem
	internalCommands   map[string]commandexecutor.CommandItem
	appViewConfig      template.FuncMap
	internalViewConfig template.FuncMap
	templateAutoLoad   map[string][]string
}

func (c *Conf) Routes() []router.ControllerAction {
	return c.routes
}

func (c *Conf) DiBindings() []DiCallback {
	return append(c.appBindings, c.internalBindings...)
}

func (c *Conf) ConsoleCommands() map[string]commandexecutor.CommandItem {
	return c.mergeCommands(c.appCommands, c.internalCommands)
}

func (c *Conf) Jobs() []cron.Job {
	return c.jobs
}

func (c *Conf) Middlewares() []interface{} {
	return c.middlewares
}

func (c *Conf) ViewConfig() template.FuncMap {
	mergedConfig := make(template.FuncMap)
	for key, value := range c.appViewConfig {
		mergedConfig[key] = value
	}

	for key, value := range c.internalViewConfig {
		mergedConfig[key] = value
	}

	return mergedConfig
}

func (c *Conf) GetTemplateAutoLoads() map[string][]string {
	return c.templateAutoLoad
}

func (*Conf) mergeCommands(maps ...map[string]commandexecutor.CommandItem) map[string]commandexecutor.CommandItem {
	mergedMap := make(map[string]commandexecutor.CommandItem)

	for _, mergeMap := range maps {
		for key, value := range mergeMap {
			mergedMap[key] = value
		}
	}

	return mergedMap
}
