package logger

import (
	"os"

	"github.com/olbrichattila/gofra/pkg/app/db"
	"github.com/olbrichattila/gofra/pkg/app/storage"

	builder "github.com/olbrichattila/gosqlbuilder/pkg"
)

type LoggerStorageResolver interface {
	Construct(db.DBer, builder.Builder)
	GetLoggerStorage() storage.Storager
}

func NewSessionStorageResolver() LoggerStorageResolver {
	return &Resolver{}
}

type Resolver struct {
	storage storage.Storager
}

func (r *Resolver) Construct(db db.DBer, sqlBuilder builder.Builder) {
	storageName := os.Getenv("LOGGER_STORAGE")
	switch storageName {
	case "file":
		r.storage = storage.NewFileStorage()
	case "redis":
		r.storage = storage.NewRedisStorage()
	case "memcached":
		r.storage = storage.NewMemcacheStorage()
	case "db":
		r.storage = storage.NewDatabaseStorage("logs", db, sqlBuilder)
	default:
		r.storage = storage.NewFileStorage()
	}
}

func (r *Resolver) GetLoggerStorage() storage.Storager {
	return r.storage
}
