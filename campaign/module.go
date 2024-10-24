package campaign

import (
	"go.uber.org/fx"
)

// ProvideModule bundles and provides all campaign-related dependencies
func ProvideModule() fx.Option {
	return fx.Options(
		fx.Provide(NewRepository),
		fx.Provide(NewService),
	)
}

// ServiceResult is the output struct for NewService
type ServiceResult struct {
	fx.Out
	Service ServiceInterface
}
