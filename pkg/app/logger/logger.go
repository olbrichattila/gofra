package logger

import (
	"encoding/json"
	"time"

	"github.com/olbrichattila/gofra/pkg/app/storage"
)

const (
	typeInfo     = "info"
	typeWarning  = "warning"
	typeError    = "error"
	typeCritical = "critical"
)

func New() Logger {
	return &Logging{
		filename: "log/app.log", // TODO make it configurable from conf.
	}
}

type Logger interface {
	Construct(LoggerStorageResolver)
	Info(string)
	Warning(string)
	Error(string)
	Critical(string)
}

type Logging struct {
	storage  storage.Storager
	filename string
}

func (s *Logging) Construct(loggerResolver LoggerStorageResolver) {
	s.storage = loggerResolver.GetLoggerStorage()
}

func (l *Logging) Info(message string) {
	l.log(typeInfo, message)
}

func (l *Logging) Warning(message string) {
	l.log(typeWarning, message)
}

func (l *Logging) Error(message string) {
	l.log(typeError, message)
}

func (l *Logging) Critical(message string) {
	l.log(typeCritical, message)
}

func (l *Logging) log(ltype, message string) {

	logStr := map[string]string{
		"type":       ltype,
		"message":    message,
		"created_at": time.Now().Format("2006-01-02 15:04:05"),
	}

	logLine, err := json.Marshal(logStr)
	if err != nil {
		return
	}

	l.storage.Append(l.filename, string(logLine)+"\n")
}
