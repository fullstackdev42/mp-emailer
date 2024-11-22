package user

import (
	"github.com/jonesrussell/mp-emailer/database"
	"github.com/jonesrussell/mp-emailer/logger"
	"go.uber.org/fx"
)

// Module defines the user module
// nolint:gochecknoglobals
var Module = fx.Options(
	fx.Provide(
		func(db database.Database) RepositoryParams {
			return RepositoryParams{
				DB: db,
			}
		},
		NewRepository,
		fx.Annotate(NewService,
			fx.As(new(ServiceInterface)),
		),
		NewHandler,
	),
	// Add module-level decoration
	fx.Decorate(
		func(base ServiceInterface, log logger.Interface) ServiceInterface {
			return NewLoggingDecorator(base, log)
		},
	),
)

// HandlerResult is the output struct for NewHandler
type HandlerResult struct {
	fx.Out
	Handler *Handler
}
