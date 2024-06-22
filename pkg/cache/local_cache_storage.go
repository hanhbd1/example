package cache

import (
	"fmt"
	"time"

	"example/pkg/log"

	"github.com/dgraph-io/ristretto"
	"github.com/vmihailenco/msgpack/v5"
)

type InMemCacheStorage struct {
	cache   *ristretto.Cache
	maxSize int64
	prefix  string
}

const defaultMaxSize int64 = 100000000
const defaultCounter int64 = 10000
const defaultBufferSize int64 = 64

// NewInMemCacheStorage create new in-memory caching instance
func NewInMemCacheStorage(prefix string) (*InMemCacheStorage, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: defaultCounter,
		MaxCost:     defaultMaxSize,
		BufferItems: defaultBufferSize,
		Metrics:     true,
	})

	if err != nil {
		return nil, err
	}

	return &InMemCacheStorage{cache: cache, maxSize: defaultCounter, prefix: prefix}, nil
}

func Must(m *InMemCacheStorage, err error) *InMemCacheStorage {
	if err != nil {
		log.Fatalf(err.Error())
	}
	return m
}

func (c *InMemCacheStorage) key(k string) string {
	return fmt.Sprintf("%s:%s", c.prefix, k)
}

func (c *InMemCacheStorage) SetStruct(k string, v interface{}, exp time.Duration) error {
	bb, err := msgpack.Marshal(v)
	if err != nil {
		return err
	}

	return c.Set(k, bb, exp)
}

func (c *InMemCacheStorage) Set(k string, d []byte, exp time.Duration) error {
	if int64(len(d)) > c.maxSize {
		return ErrDataTooLarge
	}
	c.cache.SetWithTTL(c.key(k), d, int64(len(d)), exp)
	return nil
}

// GetStruct wrapper function that get data from cache then set to v, v must be pointer
func (c *InMemCacheStorage) GetStruct(k string, v interface{}) error {
	bb, err := c.Get(k)
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(bb, v)
}

func (c *InMemCacheStorage) Get(k string) ([]byte, error) {
	d, hit := c.cache.Get(c.key(k))
	if !hit {
		return nil, ErrCacheMiss
	}

	b, ok := d.([]byte)
	if !ok {
		return nil, ErrMalformedData
	}
	return b, nil
}

func (c *InMemCacheStorage) Del(k string) error {
	c.cache.Del(c.key(k))
	return nil
}

func (c *InMemCacheStorage) TTL(k string) (time.Duration, bool) {
	return c.cache.GetTTL(c.key(k))
}

func (c *InMemCacheStorage) Close() error {
	return nil
}
