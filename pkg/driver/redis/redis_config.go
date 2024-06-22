package redis

import (
	"errors"

	"example/pkg/config"
	"example/pkg/util"

	"github.com/go-redis/redis/v8"
)

// Connection --
type Connection interface {
	BuildClient() (redis.UniversalClient, error)
}

var (
	// ErrorMissingRedisAddress --
	ErrorMissingRedisAddress = errors.New("missing redis address")

	// ErrorRedisClientNotSupported --
	ErrorRedisClientNotSupported = errors.New("redis client not supported")
)

const (
	// Sentinel type
	Sentinel = "sentinel"
	// Cluster type
	Cluster = "cluster"
	// Single type
	Single = "single"

	// DefaultPoolSize --
	DefaultPoolSize   = 100
	DefaultMaxRetries = 2
)

// DefaultRedisConnectionFromConfig -- load connection settings in config with default key
func DefaultRedisConnectionFromConfig() Connection {
	redisClientType := config.GetString("redis.clientType")
	if util.IsStringEmpty(redisClientType) {
		redisClientType = Single
	}
	poolSize := config.GetInt("redis.poolSize")
	if poolSize <= 0 {
		poolSize = DefaultPoolSize
	}
	maxRetries := config.GetInt("redis.maxRetries")
	if maxRetries <= 0 {
		maxRetries = DefaultMaxRetries
	}

	switch redisClientType {
	case Sentinel:
		return &SentinelConnection{
			MasterGroup:       config.GetString("redis.sentinel.master"),
			SentinelAddresses: config.GetStringSlice("redis.sentinel.addresses"),
			Password:          config.GetString("redis.password"),
			DB:                config.GetInt("redis.db"),
			PoolSize:          poolSize,
			MaxRetries:        maxRetries,
		}
	case Cluster:
		return &ClusterConnection{
			ClusterAddresses: config.GetStringSlice("redis.cluster.addresses"),
			Password:         config.GetString("redis.password"),
			PoolSize:         poolSize,
			MaxRetries:       maxRetries,
		}
	case Single:
		return &SingleConnection{
			Address:    config.GetString("redis.address"),
			Password:   config.GetString("redis.password"),
			DB:         config.GetInt("redis.db"),
			PoolSize:   poolSize,
			MaxRetries: maxRetries,
		}
	default:
		return &UnknownConnection{}
	}
}

// NewRedisConfig --
func NewRedisConfig(add string, db int) Connection {
	return &SingleConnection{
		Address:  add,
		DB:       db,
		PoolSize: DefaultPoolSize,
	}
}

// NewRedisConfigWithPool --
func NewRedisConfigWithPool(add string, db, poolSize int) Connection {
	return &SingleConnection{
		Address:  add,
		DB:       db,
		PoolSize: poolSize,
	}
}
