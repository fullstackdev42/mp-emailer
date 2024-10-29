package user

import (
	"fmt"

	"github.com/fullstackdev42/mp-emailer/config"
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

// GetOwnerIDFromSession retrieves the owner ID from the session
func GetOwnerIDFromSession(c echo.Context) (string, error) {
	logger, err := getLogger(c)
	if err != nil {
		return "", err
	}
	logger.Debug("GetOwnerIDFromSession: Starting")

	ownerID, ok := c.Get("user_id").(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in session or not a string")
	}

	logger.Debug("GetOwnerIDFromSession: Owner ID retrieved", "ownerID", ownerID)
	return ownerID, nil
}

// AuthMiddleware middleware to set the authenticated flag in the context
func AuthMiddleware(sessionStore sessions.Store, config *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := sessionStore.Get(c.Request(), config.SessionName)
			if err != nil {
				fmt.Printf("Session error: %v\n", err)
			}
			isAuthenticated := sess.Values["authenticated"] == true
			c.Set("IsAuthenticated", isAuthenticated)

			return next(c)
		}
	}
}
