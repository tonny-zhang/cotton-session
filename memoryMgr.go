package session

import (
	"time"
)

// Mgr session mgr
type memoryMgr struct {
	sessions map[string]ISession
}

func (instance *memoryMgr) Create() ISession {
	key := getUUID()
	ss := NewMemory(key)
	ss.Expired(expired)

	instance.sessions[key] = ss
	return ss
}
func (instance *memoryMgr) Get(id string) (ISession, bool) {
	v, ok := instance.sessions[id]
	return v, ok
}
func (instance *memoryMgr) GC() {
	now := time.Now().Second()
	for k, s := range instance.sessions {
		if s.IsExpired(now) {
			delete(instance.sessions, k)
		}
	}

	time.AfterFunc(time.Second*time.Duration(expired), func() {
		instance.GC()
	})
}

var _ imgr = &memoryMgr{}

// UseMemory use memory session
func UseMemory() {
	m := &memoryMgr{
		sessions: make(map[string]ISession, 0),
	}
	mgrObj = m
	go m.GC()
}
