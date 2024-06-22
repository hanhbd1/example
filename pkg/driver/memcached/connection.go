package memcached

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type Connection struct {
	Addr []string
}

func New(conn Connection) *memcache.Client {
	return memcache.New(conn.Addr...)
}
