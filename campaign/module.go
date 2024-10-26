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
//
//nolint:gochecknoglobals
var Module = fx.Module("campaign",
	fx.Provide(
		NewRepository,
		NewService,
		NewHandler,
		NewRepresentativeLookupService,
		NewClient,
		fx.Annotated{
			Name: "representativeLookupBaseURL",
			Target: func() string {
				return "https://represent.opennorth.ca"
			},
		},
		validator.New,
	),
)

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

// NewRepository creates a new campaign repository
type RepositoryParams struct {
	fx.In

	DB     *database.DB
	Logger loggo.LoggerInterface
}

func NewRepository(params RepositoryParams) (RepositoryInterface, error) {
	return &Repository{
		db:     params.DB,
		logger: params.Logger,
	}, nil
}

// NewService creates a new campaign service
type ServiceParams struct {
	fx.In

	Repo     RepositoryInterface
	Validate *validator.Validate
}

func NewService(params ServiceParams) (ServiceResult, error) {
	service := ServiceResult{
		Service: &Service{
			repo:     params.Repo,
			validate: params.Validate,
		},
	}
	return service, nil
}

// NewHandler initializes a new Handler
type HandlerParams struct {
	fx.In

	Service                     ServiceInterface
	Logger                      loggo.LoggerInterface
	RepresentativeLookupService RepresentativeLookupServiceInterface
	EmailService                email.Service
	Client                      ClientInterface
}

func NewHandler(params HandlerParams) *Handler {
	return &Handler{
		service:                     params.Service,
		logger:                      params.Logger,
		representativeLookupService: params.RepresentativeLookupService,
		emailService:                params.EmailService,
		client:                      params.Client,
		errorHandler:                shared.NewErrorHandler(params.Logger),
	}
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
type RepresentativeLookupServiceParams struct {
	fx.In

	Logger  loggo.LoggerInterface
	BaseURL string `name:"representativeLookupBaseURL"`
}

func NewRepresentativeLookupService(params RepresentativeLookupServiceParams) (RepresentativeLookupServiceInterface, error) {
	return &RepresentativeLookupService{
		logger:  params.Logger,
		baseURL: params.BaseURL,
	}, nil
}

// NewClient creates a new instance of ClientInterface
type ClientParams struct {
	fx.In

	Logger        loggo.LoggerInterface
	LookupService RepresentativeLookupServiceInterface
}

func NewClient(params ClientParams) (ClientInterface, error) {
	client := &DefaultClient{
		logger:        params.Logger,
		lookupService: params.LookupService,
	}
	return client, nil
}
