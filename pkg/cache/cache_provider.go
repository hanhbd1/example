package cache

import (
	"errors"
	"time"
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
