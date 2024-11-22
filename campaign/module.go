package campaign

import (
	"github.com/jonesrussell/mp-emailer/database"
	"github.com/jonesrussell/mp-emailer/logger"
	"go.uber.org/fx"
)

// Module defines the campaign module
//
//nolint:gochecknoglobals
var Module = fx.Options(
	fx.Provide(
		// Repository
		NewRepository,

		// Base service
		fx.Annotate(
			NewService,
			fx.As(new(ServiceInterface)),
		),
		fx.Annotate(
			NewRepresentativeLookupService,
			fx.As(new(RepresentativeLookupServiceInterface)),
		),
		fx.Annotate(
			NewClient,
			fx.As(new(ClientInterface)),
		),
		NewHandler,
	),
	fx.Decorate(
		func(base ServiceInterface, logger logger.Interface) ServiceInterface {
			return NewLoggingServiceDecorator(base, logger)
		},
	),
)

// RepositoryParams for dependency injection
type RepositoryParams struct {
	fx.In
	DB     database.Database
	Logger logger.Interface
}

// ClientParams for dependency injection
type ClientParams struct {
	fx.In
	Logger        logger.Interface
	LookupService RepresentativeLookupServiceInterface
}

// NewClient creates a new instance of ClientInterface
func NewClient(params ClientParams) (ClientInterface, error) {
	client := &DefaultClient{
		Logger:        params.Logger,
		lookupService: params.LookupService,
	}
	return client, nil
}

func NewLoggingServiceDecorator(service ServiceInterface, Logger logger.Interface) ServiceInterface {
	return NewLoggingDecorator(service, Logger)
}
