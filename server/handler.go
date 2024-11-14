package server

import (
	"net/http"
	"time"

	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// Handler struct
type Handler struct {
	shared.BaseHandler
	campaignService campaign.ServiceInterface
	EmailService    email.Service
	IsShuttingDown  bool
}

// HandlerInterface defines the interface for handlers
type HandlerInterface interface {
	shared.HandlerLoggable
	IndexGET(c echo.Context) error
	HealthCheck(c echo.Context) error
}

// HandlerParams defines the input parameters for Handler
type HandlerParams struct {
	fx.In
	shared.BaseHandlerParams
	CampaignService campaign.ServiceInterface
	EmailService    email.Service
}

// NewHandler creates a new Handler instance
func NewHandler(params HandlerParams) HandlerInterface {
	return &Handler{
		BaseHandler:     shared.NewBaseHandler(params.BaseHandlerParams),
		campaignService: params.CampaignService,
		EmailService:    params.EmailService,
	}
}

// IndexGET page handler
func (h *Handler) IndexGET(c echo.Context) error {
	campaigns, err := h.campaignService.GetCampaigns()
	if err != nil {
		h.Logger.Error("Error fetching campaigns", err)
		status, msg := h.MapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	return c.Render(http.StatusOK, "home", &shared.Data{
		Title:    "Home",
		PageName: "home",
		Content: map[string]interface{}{
			"Campaigns": campaigns,
		},
	})
}

// HealthCheck health check endpoint
func (h *Handler) HealthCheck(c echo.Context) error {
	status := http.StatusOK
	response := map[string]interface{}{
		"status": "healthy",
		"time":   time.Now().UTC(),
	}

	if h.IsShuttingDown {
		status = http.StatusServiceUnavailable
		response["status"] = "shutting_down"
	}

	return c.JSON(status, response)
}
