package user

import (
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type SessionManager interface {
	GetSession(c echo.Context) (*sessions.Session, error)
	SaveSession(c echo.Context, sess *sessions.Session) error
	SetSessionValues(sess *sessions.Session, user *User)
	ClearSession(sess *sessions.Session)
}

type defaultSessionManager struct {
	store  sessions.Store
	config *config.Config
}

func NewSessionManager(store sessions.Store, cfg *config.Config) SessionManager {
	return &defaultSessionManager{
		store:  store,
		config: cfg,
	}
}

func (sm *defaultSessionManager) GetSession(c echo.Context) (*sessions.Session, error) {
	return sm.store.Get(c.Request(), sm.config.Auth.SessionName)
}

func (sm *defaultSessionManager) SaveSession(c echo.Context, sess *sessions.Session) error {
	return sess.Save(c.Request(), c.Response().Writer)
}

func (sm *defaultSessionManager) SetSessionValues(sess *sessions.Session, user *User) {
	sess.Values["user_id"] = user.ID
	sess.Values["username"] = user.Username
	sess.Values["authenticated"] = true
}

func (sm *defaultSessionManager) ClearSession(sess *sessions.Session) {
	sess.Values = make(map[interface{}]interface{})
	sess.Options.MaxAge = -1
}
