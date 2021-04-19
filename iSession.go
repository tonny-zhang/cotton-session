package session

// ISession session interface
type ISession interface {
	GetID() string
	Get(key string) (value interface{}, err error)
	Set(key string, value interface{})
	Del(key string)
	Expired(expired int)
	IsExpired(toExpired int) bool
}
type sessionData struct {
	val     interface{}
	expired int
}
