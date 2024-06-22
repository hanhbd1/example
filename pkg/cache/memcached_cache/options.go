package memcached_cache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type Option func(c *Cache) error

// WithName set cache name of cache storage
func WithName(name string) Option {
	return func(c *Cache) error {
		c.name = name
		return nil
	}
}

// WithDefaultRedisConfig set redis client using default config path
func WithClient(client *memcache.Client) Option {
	return func(cache *Cache) error {
		cache.client = client
		return nil
	}
}
