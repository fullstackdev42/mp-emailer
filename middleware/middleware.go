package middleware

import (
	"net/http"

	"github.com/jonesrussell/mp-emailer/config"
	"github.com/jonesrussell/mp-emailer/session"
	"github.com/jonesrussell/mp-emailer/shared"
	"golang.org/x/time/rate"

	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jonesrussell/mp-emailer/logger"
	echojwt "github.com/labstack/echo-jwt/v4"
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
	logger         logger.Interface
	cfg            *config.Config
	errorHandler   shared.ErrorHandlerInterface
	sessionManager session.Manager
}

// ManagerParams for dependency injection
type ManagerParams struct {
	fx.In
	Logger         logger.Interface
	Cfg            *config.Config
	ErrorHandler   shared.ErrorHandlerInterface
	SessionManager session.Manager
}

// NewManager creates a new middleware manager
func NewManager(params ManagerParams) (*Manager, error) {
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
		logger:         params.Logger,
		cfg:            params.Cfg,
		errorHandler:   params.ErrorHandler,
		sessionManager: params.SessionManager,
	}, nil
}

// Register configures all middleware for the application
func (m *Manager) Register(e *echo.Echo) {
	// Add global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
		Getter: middleware.MethodFromForm("_method"),
	}))
	e.Use(m.SessionMiddleware(m.sessionManager))
	e.Use(m.CSRFMiddleware())
	m.registerRateLimiter(e)

	// Don't register JWT globally - it should be applied to specific routes
}

// JWTMiddleware adds JWT authentication middleware
func (m *Manager) JWTMiddleware() echo.MiddlewareFunc {
	config := echojwt.Config{
		SigningKey:    []byte(m.cfg.Auth.JWTSecret),
		SigningMethod: jwt.SigningMethodHS256.Name,
		ErrorHandler: func(c echo.Context, err error) error {
			return m.errorHandler.HandleHTTPError(c, err, "jwt", http.StatusUnauthorized)
		},
	}
	return echojwt.WithConfig(config)
}

// Add this method to create a JWT middleware group
func (m *Manager) JWTProtected(e *echo.Echo) *echo.Group {
	// Create a group for protected routes
	protected := e.Group("")
	protected.Use(m.JWTMiddleware())
	return protected
}

// Add this method to the Manager struct
func (m *Manager) registerRateLimiter(e *echo.Echo) {
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStore(
			rate.Limit(m.cfg.Server.RateLimiting.RequestsPerSecond),
		),
		ErrorHandler: func(c echo.Context, err error) error {
			return m.errorHandler.HandleHTTPError(c, err, "rate_limit", http.StatusTooManyRequests)
		},
	}
	e.Use(middleware.RateLimiterWithConfig(config))
}

// CSRFMiddleware returns CSRF protection middleware
func (m *Manager) CSRFMiddleware() echo.MiddlewareFunc {
	return middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "header:X-CSRF-Token,form:_csrf,form:csrf",
		ContextKey:     "csrf",
		CookieName:     "_csrf",
		CookiePath:     "/",
		CookieMaxAge:   86400,
		CookieSecure:   m.cfg.App.Env == "production",
		CookieHTTPOnly: true,
	})
}
