package ui

import (
	"fmt"
	"github.com/codr7/unid/lib"
	"github.com/codr7/unid/lib/db"
	"github.com/google/uuid"
	"net/http"
	"net/url"
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
	user *unid.User
}

func (self *Session) Cx() *db.Cx {
	return self.user.Cx()
}

func (self *Session) TimeFormat() string {
	return "2006-01-02 15:04"
}

func (self *Session) DateFormat() string {
	return "2006-01-02"
}

func (self *Session) End() {
	sessions.Delete(self.sessionKey)
}

func StartSession(user *unid.User, w http.ResponseWriter) *Session {
	k := NewSessionKey()
	s := &Session{sessionKey: k, user: user}
	sessions.Store(k, s)
	http.SetCookie(w, &http.Cookie{Name: SESSION_COOKIE_NAME, Value: k, SameSite: http.SameSiteLaxMode})
	return s
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

	http.Redirect(w, r, fmt.Sprintf("login.html?href=%v", url.QueryEscape(r.URL.Path)), http.StatusFound)
	return nil
}
