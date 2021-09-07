package session

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type redisSession struct {
	id      string
	expired int
	client  *redis.Client
	cache   map[string]interface{}
}

// Get get value
func (instance redisSession) Get(key string) (value interface{}, err error) {
	v, ok := instance.cache[key]
	if !ok {
		err = fmt.Errorf("no value")
	} else {
		value = v
	}
	return
}

// Set set value
func (instance *redisSession) Set(key string, value interface{}) {
	if nil == instance.cache {
		instance.cache = make(map[string]interface{})
	}

	instance.cache[key] = value
}

// Del delete value
func (instance *redisSession) Del(key string) {
	delete(instance.cache, key)
}

// GetID get session id
func (instance redisSession) GetID() string {
	return instance.id
}

// Expired update expired
func (instance redisSession) Expired(expired int) {
}

// IsExpired check whether expired
func (instance redisSession) IsExpired(toExpired int64) bool {
	return false
}

// Save save data to redis
func (instance redisSession) Save() error {
	b, e := json.Marshal(instance.cache)
	if e != nil {
		return e
	}
	return instance.client.Set(instance.id, b, time.Duration(instance.expired)*time.Second).Err()
}

var _ ISession = &redisSession{}

// newRedis get redis session
func newRedis(id string, client *redis.Client) *redisSession {
	return &redisSession{
		id:     id,
		client: client,
	}
}
