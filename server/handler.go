package server

import (
	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

// Handler struct
type Handler struct {
	Logger          loggo.LoggerInterface
	Store           sessions.Store
	templateManager *shared.CustomTemplateRenderer
	campaignService campaign.ServiceInterface
	errorHandler    *shared.ErrorHandler
	EmailService    email.Service
}

// HandleIndex page handler
func (h *Handler) HandleIndex(c echo.Context) error {
	h.Logger.Debug("Handling index request")
	isAuthenticated, _ := c.Get("IsAuthenticated").(bool)

	// Fetch campaigns using the campaign service
	campaigns, err := h.campaignService.GetCampaigns()
	if err != nil {
		h.Logger.Error("Error fetching campaigns", err)
		return h.errorHandler.HandleHTTPError(c, err, "Error fetching campaigns", 500)
	}

	data := map[string]interface{}{
		"Title":           "Home",
		"Campaigns":       campaigns,
		"IsAuthenticated": isAuthenticated,
	}

	return h.templateManager.Render(c.Response(), "home", data, c)
}
