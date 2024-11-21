package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/mp-emailer/config"
	"github.com/jonesrussell/mp-emailer/shared"

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
	cfg          *config.Config
	errorHandler shared.ErrorHandlerInterface
}

// ManagerParams for dependency injection
type ManagerParams struct {
	fx.In
	SessionStore sessions.Store
	Logger       loggo.LoggerInterface
	Cfg          *config.Config
	ErrorHandler shared.ErrorHandlerInterface
}

// NewManager creates a new middleware manager
func NewManager(params ManagerParams) (*Manager, error) {
	if params.SessionStore == nil {
		return nil, errors.New("session store cannot be nil")
	}
	if params.Logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if params.Cfg == nil {
		return nil, errors.New("config cannot be nil")
	}
	if params.ErrorHandler == nil {
		return nil, errors.New("error handler cannot be nil")
	}

	return &Manager{
		sessionStore: params.SessionStore,
		logger:       params.Logger,
		cfg:          params.Cfg,
		errorHandler: params.ErrorHandler,
	}, nil
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
	e.Use(m.SessionsMiddleware())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if m.sessionStore == nil {
				m.logger.Error("Session store not available", fmt.Errorf("session store initialization failed"))
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
			}

			c.Set("logger", m.logger)
			c.Set("store", m.sessionStore)
			c.Set("middleware_manager", m)
			return next(c)
		}
	})
}

// SessionsMiddleware sets up session middleware
func (m *Manager) SessionsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, err := m.sessionStore.Get(c.Request(), m.cfg.Auth.SessionName)
			if err != nil {
				m.logger.Error("Session error", err)
				// Continue with a new session instead of returning an error
				session, _ = m.sessionStore.New(c.Request(), m.cfg.Auth.SessionName)
			}

			// Store the session in context for later use
			c.Set("session", session)

			// Process the request
			err = next(c)
			if err != nil {
				return err
			}

			// Save the session after processing
			if err := session.Save(c.Request(), c.Response().Writer); err != nil {
				m.logger.Error("Failed to save session", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save session")
			}

			return nil
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
					c.Logger().Debug("Method override",
						"original_method", c.Request().Method,
						"new_method", method,
						"path", c.Request().URL.Path)
					c.Request().Method = method
				}
			}
			return next(c)
		}
	}
}

// registerMethodOverride adds method override support
func (m *Manager) registerMethodOverride(e *echo.Echo) {
	e.Pre(MethodOverride())
}

// GetSession retrieves the session from the context
func (m *Manager) GetSession(c echo.Context, sessionName string) *sessions.Session {
	m.logger.Debug("Getting session", "session_name", sessionName)

	// First try to get from context
	if sess, ok := c.Get("session").(*sessions.Session); ok {
		return sess
	}

	// Fallback to getting from store
	session, err := m.sessionStore.Get(c.Request(), sessionName)
	if err != nil {
		m.logger.Debug("Error getting session", "error", err)
		return nil
	}

	m.logger.Debug("Session retrieved successfully",
		"session_id", session.ID,
		"is_new", session.IsNew,
		"values_count", len(session.Values))
	return session
}

// GetUserIDFromSession retrieves the user ID from the session with proper type handling
func (m *Manager) GetUserIDFromSession(c echo.Context, sessionName string) (string, error) {
	m.logger.Debug("Attempting to get userID from session",
		"session_name", sessionName)

	session := m.GetSession(c, sessionName)
	if session == nil {
		m.logger.Debug("Session is nil",
			"session_name", sessionName)
		return "", ErrSessionInvalid
	}

	userIDValue := session.Values["user_id"]
	if userIDValue == nil {
		m.logger.Debug("User ID not found in session",
			"session_name", sessionName)
		return "", ErrUserNotFound
	}

	// Handle different types of user ID storage
	switch v := userIDValue.(type) {
	case string:
		m.logger.Debug("User ID found in session",
			"user_id", v,
			"session_name", sessionName)
		return v, nil
	case uuid.UUID:
		m.logger.Debug("User ID (UUID) found in session",
			"user_id", v.String(),
			"session_name", sessionName)
		return v.String(), nil
	case []byte:
		m.logger.Debug("User ID (bytes) found in session",
			"user_id", string(v),
			"session_name", sessionName)
		return string(v), nil
	default:
		m.logger.Debug("Invalid user ID type in session",
			"type", fmt.Sprintf("%T", userIDValue),
			"session_name", sessionName)
		return "", ErrUserNotFound
	}
}

// ValidateSession middleware ensures a valid session exists
func (m *Manager) ValidateSession(sessionName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, err := m.sessionStore.Get(c.Request(), sessionName)
			if err != nil {
				m.logger.Error("Session validation error", err)
				return c.Redirect(http.StatusSeeOther, "/user/login")
			}

			if session.IsNew {
				m.logger.Debug("New session detected, redirecting to login")
				return c.Redirect(http.StatusSeeOther, "/user/login")
			}

			// Check authentication status
			auth, ok := session.Values["authenticated"].(bool)
			if !ok || !auth {
				m.logger.Debug("User not authenticated")
				return c.Redirect(http.StatusSeeOther, "/user/login")
			}

			// Get and validate user ID
			userID, err := m.GetUserIDFromSession(c, sessionName)
			if err != nil {
				m.logger.Debug("Invalid user ID in session", "error", err)
				return c.Redirect(http.StatusSeeOther, "/user/login")
			}

			c.Set("user_id", userID)
			return next(c)
		}
	}
}

var (
	ErrSessionInvalid = errors.New("session is invalid")
	ErrUserNotFound   = errors.New("user not found in session")
)

// GetOwnerIDFromSession retrieves the owner ID from the session
func (m *Manager) GetOwnerIDFromSession(c echo.Context) (string, error) {
	m.logger.Debug("GetOwnerIDFromSession: Starting")

	// Use the session name from config
	ownerID, err := m.GetUserIDFromSession(c, m.cfg.Auth.SessionName)
	if err != nil {
		m.logger.Debug("GetOwnerIDFromSession: Failed to get owner ID", "error", err)
		return "", err
	}

	m.logger.Debug("GetOwnerIDFromSession: Owner ID retrieved", "ownerID", ownerID)
	return ownerID, nil
}

// JWTMiddleware creates a middleware function for JWT authentication
func (m *Manager) JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": "Missing authorization header",
				})
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization header",
				})
			}

			token := tokenParts[1]
			claims, err := shared.ValidateToken(token, m.cfg.Auth.JWTSecret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token",
				})
			}

			c.Set("user", claims)
			c.Set("user_id", claims.Username)
			return next(c)
		}
	}
}

// IsAuthenticated checks if the user is authenticated in the current session
func (m *Manager) IsAuthenticated(c echo.Context) bool {
	session, err := m.sessionStore.Get(c.Request(), m.cfg.Auth.SessionName)
	if err != nil {
		m.logger.Debug("Failed to get session for authentication check", "error", err)
		return false
	}

	// Get authentication value from session
	authValue, exists := session.Values["authenticated"]
	if !exists {
		m.logger.Debug("No authentication value found in session")
		return false
	}

	// Type assert to boolean
	auth, ok := authValue.(bool)
	if !ok {
		m.logger.Debug("Invalid authentication value type in session")
		return false
	}

	m.logger.Debug("Authentication status checked", "isAuthenticated", auth)
	return auth
}
