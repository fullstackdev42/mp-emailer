package session

import (
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

type manager struct {
	store  sessions.Store
	logger loggo.LoggerInterface
}

func NewManager(store sessions.Store, logger loggo.LoggerInterface) Manager {
	return &manager{
		store:  store,
		logger: logger,
	}
}

func (m *manager) GetSession(c echo.Context, name string) (*sessions.Session, error) {
	return m.store.Get(c.Request(), name)
}

func (m *manager) SaveSession(c echo.Context, session *sessions.Session) error {
	return session.Save(c.Request(), c.Response().Writer)
}

func (m *manager) ClearSession(c echo.Context, name string) error {
	session, err := m.GetSession(c, name)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return m.SaveSession(c, session)
}

func (m *manager) ValidateSession(name string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, err := m.GetSession(c, name)
			if err != nil {
				return err
			}
			if session.IsNew {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}

func (m *manager) SetSessionValues(sess *sessions.Session, userData interface{}) {
	if userData == nil {
		return
	}
	if ud, ok := userData.(UserData); ok {
		sess.Values["user_id"] = ud.GetID()
		sess.Values["username"] = ud.GetUsername()
	}
}
