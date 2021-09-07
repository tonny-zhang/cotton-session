package session

import (
	"encoding/json"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

type memcacheSession struct {
	id      string
	expired int
	client  *memcache.Client
	cache   map[string]interface{}
}

// Get get value
func (instance memcacheSession) Get(key string) (value interface{}, err error) {
	v, ok := instance.cache[key]
	if !ok {
		err = fmt.Errorf("no value")
	} else {
		value = v
	}
	return
}

// Set set value
func (instance *memcacheSession) Set(key string, value interface{}) {
	if nil == instance.cache {
		instance.cache = make(map[string]interface{})
	}

	instance.cache[key] = value
}

// Del delete value
func (instance *memcacheSession) Del(key string) {
	delete(instance.cache, key)
}

// GetID get session id
func (instance memcacheSession) GetID() string {
	return instance.id
}

// Expired update expired
func (instance memcacheSession) Expired(expired int) {
}

// IsExpired check whether expired
func (instance memcacheSession) IsExpired(toExpired int64) bool {
	return false
}

// Save save data to redis
func (instance memcacheSession) Save() error {
	b, e := json.Marshal(instance.cache)
	if e != nil {
		return e
	}
	return instance.client.Set(&memcache.Item{
		Key:        instance.id,
		Value:      b,
		Expiration: int32(instance.expired),
		Flags:      0,
	})
}

var _ ISession = &redisSession{}

// newMemcache get memcache session
func newMemcache(id string, client *memcache.Client) *memcacheSession {
	return &memcacheSession{
		id:     id,
		client: client,
	}
}
