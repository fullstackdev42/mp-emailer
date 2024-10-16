package user

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

// SetAuthStatusMiddleware sets the isAuthenticated status for all routes
func SetAuthStatusMiddleware(store sessions.Store, logger loggo.LoggerInterface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Debug("SetAuthStatusMiddleware: Starting")

			// Set the logger in the context
			c.Set("logger", logger)

			isAuthenticated := checkAuthentication(c, store, logger)
			logger.Debug("SetAuthStatusMiddleware: Authentication check result", "isAuthenticated", isAuthenticated)
			c.Set("isAuthenticated", isAuthenticated)
			logger.Debug("SetAuthStatusMiddleware: Set isAuthenticated in context")
			return next(c)
		}
	}
}

// RequireAuthMiddleware allows or denies access to protected routes
func RequireAuthMiddleware(store sessions.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger, ok := c.Get("logger").(loggo.LoggerInterface)
			if !ok {
				return fmt.Errorf("logger not found in context")
			}

			logger.Debug("RequireAuthMiddleware: Starting")

			sess, err := store.Get(c.Request(), "mpe")
			if err != nil {
				logger.Error("RequireAuthMiddleware: Error getting session", err)
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			authenticated, ok := sess.Values["authenticated"].(bool)
			if !ok || !authenticated {
				logger.Warn("RequireAuthMiddleware: Unauthorized access attempt")
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			logger.Debug("RequireAuthMiddleware: Access granted")
			return next(c)
		}
	}
}

func checkAuthentication(c echo.Context, store sessions.Store, logger loggo.LoggerInterface) bool {
	logger.Debug("checkAuthentication: Starting")
	sess, err := store.Get(c.Request(), "mpe")
	if err != nil {
		logger.Error("checkAuthentication: Error getting session", err)
		return false
	}
	logger.Debug("checkAuthentication: Session retrieved", "session", fmt.Sprintf("%+v", sess))

	authenticated, ok := sess.Values["authenticated"].(bool)
	if !ok || !authenticated {
		logger.Debug("checkAuthentication: User not authenticated")
		return false
	}

	username, ok := sess.Values["username"].(string)
	if !ok || username == "" {
		logger.Debug("checkAuthentication: No valid username found in session")
		return false
	}

	logger.Debug("checkAuthentication: Authentication successful", "username", username)
	return true
}

// GetOwnerIDFromSession retrieves the owner ID from the session
func GetOwnerIDFromSession(c echo.Context) (int, error) {
	logger, ok := c.Get("logger").(loggo.LoggerInterface)
	if !ok {
		// If logger is not set return an error
		return 0, fmt.Errorf("logger not found in context")
	}
	logger.Debug("GetOwnerIDFromSession: Starting")

	ownerID, ok := c.Get("user_id").(int)
	if !ok {
		err := fmt.Errorf("user_id not found in session or not an integer")
		logger.Error("GetOwnerIDFromSession: %v", err)
		return 0, err
	}

	logger.Debug("GetOwnerIDFromSession: Owner ID retrieved", "ownerID", ownerID)
	return ownerID, nil
}
