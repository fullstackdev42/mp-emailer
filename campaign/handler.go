package campaign

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service                     *Service
	logger                      loggo.LoggerInterface
	representativeLookupService *RepresentativeLookupService
	emailService                email.Service
	client                      ClientInterface
}

func NewHandler(service *Service, logger loggo.LoggerInterface, representativeLookupService *RepresentativeLookupService, emailService email.Service, client ClientInterface) *Handler {
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

	postalCode, err := h.extractAndValidatePostalCode(c)
	if err != nil {
		return err
	}

	mp, err := h.findMP(postalCode)
	if err != nil {
		return err
	}

	campaign, err := h.fetchCampaign(c.Param("id"))
	if err != nil {
		return err
	}

	userData := h.extractUserData(c)
	emailContent := h.composeEmail(mp, campaign, userData)

	return h.renderEmailTemplate(c, mp.Email, emailContent)
}

func (h *Handler) composeEmail(mp Representative, campaign *Campaign, userData map[string]string) string {
	emailTemplate := campaign.Template
	for key, value := range userData {
		placeholder := fmt.Sprintf("{{%s}}", key)
		emailTemplate = strings.ReplaceAll(emailTemplate, placeholder, value)
	}
	// Handle the token with the apostrophe
	emailTemplate = strings.ReplaceAll(emailTemplate, "{{MP's Name}}", mp.Name)
	emailTemplate = strings.ReplaceAll(emailTemplate, "{{MPEmail}}", mp.Email)
	emailTemplate = strings.ReplaceAll(emailTemplate, "{{Date}}", time.Now().Format("2006-01-02"))
	return emailTemplate
}

func (h *Handler) extractAndValidatePostalCode(c echo.Context) (string, error) {
	postalCode := c.FormValue("postal_code")
	if postalCode == "" {
		h.logger.Warn("Empty postal code submitted")
		return "", echo.NewHTTPError(http.StatusBadRequest, "Postal code is required")
	}

	postalCode = strings.ToUpper(strings.ReplaceAll(postalCode, " ", ""))
	postalCodeRegex := regexp.MustCompile(`^[ABCEGHJ-NPRSTVXY]\d[ABCEGHJ-NPRSTV-Z]\d[ABCEGHJ-NPRSTV-Z]\d$`)
	if !postalCodeRegex.MatchString(postalCode) {
		h.logger.Warn("Invalid postal code submitted", "postalCode", postalCode)
		return "", echo.NewHTTPError(http.StatusBadRequest, "Invalid postal code format")
	}

	return postalCode, nil
}

func (h *Handler) findMP(postalCode string) (Representative, error) {
	mpFinder := NewMPFinder(h.client, h.logger)
	mp, err := mpFinder.FindMP(postalCode)
	if err != nil {
		h.logger.Error("Error finding MP", err)
		return Representative{}, h.handleError(err, http.StatusInternalServerError, "Error finding MP")
	}
	return mp, nil
}

func (h *Handler) fetchCampaign(id string) (*Campaign, error) {
	campaign, err := h.service.GetCampaignByID(id)
	if err != nil {
		h.logger.Error("Error fetching campaign", err)
		return nil, h.handleError(err, http.StatusInternalServerError, "Error fetching campaign")
	}
	return campaign, nil
}

func (h *Handler) extractUserData(c echo.Context) map[string]string {
	return map[string]string{
		"First Name":    c.FormValue("first_name"),
		"Last Name":     c.FormValue("last_name"),
		"Address 1":     c.FormValue("address_1"),
		"City":          c.FormValue("city"),
		"Province":      c.FormValue("province"),
		"Postal Code":   c.FormValue("postal_code"),
		"Email Address": c.FormValue("email"),
	}
}

func (h *Handler) renderEmailTemplate(c echo.Context, email, content string) error {
	data := struct {
		Email   string
		Content template.HTML // Use template.HTML to ensure HTML content is rendered correctly
	}{
		Email:   email,
		Content: template.HTML(content), // Convert content to template.HTML
	}

	h.logger.Debug("Data for email template", "data", data)

	// Attempt to render the email template
	err := c.Render(http.StatusOK, "email.html", map[string]interface{}{
		"Data": data,
	})
	if err != nil {
		h.logger.Error("Error rendering email template", err)
		return h.handleError(err, http.StatusInternalServerError, "Error rendering email template")
	}

	h.logger.Info("Email template rendered successfully")
	return nil
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
