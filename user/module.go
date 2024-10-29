package user

import (
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
		fx.Annotate(
			NewRepository,
			fx.As(new(RepositoryInterface)),
		),
		fx.Annotate(
			func(repo RepositoryInterface, validate *validator.Validate, logger loggo.LoggerInterface, cfg *config.Config) (ServiceInterface, error) {
				serviceParams := ServiceParams{
					Repo:     repo,
					Validate: validate,
					cfg:      cfg,
				}
				serviceResult, err := NewService(serviceParams)
				if err != nil {
					return nil, err
				}
				return NewLoggingServiceDecorator(serviceResult.Service, logger), nil
			},
			fx.As(new(ServiceInterface)),
		),
		NewHandler,
	),
)

// ServiceParams for dependency injection
type ServiceParams struct {
	fx.In
	Repo     RepositoryInterface
	Validate *validator.Validate
	cfg      *config.Config
}

// ServiceResult is the output struct for NewService
type ServiceResult struct {
	fx.Out
	Service ServiceInterface
}

// NewService creates a new user service
func NewService(params ServiceParams) (ServiceResult, error) {
	service := ServiceResult{
		Service: &Service{
			repo:     params.Repo,
			validate: params.Validate,
		},
	}
	return service, nil
}

// RepositoryParams for dependency injection
type RepositoryParams struct {
	fx.In
	DB     *database.DB
	Logger loggo.LoggerInterface
}

// RepositoryResult is the output struct for NewRepository
type RepositoryResult struct {
	fx.Out
	Repository RepositoryInterface `group:"repositories"`
}

// NewRepository creates a new user repository
func NewRepository(params RepositoryParams) (RepositoryInterface, error) {
	repo := &Repository{
		db: params.DB,
	}
	return repo, nil
}

// HandlerParams for dependency injection
type HandlerParams struct {
	fx.In

	Config          *config.Config
	Logger          loggo.LoggerInterface
	Service         ServiceInterface
	Store           sessions.Store
	TemplateManager *shared.CustomTemplateRenderer
	Repo            RepositoryInterface
	ErrorHandler    *shared.ErrorHandler
	FlashHandler    *shared.FlashHandler
}

// HandlerResult is the output struct for NewHandler
type HandlerResult struct {
	fx.Out
	Handler *Handler
}

// NewHandler creates a new user handler
func NewHandler(params HandlerParams) (HandlerResult, error) {
	handler := &Handler{
		service:         params.Service,
		Logger:          params.Logger,
		Store:           params.Store,
		SessionName:     params.Config.SessionName,
		Config:          params.Config,
		templateManager: *params.TemplateManager,
		repo:            params.Repo,
		errorHandler:    params.ErrorHandler,
		flashHandler:    params.FlashHandler,
	}
	return HandlerResult{Handler: handler}, nil
}

// RegisterRoutes registers the user routes
func RegisterRoutes(h *Handler, e *echo.Echo) {
	userGroup := e.Group("/user")
	userGroup.GET("/register", h.RegisterGET)
	userGroup.POST("/register", h.RegisterPOST)
	userGroup.GET("/login", h.LoginGET)
	userGroup.POST("/login", h.LoginPOST)
	userGroup.GET("/logout", h.LogoutGET)
}
