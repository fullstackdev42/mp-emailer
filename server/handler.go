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

// Handler struct
type Handler struct {
	Store           sessions.Store
	templateManager *shared.CustomTemplateRenderer
	campaignService campaign.ServiceInterface
	errorHandler    *shared.ErrorHandler
	EmailService    email.Service
	logger          loggo.LoggerInterface
}

// HandlerInterface defines the base logging interface for handlers
type HandlerInterface interface {
	shared.HandlerLoggable
	HandleIndex(c echo.Context) error
}

// HandlerParams defines the input parameters for Handler
type HandlerParams struct {
	fx.In
	Store           sessions.Store
	TemplateManager *shared.CustomTemplateRenderer
	CampaignService campaign.ServiceInterface
	ErrorHandler    *shared.ErrorHandler
	EmailService    email.Service
	Logger          loggo.LoggerInterface
}

// NewHandler creates a new Handler instance
func NewHandler(params HandlerParams) HandlerInterface {
	return &Handler{
		Store:           params.Store,
		templateManager: params.TemplateManager,
		campaignService: params.CampaignService,
		errorHandler:    params.ErrorHandler,
		EmailService:    params.EmailService,
		logger:          params.Logger,
	}
}

// Info implements HandlerLoggable
func (h *Handler) Info(message string, params ...interface{}) {
	h.logger.Info(message, params...)
}

// Warn implements HandlerLoggable
func (h *Handler) Warn(message string, params ...interface{}) {
	h.logger.Warn(message, params...)
}

// Error implements HandlerLoggable
func (h *Handler) Error(message string, err error, params ...interface{}) {
	h.logger.Error(message, err, params...)
}

// HandleIndex page handler
func (h *Handler) HandleIndex(c echo.Context) error {
	campaigns, err := h.campaignService.GetCampaigns()
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error fetching campaigns", 500)
	}

	data := map[string]interface{}{
		"Title":     "Home",
		"Campaigns": campaigns,
	}

	return h.templateManager.Render(c.Response(), "home", data, c)
}
