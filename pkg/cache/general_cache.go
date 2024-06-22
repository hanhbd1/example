package cache

import (
	"errors"
	"time"

	"example/pkg/log"

	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/multierr"
)

type GeneralCache struct {
	distributeCache Cache
	localCache      Cache
}

func (s *GeneralCache) Close() error {
	return multierr.Combine(s.localCache.Close(), s.distributeCache.Close())
}

type Option func(c *GeneralCache) error

// WithLocalCache set local cache of cache storage
func WithLocalCache(lCache Cache) Option {
	return func(c *GeneralCache) error {
		c.localCache = lCache
		return nil
	}
}

// WithDistributeCache set local cache of cache storage
func WithDistributeCache(dCache Cache) Option {
	return func(c *GeneralCache) error {
		c.distributeCache = dCache
		return nil
	}
}

// New create a cache storage with option
func New(opts ...Option) (*GeneralCache, error) {
	c := &GeneralCache{}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (s *GeneralCache) GetStruct(k string, vStruct interface{}) error {
	return s.getStruct(k, vStruct, false)
}

func (s *GeneralCache) Get(k string) ([]byte, error) {
	return s.get(k, false)
}

func (s *GeneralCache) get(k string, remoteOnly bool) ([]byte, error) {
	var data []byte
	var ttl time.Duration
	var e error
	var fromRemote bool
	if s.localCache != nil && !remoteOnly {
		if data, e = s.localCache.Get(k); e != nil {
			if !errors.Is(e, ErrCacheMiss) {
				return nil, e
			}
		} else {
			goto decode
		}
	}

	data, e = s.distributeCache.Get(k)
	if e != nil {
		if errors.Is(e, ErrCacheMiss) {
			return nil, ErrCacheMiss
		}
		return nil, e
	}
	fromRemote = true
	ttl, _ = s.distributeCache.TTL(k)

decode:
	if fromRemote && s.localCache != nil {
		if e = s.localCache.Set(k, data, ttl); e != nil {
			log.Warnf("Cannot save data back to cache")
		}
	}

	return data, nil
}

func (s *GeneralCache) getStruct(k string, vStruct interface{}, remoteOnly bool) error {
	var e error
	var data []byte
	data, e = s.get(k, remoteOnly)
	if e != nil {
		return e
	}

	e = msgpack.Unmarshal(data, vStruct)
	if e != nil {
		return e
	}
	return nil
}

func (s *GeneralCache) SetStruct(k string, v interface{}, exp time.Duration) error {
	return s.setStruct(k, v, exp, false)
}

func (s *GeneralCache) Set(k string, d []byte, exp time.Duration) error {
	return s.set(k, d, exp, false)
}

func (s *GeneralCache) set(k string, d []byte, exp time.Duration, remoteOnly bool) error {
	var e error
	if exp <= 0 {
		return ErrInvalidCacheTime
	}
	e = s.distributeCache.Set(k, d, exp)
	if e != nil {
		return e
	}

	if s.localCache != nil && !remoteOnly {
		if e = s.localCache.Set(k, d, exp); e != nil {
			log.Warnf("Error when set local cache: %v", e)
		}
	}
	return nil
}

func (s *GeneralCache) TTL(k string) (time.Duration, bool) {
	return s.distributeCache.TTL(k)
}

func (s *GeneralCache) setStruct(k string, v interface{}, exp time.Duration, remoteOnly bool) error {
	raw, e := msgpack.Marshal(v)
	if e != nil {
		return e
	}
	return s.set(k, raw, exp, remoteOnly)
}

func (s *GeneralCache) Del(k string) error {
	return s.del(k, false)
}

func (s *GeneralCache) del(k string, remoteOnly bool) error {
	if s.localCache != nil && !remoteOnly {
		err := s.localCache.Del(k)
		if err != nil {
			log.Warnf("error when del local cache", "error", err)
		}
	}
	if err := s.distributeCache.Del(k); err != nil {
		return err
	}
	return nil
}
