package session

import (
	"time"
)

type memoryMgr struct {
	maxExpired int
	sessions   map[string]ISession
}

func (instance *memoryMgr) Create(key string) ISession {
	if key == "" {
		key = getUUID()
	}

	ss := newMemory(key)

	instance.sessions[key] = ss
	return ss
}
func (instance *memoryMgr) Get(id string) (ISession, bool) {
	v, ok := instance.sessions[id]
	return v, ok
}
func (instance *memoryMgr) GC() {
	now := time.Now().Unix()
	for k, s := range instance.sessions {
		if s.IsExpired(now) {
			delete(instance.sessions, k)
		}
	}

	time.AfterFunc(time.Second*time.Duration(expired), func() {
		instance.GC()
	})
}
func (instance *memoryMgr) SetMaxExpired(expired int) {
	instance.maxExpired = expired
}
func (instance *memoryMgr) GetMaxExpired() int {
	return instance.maxExpired
}

var _ IMgr = &memoryMgr{}

// NewMemoryMgr get memory session mgr
func NewMemoryMgr() IMgr {
	m := &memoryMgr{
		sessions: make(map[string]ISession, 0),
	}
	go m.GC()

	return m
}
