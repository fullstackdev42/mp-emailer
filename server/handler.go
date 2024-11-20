package server

import (
	"net/http"

	"github.com/jonesrussell/mp-emailer/campaign"
	"github.com/jonesrussell/mp-emailer/email"
	"github.com/jonesrussell/mp-emailer/shared"
	"github.com/jonesrussell/mp-emailer/version"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// Handler struct
type Handler struct {
	shared.BaseHandler
	campaignService campaign.ServiceInterface
	EmailService    email.Service
	IsShuttingDown  bool
	versionInfo     version.Info
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
	VersionInfo     version.Info
}

// NewHandler creates a new Handler instance
func NewHandler(params HandlerParams) HandlerInterface {
	return &Handler{
		BaseHandler:     shared.NewBaseHandler(params.BaseHandlerParams),
		campaignService: params.CampaignService,
		EmailService:    params.EmailService,
		versionInfo:     params.VersionInfo,
	}
}

// IndexGET page handler
func (h *Handler) IndexGET(c echo.Context) error {
	campaigns, err := h.campaignService.GetCampaigns(c.Request().Context())
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
	status := map[string]interface{}{
		"status":        "ok",
		"version":       h.versionInfo.Status(),
		"shutting_down": h.IsShuttingDown,
	}
	return c.JSON(http.StatusOK, status)
}
