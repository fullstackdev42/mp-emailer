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
var Module = fx.Module("campaign",
	fx.Provide(
		NewRepository,
		NewService,
		NewHandler,
		NewRepresentativeLookupService,
		NewClient,
		func() string {
			return "https://represent.opennorth.ca"
		},
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
func NewRepository(db *database.DB) (RepositoryInterface, error) {
	return &Repository{db: db}, nil
}

// NewService creates a new campaign service
func NewService(repo RepositoryInterface) (ServiceResult, error) {
	validate := validator.New()
	service := ServiceResult{
		Service: &Service{
			repo:     repo,
			validate: validate,
		},
	}
	return service, nil
}

// NewHandler initializes a new Handler
func NewHandler(
	service ServiceInterface,
	logger loggo.LoggerInterface,
	representativeLookupService RepresentativeLookupServiceInterface,
	emailService email.Service,
	client ClientInterface,
) *Handler {
	return &Handler{
		service:                     service,
		logger:                      logger,
		representativeLookupService: representativeLookupService,
		emailService:                emailService,
		client:                      client,
		errorHandler:                shared.NewErrorHandler(logger),
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
func NewRepresentativeLookupService(logger loggo.LoggerInterface, baseURL string) (RepresentativeLookupServiceInterface, error) {
	return &RepresentativeLookupService{
		logger:  logger,
		baseURL: baseURL,
	}, nil
}

// NewClient creates a new instance of ClientInterface
func NewClient(logger loggo.LoggerInterface, lookupService RepresentativeLookupServiceInterface) (ClientInterface, error) {
	client := &DefaultClient{
		logger:        logger,
		lookupService: lookupService,
	}
	return client, nil
}
