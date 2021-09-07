package session

// ISession session interface
type ISession interface {
	GetID() string
	Get(key string) (value interface{}, err error)
	Set(key string, value interface{})
	Del(key string)
	Save() error
	Expired(expired int)
	IsExpired(toExpired int64) bool
}
type sessionData struct {
	val     interface{}
	expired int
}
