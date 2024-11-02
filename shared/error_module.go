package shared

import (
	"github.com/jonesrussell/loggo"
	"go.uber.org/fx"
)

// ErrorModule provides the error handler and its decorator
//
//nolint:gochecknoglobals
var ErrorModule = fx.Options(
	fx.Provide(
		NewErrorHandler,
		fx.Annotate(
			func(logger loggo.LoggerInterface) ErrorHandlerInterface {
				baseHandler := NewErrorHandler()
				return NewLoggingErrorHandlerDecorator(baseHandler, logger)
			},
			fx.As(new(ErrorHandlerInterface)),
		),
	),
)
