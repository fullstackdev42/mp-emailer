package user

import (
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

// SetAuthStatusMiddleware sets the isAuthenticated status for all routes
func SetAuthStatusMiddleware(store sessions.Store, logger *loggo.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isAuthenticated := checkAuthentication(c, store, logger)
			c.Set("isAuthenticated", isAuthenticated)
			return next(c)
		}
	}
}

// RequireAuthMiddleware allows or denies access to protected routes
func RequireAuthMiddleware(store sessions.Store, logger *loggo.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isAuthenticated := checkAuthentication(c, store, logger)
			if !isAuthenticated {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}

func checkAuthentication(c echo.Context, store sessions.Store, logger *loggo.Logger) bool {
	sess, err := store.Get(c.Request(), "mpe")
	if err != nil {
		logger.Error("Error getting session", err)
		return false
	}

	username, ok := sess.Values["username"]
	if !ok || username == "" {
		logger.Info("No username found in session")
		return false
	}

	return true
}

// GetOwnerIDFromSession retrieves the owner ID from the session
func GetOwnerIDFromSession(c echo.Context) (int, error) {
	ownerID, ok := c.Get("user_id").(int)
	if !ok {
		return 0, fmt.Errorf("user_id not found in session or not an integer")
	}
	return ownerID, nil
}
