package cache

import (
	"os"

	"github.com/olbrichattila/gofra/pkg/app/db"
	"github.com/olbrichattila/gofra/pkg/app/storage"

	builder "github.com/olbrichattila/gosqlbuilder/pkg"
)

type CacheStorageResolver interface {
	Construct(db.DBer, builder.Builder)
	GetCacheStorage() storage.Storager
}

func NewCacheStorageResolver() CacheStorageResolver {
	return &Resolver{}
}

type Resolver struct {
	storage storage.Storager
}

func (r *Resolver) Construct(db db.DBer, sqlBuilder builder.Builder) {
	storageName := os.Getenv("CACHE_STORAGE")
	switch storageName {
	case "file":
		r.storage = storage.NewFileStorage()
	case "redis":
		r.storage = storage.NewRedisStorage()
	case "memcached":
		r.storage = storage.NewMemcacheStorage()
	case "db":
		r.storage = storage.NewDatabaseStorage("caches", db, sqlBuilder)
	default:
		r.storage = storage.NewFileStorage()
	}
}

func (r *Resolver) GetCacheStorage() storage.Storager {
	return r.storage
}
