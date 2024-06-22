package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
)

// SentinelConnection -- redis connection
type SentinelConnection struct {
	MasterGroup       string
	SentinelAddresses []string
	Password          string
	DB                int
	PoolSize          int
	MaxRetries        int
}

// BuildClient --
func (conn *SentinelConnection) BuildClient() (redis.UniversalClient, error) {
	if len(conn.SentinelAddresses) == 0 {
		return nil, ErrorMissingRedisAddress
	}

	masterGroup := conn.MasterGroup
	if masterGroup == "" {
		masterGroup = "master"
	}

	redisdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterGroup,
		SentinelAddrs: conn.SentinelAddresses,
		Password:      conn.Password,
		DB:            conn.DB,
		PoolSize:      conn.PoolSize,
		PoolTimeout:   time.Second * 4,
		MaxRetries:    conn.MaxRetries,
	})

	return redisdb, nil
}
