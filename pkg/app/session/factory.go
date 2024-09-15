package session

import (
	"os"

	"github.com/olbrichattila/gofra/pkg/app/db"
	"github.com/olbrichattila/gofra/pkg/app/storage"

	builder "github.com/olbrichattila/gosqlbuilder/pkg"
)

type SessionStorageResolver interface {
	Construct(db.DBer, builder.Builder)
	GetSessionStorage() storage.Storager
}

func NewSessionStorageResolver() SessionStorageResolver {
	return &Resolver{}
}

type Resolver struct {
	storage storage.Storager
}

func (r *Resolver) Construct(db db.DBer, sqlBuilder builder.Builder) {
	storageName := os.Getenv("SESSION_STORAGE")
	switch storageName {
	case "file":
		r.storage = storage.NewFileStorage()
	case "redis":
		r.storage = storage.NewRedisStorage()
	case "memcached":
		r.storage = storage.NewMemcacheStorage()
	case "db":
		r.storage = storage.NewDatabaseStorage("sessions", db, sqlBuilder)
	default:
		r.storage = storage.NewFileStorage()
	}
}

func (r *Resolver) GetSessionStorage() storage.Storager {
	return r.storage
}
