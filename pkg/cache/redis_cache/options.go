package redis_cache

import (
	"context"

	"example/pkg/driver/redis"
)

type Option func(c *Cache) error

// WithContext set local cache of cache storage
func WithContext(ctx context.Context) Option {
	return func(c *Cache) error {
		c.ctx = ctx
		return nil
	}
}

// WithName set cache name of cache storage
func WithName(name string) Option {
	return func(c *Cache) error {
		c.name = name
		return nil
	}
}

// WithClient set redis client using default config path
func WithClient(mcRedis *redis.McRedis) Option {
	return func(cache *Cache) error {
		cache.redisClient = mcRedis
		return nil
	}
}
