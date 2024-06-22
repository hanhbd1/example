package cache

import (
	"time"

	"example/pkg/compress"

	"github.com/vmihailenco/msgpack/v5"
)

type CompressionCache struct {
	CompressFunc   func([]byte) []byte
	DecompressFunc func([]byte) ([]byte, error)
	Cache          Cache
}

func NewCompressCache(cache Cache, compressFunc func([]byte) []byte, decompressFunc func([]byte) ([]byte, error)) *CompressionCache {
	if compressFunc == nil {
		compressFunc = compress.SnappyCompress
	}
	if decompressFunc == nil {
		decompressFunc = compress.SnappyDecompress
	}
	return &CompressionCache{
		CompressFunc:   compressFunc,
		DecompressFunc: decompressFunc,
		Cache:          cache,
	}
}

func (c *CompressionCache) SetStruct(k string, v interface{}, exp time.Duration) error {
	bb, err := msgpack.Marshal(v)
	if err != nil {
		return err
	}

	return c.Set(k, bb, exp)
}

func (c *CompressionCache) GetStruct(k string, v interface{}) error {
	bb, err := c.Get(k)
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(bb, v)
}

func (c *CompressionCache) Set(k string, d []byte, exp time.Duration) error {
	compress := c.CompressFunc(d)
	return c.Cache.Set(k, compress, exp)
}

func (c *CompressionCache) Get(k string) ([]byte, error) {
	d, err := c.Cache.Get(k)
	if err != nil {
		return nil, err
	}
	b, err := c.DecompressFunc(d)
	return b, nil
}

func (c *CompressionCache) Del(k string) error {
	return c.Cache.Del(k)
}

func (c *CompressionCache) TTL(k string) (time.Duration, bool) {
	return c.Cache.TTL(k)
}

func (c *CompressionCache) Close() error {
	return c.Cache.Close()
}
