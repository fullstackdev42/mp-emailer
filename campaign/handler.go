package campaign

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fullstackdev42/mp-emailer/pkg/services"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service                     *Service
	logger                      loggo.LoggerInterface
	representativeLookupService *RepresentativeLookupService
	emailService                services.EmailService
	client                      ClientInterface
}

func NewHandler(service *Service, logger loggo.LoggerInterface, representativeLookupService *RepresentativeLookupService, emailService services.EmailService, client ClientInterface) *Handler {
	return &Handler{
		service:                     service,
		logger:                      logger,
		representativeLookupService: representativeLookupService,
		emailService:                emailService,
		client:                      client,
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

	// Log all form values
	for key, values := range c.Request().Form {
		h.logger.Debug("Form value", "key", key, "values", values)
	}

	postalCode := c.FormValue("postal_code")
	h.logger.Debug("Raw postal code received", "postalCode", postalCode)

	if postalCode == "" {
		h.logger.Warn("Empty postal code submitted")
		return echo.NewHTTPError(http.StatusBadRequest, "Postal code is required")
	}

	postalCode = strings.ToUpper(strings.ReplaceAll(postalCode, " ", ""))
	h.logger.Debug("Processed postal code", "postalCode", postalCode)

	postalCodeRegex := regexp.MustCompile(`^[ABCEGHJ-NPRSTVXY]\d[ABCEGHJ-NPRSTV-Z]\d[ABCEGHJ-NPRSTV-Z]\d$`)
	if !postalCodeRegex.MatchString(postalCode) {
		h.logger.Warn("Invalid postal code submitted", "postalCode", postalCode, "regexPattern", postalCodeRegex.String())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid postal code format")
	}

	h.logger.Info("Valid postal code received", "postalCode", postalCode)

	mpFinder := NewMPFinder(h.client, h.logger)

	mp, err := mpFinder.FindMP(postalCode)
	if err != nil {
		h.logger.Error("Error finding MP", err)
		return h.handleError(err, http.StatusInternalServerError, "Error finding MP")
	}
	h.logger.Debug("MP found", "name", mp.Name, "email", mp.Email)

	campaignID := c.Param("id")
	h.logger.Debug("Campaign ID", "id", campaignID)
	campaign, err := h.service.GetCampaignByID(campaignID)
	if err != nil {
		h.logger.Error("Error fetching campaign", err)
		return h.handleError(err, http.StatusInternalServerError, "Error fetching campaign")
	}
	h.logger.Debug("Campaign fetched", "name", campaign.Name)

	userData := map[string]string{
		"First Name":    c.FormValue("first_name"),
		"Last Name":     c.FormValue("last_name"),
		"Address 1":     c.FormValue("address_1"),
		"City":          c.FormValue("city"),
		"Province":      c.FormValue("province"),
		"Postal Code":   c.FormValue("postal_code"),
		"Email Address": c.FormValue("email"),
	}
	h.logger.Debug("User data", "userData", userData)

	emailContent := h.composeEmail(mp, campaign, userData)
	h.logger.Debug("Composed email content", "content", emailContent)

	data := struct {
		Email   string
		Content string
	}{
		Email:   mp.Email,
		Content: emailContent,
	}
	h.logger.Debug("Data for email template", "data", data)

	err = c.Render(http.StatusOK, "email.html", data)
	if err != nil {
		h.logger.Error("Error rendering email template", err)
		return h.handleError(err, http.StatusInternalServerError, "Error rendering email template")
	}

	h.logger.Info("Email template rendered successfully")
	return nil
}

func (h *Handler) composeEmail(mp Representative, campaign *Campaign, userData map[string]string) string {
	content := campaign.Template

	// Replace MP information
	content = strings.ReplaceAll(content, "{{MP's Name}}", mp.Name)

	// Replace user information
	for key, value := range userData {
		content = strings.ReplaceAll(content, "{{"+key+"}}", value)
	}

	// Replace date
	currentDate := time.Now().Format("January 2, 2006")
	content = strings.ReplaceAll(content, "{{Date}}", currentDate)

	return content
}

func (h *Handler) HandleMPLookup(c echo.Context) error {
	postalCode := c.FormValue("postal_code")
	representatives, err := h.representativeLookupService.FetchRepresentatives(postalCode)
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error fetching representatives")
	}

	return c.JSON(http.StatusOK, representatives)
}

func (h *Handler) HandleRepresentativeLookup(c echo.Context) error {
	postalCode := c.FormValue("postal_code")
	representativeType := c.FormValue("type") // e.g., "MP", "Premier", "Prime Minister"

	representatives, err := h.representativeLookupService.FetchRepresentatives(postalCode)
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error fetching representatives")
	}

	filters := map[string]string{
		"type": representativeType,
	}

	filteredRepresentatives := h.representativeLookupService.FilterRepresentatives(representatives, filters)

	return c.JSON(http.StatusOK, filteredRepresentatives)
}
