package session

import (
	"encoding/json"

	"github.com/go-redis/redis"
)

type redisMgr struct {
	maxExpired int
	sessions   map[string]ISession
	client     *redis.Client
}

func (instance *redisMgr) Create(key string) ISession {
	if key == "" {
		key = getUUID()
	}
	ss := newRedis(key, instance.client)
	instance.sessions[key] = ss

	ss.Save()
	return ss
}
func (instance *redisMgr) Get(id string) (ISession, bool) {
	s, e := instance.client.Get(id).Result()
	if e != nil {
		return nil, false
	}
	ss := newRedis(id, instance.client)

	ee := json.Unmarshal([]byte(s), &ss.cache)
	if ee != nil {
		return nil, false
	}
	return ss, true
}
func (instance *redisMgr) GetMaxExpired() int {
	return instance.maxExpired
}
func (instance *redisMgr) SetMaxExpired(expired int) {
	instance.maxExpired = expired
}

var _ IMgr = &redisMgr{}

// NewRedisMgr get redis session mgr
func NewRedisMgr(addr string, db int, passwd string) IMgr {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       db,
		Password: passwd,
	})
	_, e := redisClient.Ping().Result()
	if e != nil {
		panic(e)
	}
	m := &redisMgr{
		client:   redisClient,
		sessions: make(map[string]ISession, 0),
	}

	return m
}
