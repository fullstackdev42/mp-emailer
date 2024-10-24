package shared

import (
	"go.uber.org/fx"
)

// Module defines the shared module
//
//nolint:gochecknoglobals
var Module = fx.Module(
	"shared",
	fx.Provide(
		NewTemplateRenderer,
		NewErrorHandler,
	),
)
