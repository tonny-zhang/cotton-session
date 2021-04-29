package session

import (
	"encoding/json"

	"github.com/bradfitz/gomemcache/memcache"
)

type memcacheMgr struct {
	maxExpired int
	client     *memcache.Client
	sessions   map[string]ISession
}

func (instance *memcacheMgr) Create(key string) ISession {
	if key == "" {
		key = getUUID()
	}
	ss := newMemcache(key, instance.client)
	instance.sessions[key] = ss

	ss.Save()
	return ss
}
func (instance *memcacheMgr) Get(id string) (ISession, bool) {
	s, e := instance.client.Get(id)
	if e != nil {
		return nil, false
	}
	ss := newMemcache(id, instance.client)

	ee := json.Unmarshal(s.Value, &ss.cache)
	if ee != nil {
		return nil, false
	}
	return ss, true
}
func (instance *memcacheMgr) GetMaxExpired() int {
	return instance.maxExpired
}
func (instance *memcacheMgr) SetMaxExpired(expired int) {
	instance.maxExpired = expired
}

var _ IMgr = &redisMgr{}

// NewMemcacheMgr get memcache session mgr
func NewMemcacheMgr(addr string) IMgr {
	client := memcache.New(addr)

	m := &memcacheMgr{
		client:   client,
		sessions: make(map[string]ISession, 0),
	}

	return m
}
