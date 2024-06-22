package redis

import (
	"time"

	"example/pkg/log"

	"github.com/go-redis/redis/v8"
)

// ClusterConnection -- redis connection
type ClusterConnection struct {
	ClusterAddresses []string
	Password         string
	PoolSize         int
	MaxRetries       int
}

// BuildClient --
func (conn *ClusterConnection) BuildClient() (redis.UniversalClient, error) {
	if len(conn.ClusterAddresses) == 0 {
		return nil, ErrorMissingRedisAddress
	}
	log.Infof("[redis] Create cluster client to %v", conn.ClusterAddresses)

	redisDB := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:       conn.ClusterAddresses,
		Password:    conn.Password,
		PoolSize:    conn.PoolSize,
		PoolTimeout: time.Second * 4,
		MaxRetries:  conn.MaxRetries,
	})

	return redisDB, nil
}
