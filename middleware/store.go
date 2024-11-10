package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type SessionStore interface {
	Get(r *http.Request, name string) (*sessions.Session, error)
	New(r *http.Request, name string) (*sessions.Session, error)
	Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error
	Delete(r *http.Request, w http.ResponseWriter) error
	Options(options *sessions.Options)
	MaxAge(age int)
	Clear(r *http.Request)
	MaxLength(length int)
	SessionID(r *http.Request) string
	Store() sessions.Store
}
