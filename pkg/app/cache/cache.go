package cache

import "github.com/olbrichattila/gofra/pkg/app/storage"

func New() Cacher {
	return &Cache{}
}

const (
	cacheStoragePath = "./cache/"
)

type CacheFunc func(pars ...interface{}) string

type Cacher interface {
	Construct(CacheStorageResolver)
	Cache(string, CacheFunc, ...interface{}) string
	Put(string, string) error
	Delete(string) error
	Get(string) (string, error)
	HasKey(string) (bool, error)
}

type Cache struct {
	storage storage.Storager
}

func (c *Cache) Cache(key string, fn CacheFunc, pars ...interface{}) string {
	fullKey := cacheStoragePath + key
	if has, _ := c.storage.HasKey(fullKey); has {
		resolved, err := c.storage.Get(fullKey)
		if err == nil {
			return resolved
		}
	}

	value := fn(pars...)
	c.storage.Put(fullKey, value)

	return value
}

func (c *Cache) Construct(resolver CacheStorageResolver) {
	c.storage = resolver.GetCacheStorage()
}

func (c *Cache) Put(key, value string) error {
	return c.storage.Put(cacheStoragePath+key, value)
}

func (c *Cache) Delete(key string) error {
	return c.storage.Delete(cacheStoragePath + key)
}

func (c *Cache) Get(key string) (string, error) {
	return c.storage.Get(cacheStoragePath + key)
}

func (c *Cache) HasKey(key string) (bool, error) {
	return c.storage.HasKey(cacheStoragePath + key)
}
