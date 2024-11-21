package session

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type Manager interface {
	GetSession(c echo.Context, name string) (*sessions.Session, error)
	SaveSession(c echo.Context, session *sessions.Session) error
	ClearSession(c echo.Context, name string) error
	ValidateSession(name string) echo.MiddlewareFunc
	SetSessionValues(sess *sessions.Session, userData interface{})
}

// UserData represents the minimal user data needed for sessions
type UserData interface {
	GetID() interface{}
	GetUsername() string
}
