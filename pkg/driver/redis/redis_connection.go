package redis

import (
	"context"
	"fmt"
	"sync"

	"example/pkg/log"

	"github.com/go-redis/redis/v8"
)

// IsInSlot -
func IsInSlot(key string, slot redis.ClusterSlot) bool {
	s := Slot(key)
	return slot.Start <= s && s <= slot.End
}

// GetSlotID -
func GetSlotID(key string, slots []redis.ClusterSlot) string {
	s := Slot(key)
	for k := range slots {
		slot := slots[k]
		if slot.Start <= s && s <= slot.End {
			return fmt.Sprintf("%v-%v", slot.Start, slot.End)
		}
	}
	return ""
}

// McRedis --
type McRedis struct {
	redis.UniversalClient
	Slots []redis.ClusterSlot
	lock  sync.Once
}

type Opts struct {
	hooks []redis.Hook
}

// NewConnection -- open connection to db
func NewConnection(conn Connection) (*McRedis, error) {
	var err error

	var redisOpt Opts

	c, err := conn.BuildClient()
	if err != nil {
		log.Error("[redis] Could not build redis client, details: ", err)
		return nil, err
	}

	pong, err := c.Ping(context.Background()).Result()
	if err != nil {
		//logger.McLog.Error("[redis] Could not ping to redis, details: ", err)
		return nil, err
	}
	log.Info("[redis] Ping to redis: ", pong)

	// hook
	if len(redisOpt.hooks) > 0 {
		for _, h := range redisOpt.hooks {
			c.AddHook(h)
		}
	}

	cs := getClusterInfo(c)
	return &McRedis{UniversalClient: c, Slots: cs}, nil
}

func getClusterInfo(c redis.UniversalClient) []redis.ClusterSlot {
	var cs = make([]redis.ClusterSlot, 0)
	if ci := c.ClusterInfo(context.Background()); ci.Err() == nil {
		csr := c.ClusterSlots(context.Background())
		var err error
		cs, err = csr.Result()
		if err != nil {
			log.Error("[redis] Cannot get cluster slots")
		}
	}
	return cs
}

// NewConnectionFromExistedClient --
func NewConnectionFromExistedClient(c redis.UniversalClient) *McRedis {
	cs := getClusterInfo(c)
	return &McRedis{UniversalClient: c, Slots: cs}
}

func (r *McRedis) Name() string {
	return "Redis"
}

func (r *McRedis) Check() error {
	return r.Ping(context.Background()).Err()
}

// Close -- close connection
func (r *McRedis) Close() error {
	var err error
	r.lock.Do(func() {
		err = r.UniversalClient.Close()
	})
	return err
}

// GetClient --
func (r *McRedis) GetClient() redis.UniversalClient {
	return r.UniversalClient
}

func (r *McRedis) IsSingle() bool {
	switch r.UniversalClient.(type) {
	case *redis.Client:
		return true
	default:
		return false
	}
}

// GetClusterSlots -
func (r *McRedis) GetClusterSlots() ([]redis.ClusterSlot, error) {
	res := r.ClusterSlots(context.Background())
	return res.Result()
}

// GetRedisSlot -
func (r *McRedis) GetRedisSlot(key string) int {
	return Slot(key)
}

// GetRedisSlotID -
func (r *McRedis) GetRedisSlotID(key string) string {
	return GetSlotID(key, r.Slots)
}
