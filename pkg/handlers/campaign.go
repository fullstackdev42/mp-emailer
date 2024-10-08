package handlers

import (
	"net/http"

	"github.com/fullstackdev42/mp-emailer/pkg/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleCreateCampaign(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "create_campaign.html", nil)
	}

	name := c.FormValue("name")
	template := c.FormValue("template")

	campaign := models.Campaign{
		Name:     name,
		Template: template,
	}

	if err := h.db.CreateCampaign(&campaign); err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error creating campaign")
	}

	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

func (h *Handler) HandleGetCampaigns(c echo.Context) error {
	campaigns, err := h.db.GetCampaigns()
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error fetching campaigns")
	}

	data := struct {
		Campaigns []models.Campaign
	}{
		Campaigns: campaigns,
	}

	return c.Render(http.StatusOK, "campaigns.html", data)
}

func (h *Handler) HandleUpdateCampaign(c echo.Context) error {
	id := c.Param("id")
	name := c.FormValue("name")
	template := c.FormValue("template")

	campaign := models.Campaign{
		ID:       id,
		Name:     name,
		Template: template,
	}

	if err := h.db.UpdateCampaign(&campaign); err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error updating campaign")
	}

	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

func (h *Handler) HandleDeleteCampaign(c echo.Context) error {
	id := c.Param("id")

	if err := h.db.DeleteCampaign(id); err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error deleting campaign")
	}

	return c.Redirect(http.StatusSeeOther, "/campaigns")
}
