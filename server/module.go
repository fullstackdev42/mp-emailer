package server

import (
	"embed"
	"fmt"
	"html/template"

	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// ProvideModule provides the server module dependencies
func ProvideModule() fx.Option {
	return fx.Options(
		fx.Provide(
			NewHandler,
			NewTemplateManager,
			shared.NewErrorHandler,
		),
	)
}

// NewTemplateManager initializes and parses templates
func NewTemplateManager(templateFiles embed.FS) (*TemplateManager, error) {
	tm := &TemplateManager{}

	// Parse all templates
	tmpl, err := template.New("").ParseFS(templateFiles, "web/templates/**/*.gohtml")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	// Ensure the "app" template exists
	if tmpl.Lookup("app") == nil {
		return nil, fmt.Errorf("layout template 'app' not found")
	}

	tm.templates = tmpl
	return tm, nil
}

// NewHandler initializes a new Handler with the necessary dependencies
func NewHandler(
	logger loggo.LoggerInterface,
	store sessions.Store,
	tm *TemplateManager,
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
