package user

import (
	"github.com/fullstackdev42/mp-emailer/config"
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
		ProvideSessionManager,
	),
	// Add module-level decoration
	fx.Decorate(
		func(base ServiceInterface, logger loggo.LoggerInterface) ServiceInterface {
			return NewLoggingDecorator(base, logger)
		},
	),
)

// HandlerResult is the output struct for NewHandler
type HandlerResult struct {
	fx.Out
	Handler *Handler
}

// Provide session manager
func ProvideSessionManager(store sessions.Store, cfg *config.Config) SessionManager {
	return NewSessionManager(store, cfg)
}
