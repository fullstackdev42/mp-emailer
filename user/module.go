package user

import (
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"go.uber.org/fx"
)

// Module defines the user module
// nolint:gochecknoglobals
var Module = fx.Options(
	fx.Provide(
		NewRepository,
		fx.Annotate(NewService,
			fx.As(new(ServiceInterface)),
		),
		NewHandler,
	),
	// Add module-level decoration
	fx.Decorate(
		func(base ServiceInterface, logger loggo.LoggerInterface) ServiceInterface {
			return NewLoggingDecorator(base, logger)
		},
	),
)

// HandlerParams for dependency injection
type HandlerParams struct {
	fx.In
	Config          *config.Config
	Service         ServiceInterface
	Store           sessions.Store
	TemplateManager shared.TemplateRendererInterface
	Repo            RepositoryInterface
	ErrorHandler    shared.ErrorHandlerInterface
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
		Service:         params.Service,
		Store:           params.Store,
		SessionName:     params.Config.SessionName,
		Config:          params.Config,
		TemplateManager: params.TemplateManager,
		Repo:            params.Repo,
		ErrorHandler:    params.ErrorHandler,
		FlashHandler:    params.FlashHandler,
	}
	return HandlerResult{Handler: handler}, nil
}
