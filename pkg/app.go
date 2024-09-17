// Package app
package app

// TODO: Refactor, split up, getting too big
import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	commandexecutor "github.com/olbrichattila/gofra/pkg/app/command"
	"github.com/olbrichattila/gofra/pkg/app/config"
	"github.com/olbrichattila/gofra/pkg/app/cron"
	"github.com/olbrichattila/gofra/pkg/app/env"
	"github.com/olbrichattila/gofra/pkg/app/router"
	"github.com/olbrichattila/gofra/pkg/app/validator"
	internalconfig "github.com/olbrichattila/gofra/pkg/internal-config"

	"github.com/olbrichattila/godi"
)

func New(
	container godi.Container,
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
) *App {
	app := &App{
		di:                  container,
		validationRules:     validationRules,
		validationRuleFuncs: validationRuleFuncs,
		conf: config.New(
			routes,
			jobs,
			middlewares,
			appBindings,
			internalconfig.DiBindings,
			consoleCommands,
			internalconfig.ConsoleCommands,
			appViewConfig,
			internalconfig.ViewFuncConfig,
			templateAutoLoad,
		),
	}

	app.initBindings()

	_, err := app.di.Get(app)
	if err != nil {
		panic(err.Error())
	}

	_, err = app.di.Call(bootstrapFunc)
	if err != nil {
		panic(err.Error())
	}

	return app
}

type App struct {
	di                  godi.Container
	router              router.Router
	conf                config.Configer
	commandExecutor     commandexecutor.CommandExecutor
	validationRules     map[string]validator.ValidationRule
	validationRuleFuncs map[string]validator.RuleFunc
}

func (a *App) Construct(
	_ env.Enver, // It will automatically loads env with it's constructor
	cron cron.JobTimer,
	router router.Router,
	ce commandexecutor.CommandExecutor,
) {
	cron.Init(a.di, a.conf.Jobs())
	a.router = router
	a.commandExecutor = ce
}

func (a *App) Serve() {
	port, err := a.getPort()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	hTTPHandler := &hTTPHandler{app: a}
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/", hTTPHandler)

	err = a.listenHttps()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("running on HTTP" + port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (a *App) getPort() (string, error) {
	port := os.Getenv("HTTP_LISTENING_PORT")
	if port == "" {
		return ":80", nil
	}

	if _, err := strconv.Atoi(port); err == nil {
		return ":" + port, nil
	}

	return "", fmt.Errorf("port %s provided is not a number", port)
}

func (a *App) getHttpsPort() (string, error) {
	port := os.Getenv("HTTPS_LISTENING_PORT")
	if port == "" {
		return ":443", nil
	}

	if _, err := strconv.Atoi(port); err == nil {
		return ":" + port, nil
	}

	return "", fmt.Errorf("port %s provided is not a number", port)
}

func (a *App) getHttpsEnabled() bool {
	enabled := os.Getenv("HTTPS_ENABLED")
	enabled = strings.ToLower(strings.TrimSpace(enabled))
	if enabled == "" {
		return false
	}

	if enabled == "1" || enabled == "true" {
		return true
	}

	return false
}

func (a *App) getCertFile() (string, error) {
	fileName := os.Getenv("HTTPS_CERT_FILE")
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		return "", fmt.Errorf("the HTTPS_CERT_FILE environment variable missing")
	}

	if !a.fileExists(fileName) {
		return "", fmt.Errorf("the certificate file %s does not exists", fileName)
	}

	return fileName, nil
}

func (a *App) getPrivateKeyFile() (string, error) {
	fileName := os.Getenv("HTTPS_PRIVATE_KEY_FILE")
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		return "", fmt.Errorf("the HTTPS_PRIVATE_KEY_FILE environment variable missing")
	}

	if !a.fileExists(fileName) {
		return "", fmt.Errorf("the private key file %s does not exists", fileName)
	}

	return fileName, nil
}

func (a *App) Command() {
	err := a.commandExecutor.Execute(a.di, a.conf.ConsoleCommands())
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (a *App) listenHttps() error {
	if a.getHttpsEnabled() {
		httpsPort, err := a.getHttpsPort()
		if err != nil {
			return err
		}

		certFile, err := a.getCertFile()
		if err != nil {
			return err
		}

		privateKeyFile, err := a.getPrivateKeyFile()
		if err != nil {
			return err
		}

		fmt.Println("running on HTTPS" + httpsPort)
		go func() {
			err = http.ListenAndServeTLS(httpsPort, certFile, privateKeyFile, nil)
			if err != nil {
				fmt.Println(err.Error())
			}
		}()
	}

	return nil
}

func (a *App) initBindings() {
	for _, cbFunc := range a.conf.DiBindings() {
		key, binding, err := cbFunc(a.di)
		if err != nil {
			panic(err.Error())
		}
		a.di.Set(key, binding)
	}
	a.di.Set("olbrichattila.gofra.pkg.app.config.Configer", a.conf)
	a.di.Set("olbrichattila.godi.Container", a.di)
}

func (a *App) fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
