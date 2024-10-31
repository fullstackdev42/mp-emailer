package server

import (
	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

// Handler struct
type Handler struct {
	Store           sessions.Store
	templateManager *shared.CustomTemplateRenderer
	campaignService campaign.ServiceInterface
	errorHandler    *shared.ErrorHandler
	EmailService    email.Service
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
