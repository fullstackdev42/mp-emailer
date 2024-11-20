package server

import (
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/mp-emailer/shared"
	"go.uber.org/fx"
)

// Module defines the server module
var Module = fx.Module("server",
	fx.Provide(
		fx.Annotate(
			NewHandler,
			fx.As(new(HandlerInterface)),
			fx.As(new(shared.HandlerLoggable)),
		),
	),
	fx.Decorate(
		func(base HandlerInterface, logger loggo.LoggerInterface) HandlerInterface {
			return shared.NewLoggingHandlerDecorator[HandlerInterface](base, logger)
		},
	),
)
