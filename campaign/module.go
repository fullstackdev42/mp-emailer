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

// ProvideModule bundles and provides all campaign-related dependencies
func ProvideModule() fx.Option {
	return fx.Options(
		fx.Provide(
			NewRepository,
			NewService,
			NewHandler,
		),
		fx.Invoke(InvokeModule),
	)
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

// NewRepository creates a new campaign repository
func NewRepository(db *database.DB) (RepositoryInterface, error) {
	return &Repository{db: db}, nil
}

// NewService creates a new campaign service
func NewService(repo RepositoryInterface) (*Service, error) {
	validate := validator.New()
	service := &Service{
		repo:     repo,
		validate: validate,
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

func InvokeModule(e *echo.Echo, handler *Handler) {
	// Register campaign-related routes
	e.GET("/campaign", handler.GetAllCampaigns)
	e.POST("/campaign", handler.CreateCampaign)
	e.GET("/campaign/:id", handler.CampaignGET)
	e.PUT("/campaign/:id", handler.EditCampaign)
	e.DELETE("/campaign/:id", handler.DeleteCampaign)
	e.POST("/campaign/:id/send", handler.SendCampaign)
}
