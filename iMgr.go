package session

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/tonny-zhang/cotton"
)

var expired = 3600

// SessionCookieName session cookie name
var SessionCookieName = "sessionid"

// IMgr mgr interface
type IMgr interface {
	Create(key string) ISession
	Get(id string) (ISession, bool)
	SetMaxExpired(int)
	GetMaxExpired() int
}

func getUUID() string {
	uuid := fmt.Sprintf("%d_%d", time.Now().Nanosecond(), time.Now().UnixNano())
	s := md5.New()
	s.Write([]byte(uuid))
	return hex.EncodeToString(s.Sum(nil))
}

// HasUsedSession check whether used session
func HasUsedSession(ctx *cotton.Context) bool {
	_, ok := ctx.Get(SessionCookieName)
	return ok
}

// GetSession get session
func GetSession(ctx *cotton.Context) (session ISession) {
	v, ok := ctx.Get(SessionCookieName)
	if ok {
		session = v.(ISession)
	} else {
		panic("session no used")
	}
	return
}

// Middleware middleware
func Middleware(mgr IMgr) cotton.HandlerFunc {
	if mgr == nil {
		panic("init first")
	}
	return func(ctx *cotton.Context) {
		sessionID, e := ctx.Cookie(SessionCookieName)

		var ss ISession
		var ok bool
		if e == nil {
			ss, ok = mgr.Get(sessionID)
			if !ok {
				ss = mgr.Create(sessionID)
				sessionID = ss.GetID()
			}
		} else {
			ss = mgr.Create(sessionID)
			sessionID = ss.GetID()
		}
		ss.Expired(mgr.GetMaxExpired())
		ss.Save()

		ctx.Set(SessionCookieName, ss)
		http.SetCookie(ctx.Response, &http.Cookie{
			Name:   SessionCookieName,
			Value:  sessionID,
			Path:   "/",
			MaxAge: expired,
		})
	}
}
