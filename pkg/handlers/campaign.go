package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"html/template"

	"github.com/fullstackdev42/mp-emailer/pkg/database"
	"github.com/fullstackdev42/mp-emailer/pkg/models"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleCreateCampaign(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "campaign_create.html", nil)
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

	// Retrieve the database connection from the context
	db := c.Get("db").(*database.DB)

	if err := db.CreateCampaign(&campaign); err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error creating campaign")
	}

	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

func (h *Handler) HandleGetCampaigns(c echo.Context) error {
	// Retrieve the database connection from the context
	db := c.Get("db").(*database.DB)

	campaigns, err := db.GetCampaigns()
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error fetching campaigns")
	}

	data := map[string]interface{}{
		"Campaigns": campaigns,
	}

	return c.Render(http.StatusOK, "campaigns.html", data)
}

func (h *Handler) HandleDeleteCampaign(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	// Retrieve the database connection from the context
	db := c.Get("db").(*database.DB)

	if err := db.DeleteCampaign(id); err != nil {
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

	// Retrieve the database connection from the context
	db := c.Get("db").(*database.DB)

	campaign, err := db.GetCampaignByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Campaign not found")
		}
		return h.handleError(err, http.StatusInternalServerError, "Error fetching campaign")
	}

	// Convert the campaign template to template.HTML
	campaignData := struct {
		ID        int
		Name      string
		Template  template.HTML
		CreatedAt time.Time
		UpdatedAt time.Time
		OwnerID   int
	}{
		ID:        campaign.ID,
		Name:      campaign.Name,
		Template:  template.HTML(campaign.Template),
		CreatedAt: campaign.CreatedAt,
		UpdatedAt: campaign.UpdatedAt,
		OwnerID:   campaign.OwnerID,
	}

	return c.Render(http.StatusOK, "campaign_detail.html", map[string]interface{}{
		"Campaign": campaignData,
	})
}

// Add this new function
func (h *Handler) HandleEditCampaign(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	// Retrieve the database connection from the context
	db := c.Get("db").(*database.DB)

	if c.Request().Method == http.MethodGet {
		campaign, err := db.GetCampaignByID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				return echo.NewHTTPError(http.StatusNotFound, "Campaign not found")
			}
			return h.handleError(err, http.StatusInternalServerError, "Error fetching campaign")
		}

		// Convert the campaign template to template.HTML
		campaignData := struct {
			ID        int
			Name      string
			Template  template.HTML
			CreatedAt time.Time
			UpdatedAt time.Time
			OwnerID   int
		}{
			ID:        campaign.ID,
			Name:      campaign.Name,
			Template:  template.HTML(campaign.Template),
			CreatedAt: campaign.CreatedAt,
			UpdatedAt: campaign.UpdatedAt,
			OwnerID:   campaign.OwnerID,
		}

		return c.Render(http.StatusOK, "campaign_edit.html", map[string]interface{}{
			"Campaign": campaignData,
		})
	} else if c.Request().Method == http.MethodPost {
		name := c.FormValue("name")
		templateContent := c.FormValue("template")

		campaign := &models.Campaign{
			ID:       id,
			Name:     name,
			Template: templateContent,
		}

		if err := db.UpdateCampaign(campaign); err != nil {
			return h.handleError(err, http.StatusInternalServerError, "Error updating campaign")
		}

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/campaigns/%d", id))
	}

	return echo.NewHTTPError(http.StatusMethodNotAllowed, "Method not allowed")
}
