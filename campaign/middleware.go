package campaign

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

// getSession retrieves the session from the context
func getSession(c echo.Context, sessionName string) *sessions.Session {
	store := c.Get("store").(sessions.Store)
	session, _ := store.Get(c.Request(), sessionName)
	return session
}

// GetUserIDFromSession safely extracts the user ID from the session
func GetUserIDFromSession(c echo.Context, sessionName string) (string, error) {
	session := getSession(c, sessionName)
	if session == nil {
		return "", ErrSessionInvalid
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		return "", ErrUserNotFound
	}

	return userID, nil
}

// ValidateSession middleware ensures a valid session exists
func ValidateSession(sessionName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, err := GetUserIDFromSession(c, sessionName)
			if err != nil {
				if err == ErrSessionInvalid || err == ErrUserNotFound {
					return c.Redirect(http.StatusSeeOther, "/user/login")
				}
				return err
			}
			return next(c)
		}
	}
}
