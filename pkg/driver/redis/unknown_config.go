package redis

import "github.com/go-redis/redis/v8"

// UnknownConnection --
type UnknownConnection struct {
}

// BuildClient --
func (conn *UnknownConnection) BuildClient() (redis.UniversalClient, error) {
	return nil, ErrorRedisClientNotSupported
}
