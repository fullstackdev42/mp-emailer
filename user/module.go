package user

import (
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/internal/database"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// ProvideModule provides the user module dependencies
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

// NewRepository creates a new user repository
func NewRepository(db *database.DB, logger loggo.LoggerInterface) RepositoryInterface {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// ServiceResult is the output struct for NewService
type ServiceResult struct {
	fx.Out
	Service ServiceInterface
}

// NewService creates a new user service
func NewService(repo RepositoryInterface, logger loggo.LoggerInterface) (ServiceResult, error) {
	service := &Service{
		repo:   repo,
		logger: logger,
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
) (HandlerResult, error) {
	handler := &Handler{
		service:     service,
		Logger:      logger,
		Store:       sessions,
		SessionName: cfg.SessionName,
		Config:      cfg,
	}
	return HandlerResult{Handler: handler}, nil
}

func InvokeModule(e *echo.Echo, handler *Handler) {
	// Register routes
	e.GET("/user/register", handler.RegisterGET)
	e.POST("/user/register", handler.RegisterPOST)
	e.GET("/user/login", handler.LoginGET)
	e.POST("/user/login", handler.LoginPOST)
	e.GET("/user/logout", handler.LogoutGET)
}
