package storage

import (
	"os"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
)

func NewMemcacheStorage() Storager {
	defaultPort := 11211
	host := os.Getenv("MEMCACHE_HOST")
	port := os.Getenv("MEMCACHE_PORT")
	if port != "" {
		configPort, err := strconv.Atoi(port)
		if err == nil {
			defaultPort = configPort
		}
	}

	memcacheConnect := host + ":" + strconv.Itoa(defaultPort)

	return &MemcacheStore{
		mc: memcache.New(memcacheConnect),
	}
}

type MemcacheStore struct {
	mc *memcache.Client
}

func (s *MemcacheStore) Append(key string, value string) error {
	return s.Put(key, value)
}

func (s *MemcacheStore) Put(key string, value string) error {
	err := s.mc.Set(&memcache.Item{Key: key, Value: []byte(value)})
	if err != nil {
		return err
	}

	return nil
}

func (s *MemcacheStore) Delete(key string) error {
	err := s.mc.Delete(key)
	if err == memcache.ErrCacheMiss {
		return nil
	} else if err != nil {
		return err
	}

	return nil
}

func (s *MemcacheStore) HasKey(key string) (bool, error) {
	_, err := s.mc.Get(key)
	if err == memcache.ErrCacheMiss {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (s *MemcacheStore) Get(key string) (string, error) {
	item, err := s.mc.Get(key)
	if err != nil {
		return "", err
	}

	return string(item.Value), nil
}
