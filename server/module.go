package server

import (
	"github.com/jonesrussell/mp-emailer/logger"
	"github.com/jonesrussell/mp-emailer/shared"
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
			fx.As(new(shared.HandlerLoggable)),
		),
	),
	fx.Decorate(
		func(base HandlerInterface, log logger.Interface) HandlerInterface {
			return shared.NewLoggingHandlerDecorator[HandlerInterface](base, log)
		},
	),
)
