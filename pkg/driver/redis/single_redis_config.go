package redis

import (
	"time"

	"example/pkg/log"

	"github.com/go-redis/redis/v8"
)

// SingleConnection -- redis connection
type SingleConnection struct {
	Address    string
	Password   string
	DB         int
	MaxRetries int
	PoolSize   int
}

// BuildClient -- build single redis client
func (conn *SingleConnection) BuildClient() (redis.UniversalClient, error) {
	if conn.Address == "" {
		return nil, ErrorMissingRedisAddress
	}

	log.Infof("[redis] single - address: %v, pass: %v, db: %v, pollSize: %v",
		conn.Address, "***", conn.DB, conn.PoolSize)

	return redis.NewClient(
		&redis.Options{
			Addr:        conn.Address,
			Password:    conn.Password, // no password set
			DB:          conn.DB,       // use default DB
			PoolSize:    conn.PoolSize,
			PoolTimeout: time.Second * 4,
			MaxRetries:  conn.MaxRetries,
		},
	), nil
}
