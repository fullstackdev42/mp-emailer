package campaign

import (
	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/internal/database"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/go-playground/validator/v10"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// Module defines the campaign module
var Module = fx.Options(
	fx.Provide(
		NewRepository,
		NewClient,
		fx.Annotate(NewRepresentativeLookupService,
			fx.ParamTags(`name:"representativeLookupBaseURL"`, `name:"representativeLogger"`),
		),
		fx.Annotate(func(repo RepositoryInterface, validate *validator.Validate, logger loggo.LoggerInterface) (ServiceInterface, error) {
			serviceParams := ServiceParams{Repo: repo, Validate: validate}
			serviceResult, err := NewService(serviceParams)
			if err != nil {
				return nil, err
			}
			return NewCampaignLoggingServiceDecorator(serviceResult.Service, logger), nil
		}, fx.As(new(ServiceInterface))),
		NewHandler,
	),
)

// NewRepository creates a new campaign repository
func NewRepository(params RepositoryParams) (RepositoryInterface, error) {
	return &Repository{db: params.DB}, nil
}

// NewService creates a new campaign service
func NewService(params ServiceParams) (ServiceResult, error) {
	service := ServiceResult{
		Service: &Service{repo: params.Repo, validate: params.Validate},
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
	TemplateRenderer            shared.TemplateRenderer
}

// NewHandler initializes a new Handler
func NewHandler(params HandlerParams) *Handler {
	return &Handler{
		service:                     params.Service,
		logger:                      params.Logger,
		representativeLookupService: params.RepresentativeLookupService,
		emailService:                params.EmailService,
		client:                      params.Client,
		errorHandler:                params.ErrorHandler,
		templateRenderer:            params.TemplateRenderer,
	}
}

// Define the missing structs
type ServiceParams struct {
	fx.In
	Repo     RepositoryInterface
	Validate *validator.Validate
}

type RepositoryParams struct {
	fx.In
	DB     *database.DB
	Logger loggo.LoggerInterface
}

type ClientParams struct {
	fx.In
	Logger        loggo.LoggerInterface
	LookupService RepresentativeLookupServiceInterface
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

// RegisterRoutes registers the campaign routes
func RegisterRoutes(h *Handler, e *echo.Echo) {
	e.GET("/campaign", h.GetAllCampaigns)
	e.POST("/campaign", h.CreateCampaign)
	e.GET("/campaign/:id", h.CampaignGET)
	e.PUT("/campaign/:id", h.EditCampaign)
	e.DELETE("/campaign/:id", h.DeleteCampaign)
	e.POST("/campaign/:id/send", h.SendCampaign)
}

// NewRepresentativeLookupService creates a new instance of RepresentativeLookupService
func NewRepresentativeLookupService(baseURL string, logger loggo.LoggerInterface) RepresentativeLookupServiceInterface {
	return &RepresentativeLookupService{
		logger:  logger,
		baseURL: baseURL,
	}
}

// NewClient creates a new instance of ClientInterface
func NewClient(params ClientParams) (ClientInterface, error) {
	client := &DefaultClient{
		logger:        params.Logger,
		lookupService: params.LookupService,
	}
	return client, nil
}
