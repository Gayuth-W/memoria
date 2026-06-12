package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache() *RedisCache {

	client := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379",
		},
	)

	return &RedisCache{
		Client: client,
	}
}

func (r *RedisCache) Get(
	key string,
) (string, error) {

	return r.Client.Get(
		context.Background(),
		key,
	).Result()
}

func (r *RedisCache) Set(
	key string,
	value string,
) error {

	return r.Client.Set(
		context.Background(),
		key,
		value,
		5*time.Minute,
	).Err()
}
