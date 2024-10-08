package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	userID, ok := sess.Values["userID"]
	if !ok {
		// Debug print
		fmt.Println("Session values:", sess.Values)
		return c.String(http.StatusInternalServerError, "User ID not found in session")
	}

	// Try to convert to int
	var ownerID int
	switch v := userID.(type) {
	case int:
		ownerID = v
	case string:
		ownerID, err = strconv.Atoi(v)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Invalid user ID in session")
		}
	default:
		// Debug print
		fmt.Printf("Unexpected user ID type: %T\n", userID)
		return c.String(http.StatusInternalServerError, "Invalid user ID type in session")
	}

	name := c.FormValue("name")
	template := c.FormValue("template")

	campaign := models.Campaign{
		Name:     name,
		Template: template,
		OwnerID:  ownerID,
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	if err := h.db.DeleteCampaign(id); err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error deleting campaign")
	}

	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

// HandleGetCampaign handles the GET request for viewing a specific campaign
func (h *Handler) HandleGetCampaign(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	campaign, err := h.db.GetCampaignByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Campaign not found")
		}
		return h.handleError(err, http.StatusInternalServerError, "Error fetching campaign")
	}

	// You might want to exclude sensitive information here
	// For example, you might not want to show the OwnerID to the public
	campaignView := struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		Template  string    `json:"template"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID:        campaign.ID,
		Name:      campaign.Name,
		Template:  campaign.Template,
		CreatedAt: campaign.CreatedAt,
		UpdatedAt: campaign.UpdatedAt,
	}

	return c.JSON(http.StatusOK, campaignView)
}
