package session

import (
	"fmt"
	"time"
)

type memorySession struct {
	id      string
	expired int64
	cache   map[string]interface{}
}

// Get get value
func (instance memorySession) Get(key string) (value interface{}, err error) {
	v, ok := instance.cache[key]
	if !ok {
		err = fmt.Errorf("no value")
	} else {
		value = v
	}
	return
}

// Set set value
func (instance *memorySession) Set(key string, value interface{}) {
	instance.cache[key] = value
}

// Del delete value
func (instance *memorySession) Del(key string) {
	delete(instance.cache, key)
}

// GetID get session id
func (instance memorySession) GetID() string {
	return instance.id
}

// Expired update expired
func (instance *memorySession) Expired(expired int) {
	var s time.Duration = time.Second * time.Duration(expired)
	instance.expired = time.Now().Add(s).Unix()
}

// IsExpired check whether expired
func (instance memorySession) IsExpired(toExpired int64) bool {
	return instance.expired < toExpired
}

// Save save data
func (instance memorySession) Save() error {
	return nil
}

var _ ISession = &memorySession{}

// newMemory use memory session
func newMemory(id string) *memorySession {
	return &memorySession{
		id:    id,
		cache: make(map[string]interface{}),
	}
}
