package campaign

import (
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/internal/database"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/go-playground/validator/v10"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// Module defines the campaign module
//
//nolint:gochecknoglobals
var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			NewRepository,
			fx.As(new(RepositoryInterface)),
		),
		fx.Annotate(
			func(repo RepositoryInterface, validate *validator.Validate, logger loggo.LoggerInterface) (ServiceInterface, error) {
				serviceParams := ServiceParams{
					Repo:     repo,
					Validate: validate,
					Logger:   logger,
				}
				serviceResult, err := NewService(serviceParams)
				if err != nil {
					return nil, err
				}
				return NewLoggingServiceDecorator(serviceResult.Service, logger), nil
			},
			fx.As(new(ServiceInterface)),
		),
		fx.Annotate(
			func(logger loggo.LoggerInterface, cfg *config.Config) RepresentativeLookupServiceInterface {
				return NewRepresentativeLookupService(cfg.RepresentativeLookupBaseURL, logger)
			},
			fx.As(new(RepresentativeLookupServiceInterface)),
		),
		fx.Annotate(
			func(params ClientParams) (ClientInterface, error) {
				return NewClient(params)
			},
			fx.As(new(ClientInterface)),
		),
		NewHandler,
	),
)

// RepositoryParams for dependency injection
type RepositoryParams struct {
	fx.In
	DB     database.Interface
	Logger loggo.LoggerInterface
}

// RepositoryResult is the output struct for NewRepository
type RepositoryResult struct {
	fx.Out
	Repository RepositoryInterface `group:"repositories"`
}

// NewRepository creates a new campaign repository
func NewRepository(params RepositoryParams) (RepositoryInterface, error) {
	repo := &Repository{
		db: params.DB,
	}
	return repo, nil
}

// ServiceParams for dependency injection
type ServiceParams struct {
	fx.In
	Repo     RepositoryInterface
	Validate *validator.Validate
	Logger   loggo.LoggerInterface
}

// ServiceResult is the output struct for NewService
type ServiceResult struct {
	fx.Out
	Service ServiceInterface
}

// NewService creates a new campaign service
func NewService(params ServiceParams) (ServiceResult, error) {
	service := ServiceResult{
		Service: &Service{
			repo:     params.Repo,
			validate: params.Validate,
			logger:   params.Logger,
		},
	}
	return service, nil
}

// HandlerParams for dependency injection
type HandlerParams struct {
	fx.In
	Service                     ServiceInterface
	Logger                      loggo.LoggerInterface
	RepresentativeLookupService RepresentativeLookupServiceInterface
	EmailService                email.Service
	Client                      ClientInterface
	ErrorHandler                *shared.ErrorHandler
	TemplateRenderer            *shared.CustomTemplateRenderer
}

// HandlerResult is the output struct for NewHandler
type HandlerResult struct {
	fx.Out
	Handler *Handler
}

// NewHandler initializes a new Handler
func NewHandler(params HandlerParams) (HandlerResult, error) {
	handler := &Handler{
		service:                     params.Service,
		logger:                      params.Logger,
		representativeLookupService: params.RepresentativeLookupService,
		emailService:                params.EmailService,
		client:                      params.Client,
		errorHandler:                params.ErrorHandler,
		templateRenderer:            params.TemplateRenderer,
	}
	return HandlerResult{Handler: handler}, nil
}

// NewRepresentativeLookupService creates a new instance of RepresentativeLookupService
func NewRepresentativeLookupService(baseURL string, logger loggo.LoggerInterface) RepresentativeLookupServiceInterface {
	return &RepresentativeLookupService{
		logger:  logger,
		baseURL: baseURL,
	}
}

// ClientParams for dependency injection
type ClientParams struct {
	fx.In
	Logger        loggo.LoggerInterface
	LookupService RepresentativeLookupServiceInterface
}

// NewClient creates a new instance of ClientInterface
func NewClient(params ClientParams) (ClientInterface, error) {
	client := &DefaultClient{
		logger:        params.Logger,
		lookupService: params.LookupService,
	}
	return client, nil
}

// RegisterRoutes registers the campaign routes
func RegisterRoutes(h *Handler, e *echo.Echo) {
	e.GET("/campaigns", h.GetCampaigns)

	campaignGroup := e.Group("/campaign")
	campaignGroup.POST("", h.CreateCampaign)
	campaignGroup.GET("/:id", h.CampaignGET)
	campaignGroup.PUT("/:id", h.EditCampaign)
	campaignGroup.DELETE("/:id", h.DeleteCampaign)
	campaignGroup.POST("/:id/send", h.SendCampaign)
}
