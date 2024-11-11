package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
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

type SessionManager interface {
	GetSession(c echo.Context, name string) (*sessions.Session, error)
	SaveSession(c echo.Context, session *sessions.Session) error
	ClearSession(c echo.Context, name string) error
	ValidateSession(name string) echo.MiddlewareFunc
}
