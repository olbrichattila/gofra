package internalcommand

import wizard "github.com/olbrichattila/gofra/pkg/app/wizards/class"

func getDefaultCreateDiMapping() map[string]wizard.ParameterInfo {
	return map[string]wizard.ParameterInfo{
		"request":    {Name: "request.Requester", Alias: "r", ImportPath: "\"github.com/olbrichattila/gofra/pkg/app/request\""},
		"db":         {Name: "db.DBer", Alias: "db", ImportPath: "\"github.com/olbrichattila/gofra/pkg/app/db\""},
		"logger":     {Name: "logger.Logger", Alias: "l", ImportPath: "\"github.com/olbrichattila/gofra/pkg/app/logger\""},
		"sqlBuilder": {Name: "builder.Builder", Alias: "sqlBuilder", ImportPath: "builder \"github.com/olbrichattila/gosqlbuilder/pkg\""},
		"session":    {Name: "session.Sessioner", Alias: "s", ImportPath: "\"github.com/olbrichattila/gofra/pkg/app/session\""},
		"view":       {Name: "view.Viewer", Alias: "v", ImportPath: "\"github.com/olbrichattila/gofra/pkg/app/view\""},
		"mail":       {Name: "mail.Mailer", Alias: "m", ImportPath: "\"github.com/olbrichattila/gofra/pkg/app/mail\""},
		"config":     {Name: "config.Configer", Alias: "c", ImportPath: "\"github.com/olbrichattila/gofra/pkg/app/config\""},
		"response":   {Name: "http.ResponseWriter", Alias: "w", ImportPath: "\"net/http\""},
		"cargs":      {Name: "args.CommandArger", Alias: "a", ImportPath: "\"github.com/olbrichattila/gofra/pkg/app/args\""},
		"queue":      {Name: "queue.Quer", Alias: "q", ImportPath: "\"github.com/olbrichattila/gofra/pkg/app/queue\""},
	}
}
