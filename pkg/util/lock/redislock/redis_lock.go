package redislock

import (
	"context"
	"fmt"
	"time"

	"example/pkg/driver/redis"
	"example/pkg/util"
	"example/pkg/util/lock"

	goredis "github.com/go-redis/redis/v8"
)

var (
	luaRelease = goredis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then` +
		` return redis.call("del", KEYS[1]) else return 0 end`)
)

type RedisLocker struct {
	client        *redis.McRedis
	prefix        string
	maxObtainTime time.Duration
}

type RedisLock struct {
	client *redis.McRedis
	key    string
	value  string
}

func NewRedisLock(client *redis.McRedis, prefix string, maxObtainTime time.Duration) *RedisLocker {
	return &RedisLocker{
		client:        client,
		prefix:        prefix,
		maxObtainTime: maxObtainTime,
	}
}

func (r *RedisLocker) Obtain(ctx context.Context, key string, value string, obtainTime ...time.Duration) (lock.Lock, error) {
	var rKey string
	if r.prefix != "" {
		rKey = fmt.Sprintf("%s:%s", r.prefix, key)
	} else {
		rKey = key
	}
	lockTime := r.maxObtainTime
	if len(obtainTime) == 1 {
		lockTime = obtainTime[0]
	}
	ok, err := r.client.SetNX(ctx, rKey, value, lockTime).Result()
	if ok {
		return &RedisLock{client: r.client, key: rKey, value: value}, nil
	} else if err == nil {
		return nil, lock.ErrLockObtained
	} else {
		return nil, err
	}
}

func (r *RedisLock) Release(ctx context.Context) error {
	if util.IsStringEmpty(r.key) {
		return nil
	}
	res, err := luaRelease.Run(ctx, r.client, []string{r.key}, r.value).Result()
	if err == goredis.Nil {
		return nil
	} else if err != nil {
		return err
	}

	if i, ok := res.(int64); !ok || i != 1 {
		return lock.ErrLockNotHeld
	}
	return nil
}
