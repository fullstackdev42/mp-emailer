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
	fx.Provide(
		NewHandler,
	),
)

// NewHandler initializes a new Handler with the necessary dependencies
func NewHandler(
	logger loggo.LoggerInterface,
	store sessions.Store,
	tm shared.TemplateRenderer,
	cs campaign.ServiceInterface,
	eh *shared.ErrorHandler,
	es email.Service,
) *Handler {
	return &Handler{
		Logger:          logger,
		Store:           store,
		templateManager: tm,
		campaignService: cs,
		errorHandler:    eh,
		EmailService:    es,
	}
}

// RegisterRoutes registers the server routes
func RegisterRoutes(h *Handler, e *echo.Echo) {
	e.GET("/", h.HandleIndex)
}
