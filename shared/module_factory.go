package shared

import (
	"go.uber.org/fx"
)

func NewModule(name string, providers []interface{}, decorators []interface{}) fx.Option {
	return fx.Module(name,
		fx.Provide(providers...),
		fx.Decorate(decorators...),
	)
}
