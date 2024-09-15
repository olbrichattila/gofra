package internalconfig

import (
	"os"

	"github.com/olbrichattila/gofra/pkg/app/args"
	"github.com/olbrichattila/gofra/pkg/app/cache"
	commandexecutor "github.com/olbrichattila/gofra/pkg/app/command"
	"github.com/olbrichattila/gofra/pkg/app/config"
	"github.com/olbrichattila/gofra/pkg/app/cron"
	"github.com/olbrichattila/gofra/pkg/app/db"
	"github.com/olbrichattila/gofra/pkg/app/env"
	"github.com/olbrichattila/gofra/pkg/app/event"
	"github.com/olbrichattila/gofra/pkg/app/logger"
	"github.com/olbrichattila/gofra/pkg/app/mail"
	"github.com/olbrichattila/gofra/pkg/app/queue"
	"github.com/olbrichattila/gofra/pkg/app/request"
	"github.com/olbrichattila/gofra/pkg/app/router"
	"github.com/olbrichattila/gofra/pkg/app/session"
	"github.com/olbrichattila/gofra/pkg/app/validator"
	"github.com/olbrichattila/gofra/pkg/app/view"
	wizard "github.com/olbrichattila/gofra/pkg/app/wizards/class"
	commandcreator "github.com/olbrichattila/gofra/pkg/app/wizards/command"

	"github.com/olbrichattila/godi"
	gosqlbuilder "github.com/olbrichattila/gosqlbuilder"
	pkg "github.com/olbrichattila/gosqlbuilder/pkg"
)

func getOpenedDb() interface{} {
	return db.New()
}

func getSqlBuilder() interface{} {
	dbConnection := os.Getenv(db.EnvdbConnection)
	builder := gosqlbuilder.New()

	switch dbConnection {
	case db.DbConnectionTypeSqLite:
		builder.SetSQLFlavour(pkg.FlavourSqLite)
	case db.DbConnectionTypeMySQL:
		builder.SetSQLFlavour(pkg.FlavourMySQL)
	case db.DbConnectionTypePgSQL:
		builder.SetSQLFlavour(pkg.FlavourPgSQL)
	case db.DbConnectionTypeFirebird:
		builder.SetSQLFlavour(pkg.FlavourFirebirdSQL)
	}

	return builder
}

var DiBindings = []config.DiCallback{
	func(di godi.Container) (string, interface{}, error) {
		env, err := di.Get(env.New())
		return "olbrichattila.gofra.pkg.app.env.Enver", env, err
	},
	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.args.CommandArger", args.New(), nil
	},
	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.router.Router", router.NewRouter(), nil
	},
	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.view.Viewer", view.New(), nil
	},
	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.request.Requester", request.New(), nil
	},
	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.validator.Validator", validator.New(), nil
	},
	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.db.DBFactoryer", db.NewDBFactory(), nil
	},
	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.db.DBer", getOpenedDb, nil
	},
	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gosqlbuilder.pkg.Builder", getSqlBuilder, nil
	},
	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.logger.LoggerStorageResolver", logger.NewSessionStorageResolver(), nil
	},
	func(di godi.Container) (string, interface{}, error) {
		l, err := di.Get(logger.New())
		return "olbrichattila.gofra.pkg.app.logger.Logger", l, err
	},
	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.session.SessionStorageResolver", session.NewSessionStorageResolver(), nil
	},
	func(di godi.Container) (string, interface{}, error) {
		session, err := di.Get(session.New())
		return "olbrichattila.gofra.pkg.app.session.Sessioner", session, err
	},

	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.cache.CacheStorageResolver", cache.NewCacheStorageResolver(), nil
	},
	func(di godi.Container) (string, interface{}, error) {
		session, err := di.Get(cache.New())
		return "olbrichattila.gofra.pkg.app.cache.Cacher", session, err
	},

	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.cron.JobTimer", cron.New(), nil
	},
	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.queue.Quer", queue.New(), nil
	},
	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.mail.Mailer", mail.New(), nil
	},

	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.wizards.command.CommandCreator", commandcreator.New(), nil
	},
	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.command.CommandExecutor", commandexecutor.New(), nil
	},
	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.wizards.class.ClassCreator", wizard.NewClassCreator(), nil
	},
	func(di godi.Container) (string, interface{}, error) {
		return "olbrichattila.gofra.pkg.app.event.Eventer", event.NewLocalEvent(), nil
	},
}
