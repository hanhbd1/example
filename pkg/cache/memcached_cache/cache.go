package memcached_cache

import (
	"errors"
	"fmt"
	"time"

	"example/pkg/cache"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/spf13/cast"
	"github.com/vmihailenco/msgpack/v5"
)

type Cache struct {
	name   string
	client *memcache.Client
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
	item, e := s.client.Get(s.key(k))
	if e != nil {
		if errors.Is(e, memcache.ErrCacheMiss) {
			return nil, cache.ErrCacheMiss
		}

		return nil, e
	}
	return item.Value, nil
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
	e = s.client.Set(&memcache.Item{
		Key:        s.key(k),
		Value:      d,
		Expiration: cast.ToInt32(time.Now().Add(exp).Unix()),
	})
	if e != nil {
		return e
	}
	return nil
}

func (s *Cache) Del(k string) error {
	return s.client.Delete(s.key(k))
}

func (s *Cache) TTL(k string) (time.Duration, bool) {
	item, err := s.client.Get(s.key(k))
	if err != nil {
		return 0, false
	}
	if item.Expiration == 0 {
		return 0, true
	}
	if item.Expiration < 60*60*24*30 {
		return time.Duration(item.Expiration) * time.Second, true
	}
	return time.Unix(int64(item.Expiration), 0).Sub(time.Now()), true
}

// New create a cache storage with option
func New(opts ...Option) (*Cache, error) {
	c := &Cache{}

	for _, opt := range opts {

		if err := opt(c); err != nil {
			return nil, err
		}
	}

	if c.client == nil {
		return nil, cache.ErrMissingConfig
	}

	if c.name == "" {
		return nil, cache.ErrMissingName
	}

	return c, nil
}

func (s *Cache) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}
