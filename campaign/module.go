package campaign

import (
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/mp-emailer/database/core"
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
	// Add module-level decoration
	fx.Decorate(
		func(base ServiceInterface, logger loggo.LoggerInterface) ServiceInterface {
			return NewLoggingServiceDecorator(base, logger)
		},
	),
)

// RepositoryParams for dependency injection
type RepositoryParams struct {
	fx.In
	DB     core.Interface
	Logger loggo.LoggerInterface
}

// ClientParams for dependency injection
type ClientParams struct {
	fx.In
	Logger        loggo.LoggerInterface
	LookupService RepresentativeLookupServiceInterface
}

// NewClient creates a new instance of ClientInterface
func NewClient(params ClientParams) (ClientInterface, error) {
	client := &DefaultClient{
		logger:        params.Logger,
		lookupService: params.LookupService,
	}
	return client, nil
}

func NewLoggingServiceDecorator(service ServiceInterface, logger loggo.LoggerInterface) ServiceInterface {
	return NewLoggingDecorator(service, logger)
}
