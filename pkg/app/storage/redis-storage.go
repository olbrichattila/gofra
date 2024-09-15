package storage

import (
	"context"
	"os"
	"strconv"

	redis "github.com/go-redis/redis/v8"
)

func NewRedisStorage() Storager {
	defaultDb := 0
	defaultPort := 6379

	db := os.Getenv("REDIS_DB")
	if db != "" {
		dbId, err := strconv.Atoi(db)
		if err == nil {
			defaultDb = dbId
		}
	}

	port := os.Getenv("REDIS_PORT")
	if port != "" {
		configPort, err := strconv.Atoi(port)
		if err == nil {
			defaultPort = configPort
		}
	}

	redisAddr := os.Getenv("REDIS_SERVER_HOST") + ":" + strconv.Itoa(defaultPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       defaultDb,
	})

	return &RedisStore{
		rdb: rdb,
	}
}

type RedisStore struct {
	rdb *redis.Client
}

func (s *RedisStore) Append(key string, value string) error {
	return s.Put(key, value)
}

func (s *RedisStore) Put(key string, value string) error {
	return s.rdb.Set(context.Background(), key, value, 0).Err()
}

func (s *RedisStore) Delete(key string) error {
	return s.rdb.Del(context.Background(), key).Err()
}

func (s *RedisStore) HasKey(key string) (bool, error) {
	s.rdb.Exists(context.Background(), key).Result()
	exists, err := s.rdb.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}

	if exists == 1 {
		return true, nil
	}

	return false, nil
}

func (s *RedisStore) Get(key string) (string, error) {
	return s.rdb.Get(context.Background(), key).Result()
}
