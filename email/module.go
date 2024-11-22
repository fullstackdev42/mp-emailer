package email

import (
	"github.com/jonesrussell/mp-emailer/logger"
	"go.uber.org/fx"
)

// Module provides email services to the application
//
//nolint:gochecknoglobals
var Module = fx.Options(
	fx.Provide(
		// Provide the email service factory
		fx.Annotate(
			NewEmailService,
			fx.As(new(Service)),
		),
		// Provide the SMTP client implementation
		fx.Annotate(
			func() SMTPClient {
				return &SMTPClientImpl{}
			},
			fx.As(new(SMTPClient)),
		),
	),
)

// Params holds the dependencies needed to create an email service
type Params struct {
	fx.In

	Config Config
	Logger logger.Interface
}

// Result holds the email service instance
type Result struct {
	fx.Out

	EmailService Service
}
