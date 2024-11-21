package middleware

import (
	"net/http"

	"github.com/jonesrussell/mp-emailer/config"
	"github.com/jonesrussell/mp-emailer/shared"
	"golang.org/x/time/rate"

	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jonesrussell/loggo"
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
	logger       loggo.LoggerInterface
	cfg          *config.Config
	errorHandler shared.ErrorHandlerInterface
}

// ManagerParams for dependency injection
type ManagerParams struct {
	fx.In
	Logger       loggo.LoggerInterface
	Cfg          *config.Config
	ErrorHandler shared.ErrorHandlerInterface
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
		logger:       params.Logger,
		cfg:          params.Cfg,
		errorHandler: params.ErrorHandler,
	}, nil
}

// Register configures all middleware for the application
func (m *Manager) Register(e *echo.Echo) {
	// Add global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.MethodOverride())
	m.registerRateLimiter(e)
	// Don't register JWT globally - it should be applied to specific routes
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
