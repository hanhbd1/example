package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestObject struct {
	Name string
}

func TestCacheUseCase(t *testing.T) {
	var cacheInstance Cache
	cacheInstance, err := NewInMemCacheStorage("test")
	assert.Nil(t, err)

	to := &TestObject{
		Name: "test",
	}
	err = cacheInstance.SetStruct("test", to, 5*time.Minute)
	assert.Nil(t, err)

	time.Sleep(1 * time.Millisecond) // sleep for 1 millisecond to make sure the cache is set

	var to2 TestObject
	err = cacheInstance.GetStruct("test", &to2)
	assert.Nil(t, err)
	assert.Equal(t, to.Name, to2.Name)

}
