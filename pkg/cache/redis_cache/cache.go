package redis_cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"example/pkg/cache"
	"example/pkg/driver/redis"

	goredis "github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
)

type Cache struct {
	ctx         context.Context
	name        string
	redisClient *redis.McRedis
}

func (s *Cache) key(k string) string {
	return fmt.Sprintf("%s:%s", s.name, k)
}

func (s *Cache) GetStruct(k string, vStruct interface{}) error {
	var e error
	var data []byte
	data, e = s.Get(k)
	if e != nil {
		return e
	}

	e = msgpack.Unmarshal(data, vStruct)
	if e != nil {
		return e
	}
	return nil
}

func (s *Cache) Get(k string) ([]byte, error) {
	var data []byte
	var e error
	data, e = s.redisClient.Get(s.ctx, s.key(k)).Bytes()
	if e != nil {
		if errors.Is(e, goredis.Nil) {
			return nil, cache.ErrCacheMiss
		}

		return nil, e
	}
	return data, nil
}

func (s *Cache) SetStruct(k string, v interface{}, exp time.Duration) error {
	raw, e := msgpack.Marshal(v)
	if e != nil {
		return e
	}
	return s.Set(k, raw, exp)
}

func (s *Cache) Set(k string, d []byte, exp time.Duration) error {
	var e error
	if exp <= 0 {
		return cache.ErrInvalidCacheTime
	}
	e = s.redisClient.Set(s.ctx, s.key(k), d, exp).Err()
	if e != nil {
		return e
	}
	return nil
}

func (s *Cache) Del(k string) error {
	return s.redisClient.Del(s.ctx, s.key(k)).Err()
}

func (s *Cache) TTL(k string) (time.Duration, bool) {
	t, err := s.redisClient.TTL(s.ctx, s.key(k)).Result()
	if err != nil {
		return 0, false
	}
	return t, t == -2
}

// New create a cache storage with option
func New(opts ...Option) (*Cache, error) {
	c := &Cache{}

	for _, opt := range opts {

		if err := opt(c); err != nil {
			return nil, err
		}
	}
	if c.ctx == nil {
		c.ctx = context.Background()
	}
	if c.redisClient == nil {
		return nil, cache.ErrMissingConfig
	}
	if c.name == "" {
		return nil, cache.ErrMissingName
	}
	return c, nil
}

func (s *Cache) Close() error {
	if s.redisClient != nil {
		return s.redisClient.Close()
	}
	return nil
}
