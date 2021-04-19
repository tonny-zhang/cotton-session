package session

import (
	"fmt"
	"time"
)

// MemorySession memory session
type MemorySession struct {
	id      string
	expired int
	cache   map[string]interface{}
}

// Get get value
func (instance MemorySession) Get(key string) (value interface{}, err error) {
	v, ok := instance.cache[key]
	if !ok {
		err = fmt.Errorf("no value")
	} else {
		value = v
	}
	return
}

// Set set value
func (instance MemorySession) Set(key string, value interface{}) {
	instance.cache[key] = value
}

// Del delete value
func (instance MemorySession) Del(key string) {
	delete(instance.cache, key)
}

// GetID get session id
func (instance MemorySession) GetID() string {
	return instance.id
}

// Expired update expired
func (instance MemorySession) Expired(expired int) {
	var s time.Duration = time.Second * time.Duration(expired)
	instance.expired = time.Now().Add(s).Second()
}

// IsExpired check whether expired
func (instance MemorySession) IsExpired(toExpired int) bool {
	return instance.expired < toExpired
}

var _ ISession = MemorySession{}

// NewMemory use memory session
func NewMemory(id string) *MemorySession {
	return &MemorySession{
		id:    id,
		cache: make(map[string]interface{}),
	}
}
