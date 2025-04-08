package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisClient(addr string, ttl time.Duration) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	//Нужно добавить проверку Redis`a

	return &RedisCache{
		client: client,
		ttl:    ttl,
	}, nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key string, value string) error {
	return r.client.Set(ctx, key, value, r.ttl).Err()
}

func (r *RedisCache) GetClient() *redis.Client {
	return r.client
}
