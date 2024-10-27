package user

import (
	"strconv"
	"time"

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
var Module = fx.Module("user",
	fx.Provide(
		NewRepository,
		NewService,
		NewHandler,
	),
)

// ServiceResult is the output struct for NewService
type ServiceResult struct {
	fx.Out
	Service ServiceInterface
}

// NewService creates a new user service
func NewService(repo RepositoryInterface, logger loggo.LoggerInterface, cfg *config.Config) (ServiceResult, error) {
	expiry, err := strconv.Atoi(cfg.JWTExpiry)
	if err != nil {
		return ServiceResult{}, fmt.Errorf("invalid JWTExpiry: %w", err)
	}

	service := &Service{
		repo:   repo,
		logger: logger,
		config: &Config{
			JWTSecret: cfg.JWTSecret,
			JWTExpiry: time.Duration(expiry) * time.Minute,
		},
	}
	return ServiceResult{Service: service}, nil
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
	service ServiceInterface,
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
