package database

import (
	"github.com/bradfitz/gomemcache/memcache"
)

func MemcachedStore(key string, val string) {
	mc := memcache.New("localhost:11211")
	mc.Set(&memcache.Item{Key: key, Value: []byte(val)})
}

func MemcachedRetrieve(key string) (string, error) {
	mc := memcache.New("localhost:11211")
	val, err := mc.Get(key)
	return string(val.Value), err
}