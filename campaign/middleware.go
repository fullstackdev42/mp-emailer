package campaign

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

// getSession retrieves the session from the context
func getSession(c echo.Context, sessionName string, logger loggo.LoggerInterface) *sessions.Session {
	logger.Debug("Getting session",
		"session_name", sessionName)

	store := c.Get("store").(sessions.Store)
	session, err := store.Get(c.Request(), sessionName)
	if err != nil {
		logger.Debug("Error getting session",
			"error", err,
			"session_name", sessionName)
		return nil
	}

	// Log specific session details instead of the entire session object
	logger.Debug("Session retrieved successfully",
		"session_id", session.ID,
		"is_new", session.IsNew,
		"values_count", len(session.Values))
	return session
}

// GetUserIDFromSession safely extracts the user ID from the session
func GetUserIDFromSession(c echo.Context, sessionName string, logger loggo.LoggerInterface) (uuid.UUID, error) {
	logger.Debug("Attempting to get userID from session",
		"session_name", sessionName)

	session := getSession(c, sessionName, logger)
	if session == nil {
		logger.Debug("Session is nil")
		return uuid.UUID{}, ErrSessionInvalid
	}

	// Log session values in a safe way
	var sessionKeys []string
	for k := range session.Values {
		if key, ok := k.(string); ok {
			sessionKeys = append(sessionKeys, key)
		}
	}
	logger.Debug("Session values",
		"available_keys", strings.Join(sessionKeys, ", "))

	userIDValue := session.Values["user_id"]
	logger.Debug("Raw user_id value",
		"type", fmt.Sprintf("%T", userIDValue))

	switch v := userIDValue.(type) {
	case uuid.UUID:
		logger.Debug("Found UUID directly",
			"uuid", v.String())
		return v, nil
	case string:
		parsedUUID, err := uuid.Parse(v)
		if err != nil {
			logger.Debug("Failed to parse UUID string",
				"error", err,
				"value", v)
			return uuid.UUID{}, err
		}
		return parsedUUID, nil
	case []byte:
		userID := string(v)
		parsedUUID, err := uuid.Parse(userID)
		if err != nil {
			logger.Debug("Failed to parse UUID from bytes",
				"error", err,
				"value", userID)
			return uuid.UUID{}, err
		}
		return parsedUUID, nil
	default:
		logger.Debug("Unexpected type for user_id",
			"type", fmt.Sprintf("%T", userIDValue))
		return uuid.UUID{}, ErrUserNotFound
	}
}

// ValidateSessionWithLogging middleware ensures a valid session exists
func ValidateSessionWithLogging(sessionName string, logger loggo.LoggerInterface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Debug("Validating session",
				"session_name", sessionName,
				"path", c.Path())

			userID, err := GetUserIDFromSession(c, sessionName, logger)
			if err != nil {
				logger.Debug("Session validation error",
					"error", err,
					"session_name", sessionName)

				if err == ErrSessionInvalid || err == ErrUserNotFound {
					logger.Debug("Redirecting to login page")
					return c.Redirect(http.StatusSeeOther, "/user/login")
				}
				return err
			}

			logger.Debug("Session validated successfully",
				"user_id", userID.String())

			return next(c)
		}
	}
}
