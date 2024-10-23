package user

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

// getLogger retrieves the logger from the context
func getLogger(c echo.Context) (loggo.LoggerInterface, error) {
	logger, ok := c.Get("logger").(loggo.LoggerInterface)
	if !ok {
		return nil, fmt.Errorf("logger not found in context")
	}
	return logger, nil
}

// getSession retrieves the session from the request
func getSession(c echo.Context, store sessions.Store, sessionName string) (*sessions.Session, error) {
	return store.Get(c.Request(), sessionName)
}

// isAuthenticated checks if the user is authenticated
func isAuthenticated(sess *sessions.Session) bool {
	authenticated, ok := sess.Values["authenticated"].(bool)
	return ok && authenticated
}

// SetAuthStatusMiddleware sets the isAuthenticated status for all routes
func SetAuthStatusMiddleware(store sessions.Store, logger loggo.LoggerInterface, sessionName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Debug("SetAuthStatusMiddleware: Starting")

			c.Set("logger", logger)

			sess, err := getSession(c, store, sessionName)
			if err != nil {
				logger.Error("SetAuthStatusMiddleware: Error getting session", err)
				c.Set("isAuthenticated", false)
			} else {
				authenticated, ok := sess.Values["authenticated"].(bool)
				isAuthenticated := ok && authenticated && sess.Values["username"] != ""
				logger.Debug("SetAuthStatusMiddleware: Authentication check result", "isAuthenticated", isAuthenticated)
				c.Set("isAuthenticated", isAuthenticated)
				c.Set("username", sess.Values["username"])
				c.Set("user_id", sess.Values["user_id"])
			}

			logger.Debug("SetAuthStatusMiddleware: Set isAuthenticated in context")
			return next(c)
		}
	}
}

// RequireAuthMiddleware allows or denies access to protected routes
func RequireAuthMiddleware(store sessions.Store, sessionName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger, err := getLogger(c)
			if err != nil {
				return err
			}

			logger.Debug("RequireAuthMiddleware: Starting")

			sess, err := getSession(c, store, sessionName)
			if err != nil {
				logger.Error("RequireAuthMiddleware: Error getting session", err)
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			if !isAuthenticated(sess) {
				logger.Warn("RequireAuthMiddleware: Unauthorized access attempt")
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			logger.Debug("RequireAuthMiddleware: Access granted")
			return next(c)
		}
	}
}

// GetOwnerIDFromSession retrieves the owner ID from the session
func GetOwnerIDFromSession(c echo.Context) (string, error) {
	logger, err := getLogger(c)
	if err != nil {
		return "", err
	}
	logger.Debug("GetOwnerIDFromSession: Starting")

	ownerID, ok := c.Get("user_id").(string)
	if !ok {
		err := fmt.Errorf("user_id not found in session or not a string")
		logger.Error("GetOwnerIDFromSession: %v", err)
		return "", err
	}

	logger.Debug("GetOwnerIDFromSession: Owner ID retrieved", "ownerID", ownerID)
	return ownerID, nil
}
