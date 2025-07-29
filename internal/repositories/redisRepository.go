package repositories

import (
	"GinBox/config"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// IRedisRepository defines the interface how to communicate with Redis via its client.
type IRedisRepository interface {
	Set(ctx context.Context, key string, value string, ttlSeconds int64) error
	Get(ctx context.Context, key string) (string, error)
}

// RedisRepository implements IRedisRepository and provides high-level usable methods.
type RedisRepository struct {
	client *redis.Client
}

// NewRedisRepository returns a new RedisRepository with input config.
func NewRedisRepository(cfg *config.Config) *RedisRepository {
	return &RedisRepository{
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Host,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}),
	}
}

// Set associates particular key with particular value.
func (r *RedisRepository) Set(ctx context.Context, key string, value string, ttlSeconds int64) error {
	return r.client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

// Get returns a result of getting value by input key.
func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}
