package middleware

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"

	"errors"

	"github.com/google/uuid"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

// Module provides middleware functionality
//
//nolint:gochecknoglobals
var Module = fx.Options(
	fx.Provide(
		NewManager,
	),
)

// Manager handles middleware registration and configuration
type Manager struct {
	sessionStore sessions.Store
	logger       loggo.LoggerInterface
}

// ManagerParams for dependency injection
type ManagerParams struct {
	fx.In
	SessionStore sessions.Store
	Logger       loggo.LoggerInterface
}

// NewManager creates a new middleware manager
func NewManager(params ManagerParams) *Manager {
	return &Manager{
		sessionStore: params.SessionStore,
		logger:       params.Logger,
	}
}

// Register configures all middleware for the application
func (m *Manager) Register(e *echo.Echo) {
	m.registerContextMiddleware(e)
	m.registerLogging(e)
	m.registerRateLimiter(e)
	m.registerMethodOverride(e)
}

// registerContextMiddleware adds session store and logger to context
func (m *Manager) registerContextMiddleware(e *echo.Echo) {
	e.Use(sessionsMiddleware(m.sessionStore))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if m.sessionStore == nil {
				m.logger.Error("Session store not available", fmt.Errorf("session store initialization failed"))
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
			}

			c.Set("logger", m.logger)
			c.Set("store", m.sessionStore)
			return next(c)
		}
	})
}

// sessionsMiddleware sets up session middleware
func sessionsMiddleware(store sessions.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, err := store.Get(c.Request(), "session-name")
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Error retrieving session")
			}
			c.Set("session", session)
			return next(c)
		}
	}
}

// registerLogging configures request logging
func (m *Manager) registerLogging(e *echo.Echo) {
	e.Use(middleware.Logger())
}

// registerRateLimiter configures rate limiting
func (m *Manager) registerRateLimiter(e *echo.Echo) {
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStore(20),
		DenyHandler: func(_ echo.Context, identifier string, err error) error {
			m.logger.Error("Rate limit exceeded", err, "identifier", identifier)
			return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
		},
	}))
}

// MethodOverride middleware converts POST requests with _method parameter to the specified method
func MethodOverride() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == "POST" {
				method := c.FormValue("_method")
				if method != "" {
					c.Request().Method = method
				}
			}
			return next(c)
		}
	}
}

// registerMethodOverride adds method override support
func (m *Manager) registerMethodOverride(e *echo.Echo) {
	e.Use(MethodOverride())
}

// GetSession retrieves the session from the context
func (m *Manager) GetSession(c echo.Context, sessionName string) *sessions.Session {
	m.logger.Debug("Getting session", "session_name", sessionName)

	store := c.Get("store").(sessions.Store)
	session, err := store.Get(c.Request(), sessionName)
	if err != nil {
		m.logger.Debug("Error getting session", "error", err)
		return nil
	}

	fmt.Printf("Session ID: %v, IsNew: %v, Values: %v\n", session.ID, session.IsNew, session.Values) // Debugging statement

	// Ensure the session is saved to persist values
	if err := session.Save(c.Request(), c.Response()); err != nil {
		m.logger.Error("Failed to save session", err)
		return nil
	}

	m.logger.Debug("Session retrieved successfully",
		"session_id", session.ID,
		"is_new", session.IsNew,
		"values_count", len(session.Values))
	return session
}

func (m *Manager) GetUserIDFromSession(c echo.Context, sessionName string) (uuid.UUID, error) {
	m.logger.Debug("Attempting to get userID from session", "session_name", sessionName)

	session := m.GetSession(c, sessionName)
	if session == nil {
		m.logger.Debug("Session is nil")
		return uuid.UUID{}, ErrSessionInvalid
	}

	userIDValue := session.Values["user_id"]
	fmt.Printf("Raw user_id value from session: %v, type: %T\n", userIDValue, userIDValue) // Debugging statement
	m.logger.Debug("Raw user_id value", "type", fmt.Sprintf("%T", userIDValue))

	switch v := userIDValue.(type) {
	case uuid.UUID:
		return v, nil
	case string:
		return uuid.Parse(v)
	case []byte:
		return uuid.Parse(string(v))
	default:
		fmt.Printf("Unexpected type for user_id: %T\n", userIDValue) // Debugging statement
		m.logger.Debug("Unexpected type for user_id", "type", fmt.Sprintf("%T", userIDValue))
		return uuid.UUID{}, ErrUserNotFound
	}
}

// ValidateSession middleware ensures a valid session exists
func (m *Manager) ValidateSession(sessionName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			m.logger.Debug("Validating session",
				"session_name", sessionName,
				"path", c.Path())

			userID, err := m.GetUserIDFromSession(c, sessionName)
			if err != nil {
				m.logger.Debug("Session validation error", "error", err)

				if err == ErrSessionInvalid || err == ErrUserNotFound {
					return c.Redirect(http.StatusSeeOther, "/user/login")
				}
				return err
			}

			m.logger.Debug("Session validated successfully", "user_id", userID)
			return next(c)
		}
	}
}

var (
	ErrSessionInvalid = errors.New("session is invalid")
	ErrUserNotFound   = errors.New("user not found in session")
)
