package shared

import (
	"github.com/jonesrussell/mp-emailer/logger"
	"go.uber.org/fx"
)

// ErrorModule provides the error handler and its decorator
//
//nolint:gochecknoglobals
var ErrorModule = fx.Options(
	fx.Provide(
		NewErrorHandler,
		fx.Annotate(
			func(log logger.Interface) ErrorHandlerInterface {
				baseHandler := NewErrorHandler()
				return NewLoggingErrorHandlerDecorator(baseHandler, log)
			},
			fx.As(new(ErrorHandlerInterface)),
		),
	),
)
