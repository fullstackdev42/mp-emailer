package user

import (
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/internal/database"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"go.uber.org/fx"
)

// ProvideModule provides the user module dependencies
func ProvideModule() fx.Option {
	return fx.Options(
		fx.Provide(
			NewRepository,
			NewService,
			// NewHandler,
		),
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
func NewService(repo RepositoryInterface, logger loggo.LoggerInterface) (ServiceInterface, error) {
	service := &Service{
		repo:   repo,
		logger: logger,
	}
	return service, nil
}

// HandlerResult is the output struct for NewHandler
type HandlerResult struct {
	fx.Out
	Handler *Handler
}

// NewHandler creates a new user handler
func NewHandler(
	repo RepositoryInterface,
	service ServiceInterface,
	logger loggo.LoggerInterface,
	store sessions.Store,
	config *config.Config,
) (HandlerResult, error) {
	handler := &Handler{
		repo:        repo,
		service:     service,
		Logger:      logger,
		Store:       store,
		SessionName: config.SessionName,
		Config:      config,
	}
	return HandlerResult{Handler: handler}, nil
}
