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
			logger.Debug("SetAuthStatusMiddleware: Starting")
			isAuthenticated := checkAuthentication(c, store, logger)
			logger.Debug("SetAuthStatusMiddleware: Authentication check result", "isAuthenticated", isAuthenticated)
			c.Set("isAuthenticated", isAuthenticated)
			logger.Debug("SetAuthStatusMiddleware: Set isAuthenticated in context")
			return next(c)
		}
	}
}

// RequireAuthMiddleware allows or denies access to protected routes
func RequireAuthMiddleware(store sessions.Store, logger *loggo.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Debug("RequireAuthMiddleware: Starting")
			isAuthenticated := checkAuthentication(c, store, logger)
			logger.Debug("RequireAuthMiddleware: Authentication check result", "isAuthenticated", isAuthenticated)
			if !isAuthenticated {
				logger.Warn("RequireAuthMiddleware: Unauthorized access attempt")
				return echo.ErrUnauthorized
			}
			logger.Debug("RequireAuthMiddleware: Access granted")
			return next(c)
		}
	}
}

func checkAuthentication(c echo.Context, store sessions.Store, logger *loggo.Logger) bool {
	logger.Debug("checkAuthentication: Starting")
	sess, err := store.Get(c.Request(), "mpe")
	if err != nil {
		logger.Error("checkAuthentication: Error getting session", err)
		return false
	}
	logger.Debug("checkAuthentication: Session retrieved", "session", fmt.Sprintf("%+v", sess))

	username, ok := sess.Values["username"]
	if !ok {
		logger.Debug("checkAuthentication: No username found in session")
		return false
	}
	logger.Debug("checkAuthentication: Username found in session", "username", username)

	if username == "" {
		logger.Debug("checkAuthentication: Username is empty")
		return false
	}

	logger.Debug("checkAuthentication: Authentication successful")
	return true
}

// GetOwnerIDFromSession retrieves the owner ID from the session
func GetOwnerIDFromSession(c echo.Context) (int, error) {
	logger := c.Get("logger").(*loggo.Logger)
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
