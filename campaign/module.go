package campaign

import (
	"github.com/fullstackdev42/mp-emailer/database"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// Module defines the campaign module
//
//nolint:gochecknoglobals
var Module = fx.Options(
	fx.Provide(
		// Repository
		fx.Annotate(
			NewRepository,
			fx.As(new(RepositoryInterface)),
		),
		// Base service
		fx.Annotate(
			NewService,
			fx.As(new(ServiceInterface)),
		),
		fx.Annotate(
			NewRepresentativeLookupService,
			fx.As(new(RepresentativeLookupServiceInterface)),
		),
		fx.Annotate(
			NewClient,
			fx.As(new(ClientInterface)),
		),
		NewHandler,
	),
	// Add module-level decoration
	fx.Decorate(
		func(base ServiceInterface, logger loggo.LoggerInterface) ServiceInterface {
			return NewLoggingServiceDecorator(base, logger)
		},
	),
)

// RepositoryParams for dependency injection
type RepositoryParams struct {
	fx.In
	DB     database.Interface
	Logger loggo.LoggerInterface
}

// NewRepository creates a new campaign repository
func NewRepository(params RepositoryParams) (RepositoryInterface, error) {
	repo := &Repository{
		db: params.DB,
	}
	return repo, nil
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

func NewLoggingServiceDecorator(service ServiceInterface, logger loggo.LoggerInterface) ServiceInterface {
	return NewLoggingDecorator(service, logger)
}
