package ui

import (
	"github.com/codr7/unid/lib"
	"github.com/google/uuid"
	"net/http"
	"sync"
)

const (
	SESSION_COOKIE_NAME = "unid-session"
)

var (
	sessions sync.Map
)

type Session struct {
	sessionKey string
	userName string
}

func (self *Session) End() {
	sessions.Delete(self.sessionKey)
}

func StartSession(user *unid.User, w http.ResponseWriter) {
	k := NewSessionKey()
	s := &Session{sessionKey: k, userName: user.Name}
	sessions.Store(k, s)
	http.SetCookie(w, &http.Cookie{Name: SESSION_COOKIE_NAME, Value: k})
}

func NewSessionKey() string {
	return uuid.New().String()
}

func FindSession(r *http.Request) *Session {
	c, err := r.Cookie(SESSION_COOKIE_NAME)

	if err != nil {
		return nil
	}
	
	s, _ := sessions.Load(c.Value)

	if s == nil {
		return nil
	}
	
	return s.(*Session)
}

func CurrentSession(w http.ResponseWriter, r *http.Request) *Session {
	if s := FindSession(r); s != nil {
		return s
	}

	http.Redirect(w, r, "login.html", http.StatusFound)
	return nil
}