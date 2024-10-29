package shared

import (
	"go.uber.org/fx"
)

// Module defines the shared module
//
//nolint:gochecknoglobals
var Module = fx.Options(
	fx.Provide(
		NewErrorHandler,
		NewFlashHandler,
		NewCustomTemplateRenderer,
	),
)
