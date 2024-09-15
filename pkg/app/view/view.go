package view

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/olbrichattila/gofra/pkg/app/config"
	"github.com/olbrichattila/gofra/pkg/app/session"
)

const (
	ViewTypeHTML  = "view"
	ViewTypeEmail = "mail"
)

type Viewer interface {
	Construct(config.Configer, session.Sessioner)
	Render(string, any) string
	RenderView(string, any) string
	RenderViewWithSessionError(string, any) string
	RenderMail(string, any) string
	NewPath(...string)
	RenderToFile(string, string, any) error
	RenderMailToFile(string, string, any) error
	Funcs(template.FuncMap) Viewer
	LoadTemplateParts([]string) Viewer
}

func New() Viewer {
	return &View{
		path: []string{"app", "views"},
	}
}

type View struct {
	config                   config.Configer
	session                  session.Sessioner
	path                     []string
	funcs                    template.FuncMap
	extraTemplatesToAutoLoad []string
	viewType                 string
}

func (v *View) Construct(conf config.Configer, s session.Sessioner) {
	v.config = conf
	v.session = s
}

func (v *View) RenderViewToFile(fileName string, templateFileName string, params any) error {
	content := v.RenderView(templateFileName, params)
	return v.toFile(fileName, content)
}

func (v *View) RenderToFile(fileName string, templateFileName string, params any) error {
	content := v.Render(templateFileName, params)
	return v.toFile(fileName, content)
}

func (v *View) RenderMailToFile(fileName string, templateFileName string, params any) error {
	content := v.RenderMail(templateFileName, params)
	return v.toFile(fileName, content)
}

func (v *View) toFile(fileName, content string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func (v *View) RenderView(templateFileName string, params any) string {
	v.viewType = ViewTypeHTML
	return v.renderTemplate(templateFileName, params)
}

func (v *View) RenderViewWithSessionError(templateFileName string, params any) string {
	wrapper := map[string]interface{}{
		"lastError":           v.session.Get("lastError"),
		"lastValidationError": v.session.Get("lastValidationError"),
		"lastRequest":         v.session.Get("lastRequest"),
		"data":                params,
	}
	v.viewType = ViewTypeHTML
	return v.renderTemplate(templateFileName, wrapper)
}

func (v *View) Render(templateFileName string, params any) string {
	v.viewType = ""
	return v.renderTemplate(templateFileName, params)
}

func (v *View) renderTemplate(templateFileName string, params any) string {
	if len(templateFileName) == 0 {
		return ""
	}
	templates := []string{templateFileName}
	if v.extraTemplatesToAutoLoad != nil {
		templates = append(templates, v.extraTemplatesToAutoLoad...)
	}

	if v.viewType != "" {
		templateAutoLoads := v.config.GetTemplateAutoLoads()
		if templateConfig, ok := templateAutoLoads[v.viewType]; ok {
			templates = append(templates, templateConfig...)
		}
	}

	paths := v.addPath(templates)

	funcs := v.mergeFuncMap()
	tmpl, err := template.New("example").Funcs(funcs).ParseFiles(paths...)
	if err != nil {
		return err.Error()
	}

	// Create a bytes.Buffer to hold the rendered template
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, templates[0], params)
	if err != nil {
		return err.Error()
	}
	v.funcs = make(template.FuncMap, 0)
	v.extraTemplatesToAutoLoad = nil

	// Return the rendered template as a string
	return buf.String()
}

func (v *View) RenderMail(templateFileName string, params any) string {
	v.NewPath("app", "mails")
	v.viewType = ViewTypeEmail
	result := v.renderTemplate(templateFileName, params)
	v.NewPath("app", "views")
	return result
}

func (v *View) Funcs(fm template.FuncMap) Viewer {
	v.funcs = fm
	return v
}

func (v *View) LoadTemplateParts(templates []string) Viewer {
	v.extraTemplatesToAutoLoad = templates
	return v
}

func (v *View) addPath(templates []string) []string {
	result := make([]string, len(templates))
	viewsDir := filepath.Join(v.path...)

	for i, templateFileName := range templates {
		result[i] = viewsDir + "/" + templateFileName
	}

	return result
}

func (v *View) NewPath(p ...string) {
	v.path = p
}

func (v *View) mergeFuncMap() template.FuncMap {
	merged := make(template.FuncMap, 0)
	for funcName, value := range v.config.ViewConfig() {
		merged[funcName] = value
	}

	for funcName, value := range v.funcs {
		merged[funcName] = value
	}

	return merged
}
