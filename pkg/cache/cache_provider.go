package cache

import (
	"errors"
	"time"

	"example/pkg/log"
)

var (
	ErrCacheMiss        = errors.New("cache miss")
	ErrDataTooLarge     = errors.New("data too large")
	ErrMalformedData    = errors.New("malformed data")
	ErrInvalidCacheTime = errors.New("invalid cache time")
	ErrMissingConfig    = errors.New("missing config")
	ErrMissingName      = errors.New("missing name")
)

type Cache interface {
	SetStruct(k string, v interface{}, exp time.Duration) error
	GetStruct(k string, v interface{}) error
	Set(k string, d []byte, exp time.Duration) error
	Get(k string) ([]byte, error)
	Del(k string) error
	TTL(k string) (time.Duration, bool)
	Close() error
}

type MetricAbleCache struct {
	cache          Cache
	cacheHitFunc   func()
	cacheMissFunc  func()
	totalCountFunc func()
}

func (m *MetricAbleCache) SetStruct(k string, v interface{}, exp time.Duration) error {
	return m.cache.SetStruct(k, v, exp)
}

func (m *MetricAbleCache) GetStruct(k string, v interface{}) error {
	defer m.totalCountFunc()
	err := m.cache.GetStruct(k, v)
	if err == nil {
		m.cacheHitFunc()
	} else if errors.Is(err, ErrCacheMiss) {
		m.cacheMissFunc()
	}
	return err
}

func (m *MetricAbleCache) Set(k string, d []byte, exp time.Duration) error {
	return m.cache.Set(k, d, exp)
}

func (m *MetricAbleCache) Get(k string) ([]byte, error) {
	defer m.totalCountFunc()
	res, err := m.cache.Get(k)
	if err == nil {
		m.cacheHitFunc()
	} else if errors.Is(err, ErrCacheMiss) {
		m.cacheMissFunc()
	}
	return res, err
}

func (m *MetricAbleCache) Del(k string) error {
	return m.cache.Del(k)
}

func (m *MetricAbleCache) TTL(k string) (time.Duration, bool) {
	return m.cache.TTL(k)
}

func (m *MetricAbleCache) Close() error {
	return m.cache.Close()
}

func defaultCacheHitFunc() {
}

func defaultCacheMissFunc() {

}

func defaultTotalCountFunc() {

}

func NewMetricAbleCache(cache Cache, cacheHitFunc, cacheMissFunc, totalCountFunc func()) *MetricAbleCache {
	if cacheHitFunc == nil {
		cacheHitFunc = defaultCacheHitFunc
	}
	if cacheMissFunc == nil {
		cacheMissFunc = defaultCacheMissFunc
	}
	if totalCountFunc == nil {
		totalCountFunc = defaultTotalCountFunc
	}
	if cache == nil {
		log.Fatalw("cache layer cannot be nil")
	}
	return &MetricAbleCache{
		cache:          cache,
		cacheHitFunc:   cacheHitFunc,
		cacheMissFunc:  cacheMissFunc,
		totalCountFunc: totalCountFunc,
	}
}
