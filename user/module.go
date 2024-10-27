package user

import (
	"fmt" // Import fmt for logging

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/internal/database"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// Module defines the user module
//
//nolint:gochecknoglobals
var Module = fx.Options(
	fx.Provide(
		NewRepository,
		fx.Annotate(func(repo RepositoryInterface, validate *validator.Validate, logger loggo.LoggerInterface, cfg *config.Config) (ServiceInterface, error) {
			fmt.Println("Creating ServiceInterface")
			serviceParams := ServiceParams{Repo: repo, Validate: validate, cfg: cfg}
			serviceResult, err := NewService(serviceParams)
			if err != nil {
				fmt.Println("Error creating ServiceInterface:", err)
				return nil, err
			}
			return NewLoggingServiceDecorator(serviceResult.Service, logger), nil
		}, fx.As(new(ServiceInterface))),
		NewHandler,
	),
)

// NewService creates a new user service
func NewService(params ServiceParams) (ServiceResult, error) {
	fmt.Println("Creating new user service")
	service := ServiceResult{Service: &Service{repo: params.Repo, validate: params.Validate}}
	return service, nil
}

// NewRepository creates a new user repository
func NewRepository(params RepositoryParams) (RepositoryInterface, error) {
	fmt.Println("Creating new repository")
	return &Repository{db: params.DB}, nil
}

// ServiceResult is the output struct for NewService
type ServiceResult struct {
	fx.Out
	Service ServiceInterface
}

// HandlerResult is the output struct for NewHandler
type HandlerResult struct {
	fx.Out
	Handler *Handler
}

// NewHandler creates a new user handler
func NewHandler(cfg *config.Config, logger loggo.LoggerInterface, service ServiceInterface, // This will now receive the decorated service
	sessions sessions.Store, templateManager shared.TemplateRenderer, repo RepositoryInterface, errorHandler *shared.ErrorHandler,
) (HandlerResult, error) {
	fmt.Println("Creating new user handler")
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

// ServiceParams for dependency injection
type ServiceParams struct {
	fx.In
	Repo     RepositoryInterface
	Validate *validator.Validate
	cfg      *config.Config
}

// RepositoryParams for dependency injection
type RepositoryParams struct {
	fx.In
	DB     *database.DB
	Logger loggo.LoggerInterface
}

// RegisterRoutes registers the user routes
func RegisterRoutes(h *Handler, e *echo.Echo) {
	fmt.Println("Registering user routes")
	e.GET("/user/register", h.RegisterGET)
	e.POST("/user/register", h.RegisterPOST)
	e.GET("/user/login", h.LoginGET)
	e.POST("/user/login", h.LoginPOST)
	e.GET("/user/logout", h.LogoutGET)
}
