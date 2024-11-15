package user

import (
	"github.com/fullstackdev42/mp-emailer/shared"
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
	shared.BaseHandlerParams
	Service      ServiceInterface
	FlashHandler *shared.FlashHandler
	Repo         RepositoryInterface
}

// HandlerResult is the output struct for NewHandler
type HandlerResult struct {
	fx.Out
	Handler *Handler
}

// NewHandler creates a new user handler
func NewHandler(params HandlerParams) (HandlerResult, error) {
	handler := &Handler{
		BaseHandler:  shared.NewBaseHandler(params.BaseHandlerParams),
		Service:      params.Service,
		FlashHandler: params.FlashHandler,
		Repo:         params.Repo,
	}
	return HandlerResult{Handler: handler}, nil
}
