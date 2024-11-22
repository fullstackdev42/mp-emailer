package middleware

import (
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/mp-emailer/session"
	"github.com/labstack/echo/v4"
)

// SessionMiddleware injects the session manager into the context
func (m *Manager) SessionMiddleware(sessionManager session.Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Store session manager in context
			c.Set("session_manager", sessionManager)

			// Get or create session
			sess, err := sessionManager.GetSession(c, m.cfg.Auth.SessionName)
			if err != nil {
				m.logger.Error("Failed to get/create session", err)
				return next(c)
			}

			// Get session values based on type
			values := getSessionValues(sess)

			// Debug log session state
			m.logger.Debug("Session state",
				"session_id", getSessionID(sess),
				"is_new", isNewSession(sess),
				"user_id", values["user_id"],
				"is_authenticated", values["is_authenticated"])

			// Store session in context for easy access
			c.Set("session", sess)

			// Add flash messages to context if they exist
			if flashes := sessionManager.GetFlashes(sess); len(flashes) > 0 {
				c.Set("flashes", flashes)
			}

			// Call next handler
			err = next(c)

			// Debug log session state after handler
			m.logger.Debug("Session state after handler",
				"session_id", getSessionID(sess),
				"user_id", values["user_id"],
				"is_authenticated", values["is_authenticated"])

			// Save session after processing request
			if saveErr := sessionManager.SaveSession(c, sess); saveErr != nil {
				m.logger.Error("Failed to save session", saveErr)
				return saveErr
			}

			return err
		}
	}
}

// Helper functions to handle different session types
func getSessionValues(sess interface{}) map[interface{}]interface{} {
	switch s := sess.(type) {
	case session.Interface:
		return s.Values()
	case *sessions.Session:
		return s.Values
	default:
		return make(map[interface{}]interface{})
	}
}

func getSessionID(sess interface{}) string {
	switch s := sess.(type) {
	case session.Interface:
		return s.GetID()
	case *sessions.Session:
		return s.ID
	default:
		return ""
	}
}

func isNewSession(sess interface{}) bool {
	switch s := sess.(type) {
	case session.Interface:
		return s.IsNew()
	case *sessions.Session:
		return s.IsNew
	default:
		return false
	}
}
