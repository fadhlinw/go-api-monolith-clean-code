package lib

import (
	"context"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// RedisClient struct holds the Redis client instance and context
type RedisClient struct {
	logger Logger
	client *redis.Client
	ctx    context.Context
}

// NewRedisClient creates a new Redis client instance
func NewRedisClient(env Env, logger Logger) RedisClient {
	return RedisClient{
		logger: Logger{},
		client: redis.NewClient(&redis.Options{
			Addr:     env.RedisHost + ":" + env.RedisPort,
			Password: env.RedisPassword,
			DB:       0,
		}),
		ctx: context.Background(),
	}
}

// SetData stores a value with a key in Redis with a TTL
func (r RedisClient) SetData(key string, value string) error {
	err := r.client.Set(r.ctx, key, value, time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetData retrieves the value of a given key from Redis
func (r RedisClient) GetData(key string) ([]string, error) {
	val, err := r.client.LRange(r.ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

// RemoveFromList removes a value from a list in Redis
func (r RedisClient) RemoveFromList(key string, value string) error {
	err := r.client.LRem(r.ctx, key, 0, value).Err()
	if err != nil {
		return err
	}
	return nil
}
