package campaign

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetCampaign(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	campaign, err := h.service.GetCampaignByID(id)
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error fetching campaign")
	}

	// Convert the campaign template to template.HTML
	campaignData := struct {
		ID        int
		Name      string
		Template  template.HTML
		CreatedAt string
		UpdatedAt string
		OwnerID   int
	}{
		ID:        campaign.ID,
		Name:      campaign.Name,
		Template:  template.HTML(campaign.Template),
		CreatedAt: campaign.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: campaign.UpdatedAt.Format("2006-01-02 15:04:05"),
		OwnerID:   campaign.OwnerID,
	}

	return c.Render(http.StatusOK, "campaign_detail.html", map[string]interface{}{
		"Campaign": campaignData,
	})
}

func (h *Handler) GetAllCampaigns(c echo.Context) error {
	campaigns, err := h.service.GetAllCampaigns()
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error fetching campaigns")
	}

	return c.Render(http.StatusOK, "campaigns.html", map[string]interface{}{
		"Campaigns": campaigns,
	})
}

func (h *Handler) CreateCampaignForm(c echo.Context) error {
	return c.Render(http.StatusOK, "campaign_create.html", nil)
}

func (h *Handler) CreateCampaign(c echo.Context) error {
	name := c.FormValue("name")
	template := c.FormValue("template")
	ownerID, err := h.getOwnerIDFromSession(c)
	if err != nil {
		return err
	}

	campaign := &Campaign{
		Name:     name,
		Template: template,
		OwnerID:  ownerID,
	}

	if err := h.service.CreateCampaign(campaign); err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error creating campaign")
	}

	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

func (h *Handler) DeleteCampaign(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	if err := h.service.DeleteCampaign(id); err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error deleting campaign")
	}

	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

func (h *Handler) EditCampaignForm(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	campaign, err := h.service.GetCampaignByID(id)
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error fetching campaign")
	}

	campaignData := struct {
		ID       int
		Name     string
		Template template.HTML
	}{
		ID:       campaign.ID,
		Name:     campaign.Name,
		Template: template.HTML(campaign.Template),
	}

	return c.Render(http.StatusOK, "campaign_edit.html", map[string]interface{}{
		"Campaign": campaignData,
	})
}

func (h *Handler) EditCampaign(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}

	name := c.FormValue("name")
	templateContent := c.FormValue("template")

	campaign := &Campaign{
		ID:       id,
		Name:     name,
		Template: templateContent,
	}

	if err := h.service.UpdateCampaign(campaign); err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error updating campaign")
	}

	return c.Redirect(http.StatusSeeOther, "/campaigns/"+strconv.Itoa(id))
}

func (h *Handler) handleError(err error, statusCode int, message string) error {
	// Log the error here (implement proper logging)
	return echo.NewHTTPError(statusCode, fmt.Sprintf("%s: %v", message, err))
}

func (h *Handler) getOwnerIDFromSession(c echo.Context) (int, error) {
	// Get the owner ID from the session
	ownerID, ok := c.Get("user_id").(int)
	if !ok {
		return 0, fmt.Errorf("user_id not found in session or not an integer")
	}
	return ownerID, nil
}
