package server

import (
	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// Module defines the server module
//
//nolint:gochecknoglobals
var Module = fx.Module("server",
	// Provide the base handler
	fx.Provide(
		fx.Annotate(
			NewHandler,
			fx.As(new(HandlerInterface)),
		),
	),
	// Decorate at module level
	fx.Decorate(
		func(base HandlerInterface, logger loggo.LoggerInterface) HandlerInterface {
			return NewLoggingHandlerDecorator(base, logger)
		},
	),
)

// NewHandler initializes a new Handler with the necessary dependencies
func NewHandler(
	store sessions.Store,
	tm *shared.CustomTemplateRenderer,
	cs campaign.ServiceInterface,
	eh *shared.ErrorHandler,
	es email.Service,
) HandlerInterface {
	return &Handler{
		Store:           store,
		templateManager: tm,
		campaignService: cs,
		errorHandler:    eh,
		EmailService:    es,
	}
}

// RegisterRoutes registers the server routes
func RegisterRoutes(h HandlerInterface, e *echo.Echo) {
	e.GET("/", h.HandleIndex)
}
