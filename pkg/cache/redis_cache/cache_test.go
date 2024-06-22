package redis_cache

import (
	"testing"
	"time"

	"example/pkg/driver/redis"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

type TestObject struct {
	Name string
}

func TestRedisCacheUseCase(t *testing.T) {
	//this is a mock redis server
	redisServer, _ := miniredis.Run()
	defer redisServer.Close()
	// init redis client

	redisClient, err := redis.NewConnection(redis.NewRedisConfig(redisServer.Addr(), 0))
	assert.Nil(t, err)

	redisCache, err := New(WithClient(redisClient), WithName("test"))
	assert.Nil(t, err)
	to := &TestObject{
		Name: "test",
	}
	err = redisCache.SetStruct("test", to, 5*time.Minute)
	assert.Nil(t, err)

	var to2 TestObject
	err = redisCache.GetStruct("test", &to2)
	assert.Nil(t, err)
	assert.Equal(t, to.Name, to2.Name)

	// Test time Fly for expire cache
	redisServer.FastForward(6 * time.Minute)

	var to3 TestObject

	err = redisCache.GetStruct("test", &to3)
	assert.NotNil(t, err)

}
