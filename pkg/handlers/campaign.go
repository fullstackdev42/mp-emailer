package handlers

import (
	"net/http"

	"github.com/fullstackdev42/mp-emailer/pkg/models"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleCreateCampaign(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "create_campaign.html", nil)
	}

	// Extract user ID from session
	sess, err := session.Get("mpe", c)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get session")
	}
	userID, ok := sess.Values["userID"].(string)
	if !ok {
		return c.String(http.StatusInternalServerError, "Invalid user ID in session")
	}

	name := c.FormValue("name")
	template := c.FormValue("template")

	campaign := models.Campaign{
		Name:     name,
		Template: template,
		OwnerID:  userID, // Set the owner ID here
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

	data := map[string]interface{}{
		"Campaigns": campaigns,
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
