package user

import (
	"fmt"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// Module defines the user module
//
//nolint:gochecknoglobals
var Module = fx.Module(
	"user",
	fx.Provide(
		NewRepository,
		// Provide the base service with a different name to avoid confusion
		fx.Annotated{Name: "base_service", Target: NewService},
		// Provide the decorated service as the main ServiceInterface
		fx.Annotated{Target: func(base ServiceInterface, logger loggo.LoggerInterface) (ServiceInterface, error) {
			decorator := NewLoggingServiceDecorator(base, logger)
			return decorator, nil
		}},
		NewHandler,
	),
)

// NewService creates a new user service
func NewService(repo RepositoryInterface, logger loggo.LoggerInterface, cfg *config.Config) (*Service, error) {
	if repo == nil {
		return nil, fmt.Errorf("repository cannot be nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	return &Service{
		repo:   repo,
		logger: logger,
	}, nil
}

// HandlerResult is the output struct for NewHandler
type HandlerResult struct {
	fx.Out
	Handler *Handler
}

// NewHandler creates a new user handler
func NewHandler(
	cfg *config.Config,
	logger loggo.LoggerInterface,
	service ServiceInterface, // This will now receive the decorated service
	sessions sessions.Store,
	templateManager shared.TemplateRenderer,
	repo RepositoryInterface,
	errorHandler *shared.ErrorHandler,
) (HandlerResult, error) {
	handler := &Handler{
		service:         service,
		Logger:          logger,
		Store:           sessions,
		SessionName:     cfg.SessionName,
		Config:          cfg,
		templateManager: templateManager,
		repo:            repo,
		errorHandler:    errorHandler,
	}
	return HandlerResult{Handler: handler}, nil
}

// RegisterRoutes registers the user routes
func RegisterRoutes(h *Handler, e *echo.Echo) {
	e.GET("/user/register", h.RegisterGET)
	e.POST("/user/register", h.RegisterPOST)
	e.GET("/user/login", h.LoginGET)
	e.POST("/user/login", h.LoginPOST)
	e.GET("/user/logout", h.LogoutGET)
}
