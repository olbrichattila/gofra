package internalconfig

import (
	"html/template"
	"net/url"
	"os"

	internalviewfunction "github.com/olbrichattila/gofra/pkg/app/view-functions"
)

var ViewFuncConfig = template.FuncMap{
	"urlEscape":    url.QueryEscape,
	"envVar":       os.Getenv,
	"renderErrors": internalviewfunction.RenderErrors,
	"renderError":  internalviewfunction.RenderError,
	"lastRequest":  internalviewfunction.RenderRequest,
}
