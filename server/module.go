package server

import (
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/jonesrussell/loggo"
	"go.uber.org/fx"
)

// Module defines the server module
//
//nolint:gochecknoglobals
var Module = fx.Module("server",
	fx.Provide(
		fx.Annotate(
			NewHandler,
			fx.As(new(HandlerInterface)),
		),
	),
	fx.Decorate(
		func(base HandlerInterface, logger loggo.LoggerInterface) HandlerInterface {
			return shared.NewLoggingHandlerDecorator[HandlerInterface](base, logger)
		},
	),
)
