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

// RequireAuthMiddleware middleware to require authentication
func (h *Handler) RequireAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger, err := getLogger(c)
		if err != nil {
			return err
		}

		sess, err := h.Store.Get(c.Request(), h.SessionName)
		if err != nil {
			logger.Error("RequireAuthMiddleware: Failed to get session", err)
			return c.Redirect(http.StatusSeeOther, "/user/login")
		}

		if auth, ok := sess.Values["authenticated"].(bool); !ok || !auth {
			logger.Debug("RequireAuthMiddleware: User not authenticated")
			return c.Redirect(http.StatusSeeOther, "/user/login")
		}

		// User is authenticated, set user ID in context
		userID, ok := sess.Values["user_id"].(string)
		if !ok {
			logger.Info("RequireAuthMiddleware: User ID not found in session")
			return c.Redirect(http.StatusSeeOther, "/user/login")
		}

		c.Set("user_id", userID)
		logger.Debug("RequireAuthMiddleware: User authenticated", "userID", userID)

		return next(c)
	}
}

// GetAuthenticatedUser retrieves the authenticated user from the context
func GetAuthenticatedUser(c echo.Context) *User {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return nil
	}

	// Here you would typically fetch the user from your database
	// For this example, we'll just return a simple User struct
	return &User{ID: userID}
}

func AuthMiddleware(sessionStore sessions.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := sessionStore.Get(c.Request(), "session")
			if auth, ok := sess.Values["authenticated"].(bool); ok && auth {
				c.Set("IsAuthenticated", true)
			} else {
				c.Set("IsAuthenticated", false)
			}
			return next(c)
		}
	}
}
