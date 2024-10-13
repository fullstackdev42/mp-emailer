package campaign

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fullstackdev42/mp-emailer/pkg/api"
	"github.com/fullstackdev42/mp-emailer/pkg/models"
	"github.com/fullstackdev42/mp-emailer/pkg/services"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service      *Service
	logger       loggo.LoggerInterface
	client       api.ClientInterface
	emailService services.EmailService
}

func NewHandler(service *Service, logger loggo.LoggerInterface, client api.ClientInterface, emailService services.EmailService) *Handler {
	return &Handler{
		service:      service,
		logger:       logger,
		client:       client,
		emailService: emailService,
	}
}

func (h *Handler) GetCampaign(c echo.Context) error {
	id := c.Param("id")
	campaign, err := h.service.GetCampaignByID(id)
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error fetching campaign")
	}

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
	id := c.Param("id")
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
	h.logger.Error(message, err)
	return echo.NewHTTPError(statusCode, message)
}

func (h *Handler) getOwnerIDFromSession(c echo.Context) (int, error) {
	// Get the owner ID from the session
	ownerID, ok := c.Get("user_id").(int)
	if !ok {
		return 0, fmt.Errorf("user_id not found in session or not an integer")
	}
	return ownerID, nil
}

func (h *Handler) SendCampaign(c echo.Context) error {
	h.logger.Info("Handling campaign submit request")

	postalCode := c.FormValue("postalCode")
	postalCode = strings.ToUpper(strings.ReplaceAll(postalCode, " ", ""))

	postalCodeRegex := regexp.MustCompile(`^[ABCEGHJ-NPRSTVXY]\d[ABCEGHJ-NPRSTV-Z]\d[ABCEGHJ-NPRSTV-Z]\d$`)
	if !postalCodeRegex.MatchString(postalCode) {
		h.logger.Warn("Invalid postal code submitted", "postalCode", postalCode)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid postal code format")
	}

	mpFinder := services.NewMPFinder(h.client, h.logger)

	mp, err := mpFinder.FindMP(postalCode)
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error finding MP")
	}

	campaignID := c.Param("id")
	campaign, err := h.service.GetCampaignByID(campaignID)
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error fetching campaign")
	}

	emailContent := h.composeEmail(mp, campaign)

	data := struct {
		Email   string
		Content string
	}{
		Email:   mp.Email,
		Content: emailContent,
	}

	return c.Render(http.StatusOK, "email.html", data)
}

func (h *Handler) composeEmail(mp models.Representative, campaign *Campaign) string {
	// Here you would use the campaign template and replace any placeholders
	// with the MP's information. For now, we'll use a simple format:
	return fmt.Sprintf("Dear %s,\n\n%s\n\nBest regards,\nYour constituent", mp.Name, campaign.Template)
}
