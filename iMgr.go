package session

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/tonny-zhang/cotton"
)

var mgrObj imgr
var expired = 3600

// SessionName session cookie name
var SessionName = "sessionid"

type imgr interface {
	Create() ISession
	Get(id string) (ISession, bool)
}

func getUUID() string {
	uuid := fmt.Sprintf("%d_%d", time.Now().Nanosecond(), time.Now().UnixNano())
	s := md5.New()
	s.Write([]byte(uuid))
	return hex.EncodeToString(s.Sum(nil))
}

// HasUsedSession check whether used session
func HasUsedSession(ctx *cotton.Context) bool {
	_, ok := ctx.Get(SessionName)
	return ok
}

// GetSession get session
func GetSession(ctx *cotton.Context) (session ISession) {
	v, ok := ctx.Get(SessionName)
	if ok {
		session = v.(ISession)
	} else {
		panic("session no used")
	}
	return
}

// Middleware middleware
func Middleware(name string, options ...string) cotton.HandlerFunc {
	switch name {
	default:
		UseMemory()
	}
	if mgrObj == nil {
		panic("init first")
	}
	return func(ctx *cotton.Context) {
		sessionID, e := ctx.Cookie(SessionName)

		var ss ISession
		var ok bool
		if e == nil {
			ss, ok = mgrObj.Get(sessionID)
			if !ok {
				ss = mgrObj.Create()
				sessionID = ss.GetID()
			}
		} else {
			ss = mgrObj.Create()
			sessionID = ss.GetID()
		}
		ctx.Set(SessionName, ss)
		http.SetCookie(ctx.Response, &http.Cookie{
			Name:   SessionName,
			Value:  sessionID,
			Path:   "/",
			MaxAge: expired,
		})
	}
}
