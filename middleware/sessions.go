package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/mp-emailer/session"
	"github.com/jonesrussell/mp-emailer/shared"
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
	SetSessionValues(sess *sessions.Session, userData interface{})
}

type SessionMiddleware struct {
	store        sessions.Store
	logger       loggo.LoggerInterface
	errorHandler shared.ErrorHandlerInterface
}

// NewSessionManager creates a new session manager implementation
func NewSessionManager(store sessions.Store, logger loggo.LoggerInterface, errorHandler shared.ErrorHandlerInterface) SessionManager {
	return &SessionMiddleware{
		store:        store,
		logger:       logger,
		errorHandler: errorHandler,
	}
}

// GetSession implements SessionManager interface
func (sm *SessionMiddleware) GetSession(c echo.Context, name string) (*sessions.Session, error) {
	sm.logger.Debug("Getting session", "name", name)
	return sm.store.Get(c.Request(), name)
}

// SaveSession implements SessionManager interface
func (sm *SessionMiddleware) SaveSession(c echo.Context, session *sessions.Session) error {
	sm.logger.Debug("Saving session")
	return sm.store.Save(c.Request(), c.Response().Writer, session)
}

// ClearSession implements SessionManager interface
func (sm *SessionMiddleware) ClearSession(c echo.Context, name string) error {
	sm.logger.Debug("Clearing session", "name", name)
	session, err := sm.GetSession(c, name)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return sm.SaveSession(c, session)
}

// ValidateSession implements SessionManager interface
func (sm *SessionMiddleware) ValidateSession(name string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sm.logger.Debug("Validating session", "name", name)

			session, err := sm.GetSession(c, name)
			if err != nil {
				return sm.errorHandler.HandleHTTPError(c, err, "Error getting session", http.StatusInternalServerError)
			}

			// Store session in context for handlers to use
			c.Set("session", session)

			// Store original values to detect changes
			originalValues := make(map[interface{}]interface{})
			for k, v := range session.Values {
				originalValues[k] = v
			}

			if err = next(c); err != nil {
				return err
			}

			// Check if values have changed
			valuesChanged := false
			if len(originalValues) != len(session.Values) {
				valuesChanged = true
			} else {
				for k, v := range session.Values {
					if originalV, exists := originalValues[k]; !exists || originalV != v {
						valuesChanged = true
						break
					}
				}
			}

			// Only save if values changed
			if valuesChanged {
				sm.logger.Debug("Session values changed, saving session")
				if err := sm.SaveSession(c, session); err != nil {
					return sm.errorHandler.HandleHTTPError(c, err, "Error saving session", http.StatusInternalServerError)
				}
			}

			return nil
		}
	}
}
func (sm *SessionMiddleware) SetSessionValues(sess *sessions.Session, userData interface{}) {
	if u, ok := userData.(session.UserData); ok {
		sess.Values["user_id"] = u.GetID()
		sess.Values["username"] = u.GetUsername()
		sess.Values["authenticated"] = true
	}
}
