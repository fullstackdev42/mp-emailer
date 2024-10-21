package campaign

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

// Handler handles the HTTP requests for the campaign service
type Handler struct {
	service                     ServiceInterface
	logger                      loggo.LoggerInterface
	representativeLookupService RepresentativeLookupServiceInterface
	emailService                email.Service
	client                      ClientInterface
	errorHandler                *shared.ErrorHandler
}

// NewHandler initializes a new Handler
func NewHandler(
	service ServiceInterface,
	logger loggo.LoggerInterface,
	representativeLookupService RepresentativeLookupServiceInterface,
	emailService email.Service,
	client ClientInterface,
) *Handler {
	return &Handler{
		service:                     service,
		logger:                      logger,
		representativeLookupService: representativeLookupService,
		emailService:                emailService,
		client:                      client,
		errorHandler:                shared.NewErrorHandler(logger),
	}
}

// CampaignGET handles GET requests for campaign details
func (h *Handler) CampaignGET(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.errorHandler.HandleError(c, err, http.StatusBadRequest, "Invalid campaign ID")
	}

	campaign, err := h.service.FetchCampaign(id)
	if err != nil {
		if errors.Is(err, ErrCampaignNotFound) {
			return h.errorHandler.HandleError(c, err, http.StatusNotFound, "Campaign not found")
		}
		return h.errorHandler.HandleError(c, err, http.StatusInternalServerError, "Error fetching campaign")
	}

	pageData := shared.PageData{
		Title:   "Campaign Details",
		Content: campaign,
	}

	return c.Render(http.StatusOK, "campaign_details.html", pageData)
}

// GetAllCampaigns handles GET requests for all campaigns
func (h *Handler) GetAllCampaigns(c echo.Context) error {
	campaigns, err := h.service.GetAllCampaigns()
	if err != nil {
		return h.errorHandler.HandleError(c, err, http.StatusInternalServerError, "Error fetching campaigns")
	}
	return c.Render(http.StatusOK, "campaigns.html", map[string]interface{}{"Campaigns": campaigns})
}

// CreateCampaignForm handles GET requests for the campaign creation form
func (h *Handler) CreateCampaignForm(c echo.Context) error {
	return c.Render(http.StatusOK, "campaign_create.html", nil)
}

// CreateCampaign handles POST requests for creating a new campaign
func (h *Handler) CreateCampaign(c echo.Context) error {
	name := c.FormValue("name")
	template := c.FormValue("template")
	ownerID, err := user.GetOwnerIDFromSession(c)
	if err != nil {
		return h.errorHandler.HandleError(c, err, http.StatusUnauthorized, "Unauthorized")
	}
	campaign := &Campaign{
		Name:     name,
		Template: template,
		OwnerID:  ownerID,
	}
	if err := h.service.CreateCampaign(campaign); err != nil {
		return h.errorHandler.HandleError(c, err, http.StatusInternalServerError, "Error creating campaign")
	}
	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

// DeleteCampaign handles DELETE requests for deleting a campaign
func (h *Handler) DeleteCampaign(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}
	if err := h.service.DeleteCampaign(id); err != nil {
		return h.errorHandler.HandleError(c, err, http.StatusInternalServerError, "Error deleting campaign")
	}
	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

// EditCampaignForm handles GET requests for the campaign edit form
func (h *Handler) EditCampaignForm(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.errorHandler.HandleError(c, err, http.StatusBadRequest, "Invalid campaign ID")
	}
	campaign, err := h.service.FetchCampaign(id)
	if err != nil {
		return h.errorHandler.HandleError(c, err, http.StatusInternalServerError, "Error fetching campaign")
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
	return c.Render(http.StatusOK, "campaign_edit.html", map[string]interface{}{"Campaign": campaignData})
}

// EditCampaign handles POST requests for updating a campaign
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
		return h.errorHandler.HandleError(c, err, http.StatusInternalServerError, "Error updating campaign")
	}
	return c.Redirect(http.StatusSeeOther, "/campaigns/"+strconv.Itoa(id))
}

// SendCampaign handles POST requests for sending a campaign
func (h *Handler) SendCampaign(c echo.Context) error {
	h.logger.Info("Handling campaign submit request")
	postalCode, err := extractAndValidatePostalCode(c)
	if err != nil {
		h.logger.Warn("Invalid postal code submitted", "error", err)
		return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
			"Error": "Invalid postal code",
		})
	}
	mp, err := h.representativeLookupService.FetchRepresentatives(postalCode)
	if err != nil {
		return h.errorHandler.HandleError(c, err, http.StatusInternalServerError, "Error finding MP")
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.errorHandler.HandleError(c, err, http.StatusBadRequest, "Invalid campaign ID")
	}
	campaign, err := h.service.FetchCampaign(id)
	if err != nil {
		return h.errorHandler.HandleError(c, err, http.StatusInternalServerError, "Error fetching campaign")
	}
	userData := extractUserData(c)

	if len(mp) == 0 {
		return h.errorHandler.HandleError(c, errors.New("no representatives found"), http.StatusNotFound, "No representatives found")
	}
	representative := mp[0]

	emailContent := h.service.ComposeEmail(representative, campaign, userData)
	return h.RenderEmailTemplate(c, representative.Email, emailContent)
}

// HandleError renders an error page
func (h *Handler) HandleError(c echo.Context, err error, statusCode int, message string) error {
	return c.Render(statusCode, "error.html", map[string]interface{}{"Error": message, "Details": err.Error()})
}

// RenderEmailTemplate renders the email template
func (h *Handler) RenderEmailTemplate(c echo.Context, email, content string) error {
	data := struct {
		Email   string
		Content template.HTML
	}{
		Email:   email,
		Content: template.HTML(content),
	}

	h.logger.Debug("Data for email template", "data", data)
	err := c.Render(http.StatusOK, "email.html", map[string]interface{}{"Data": data})
	if err != nil {
		return h.errorHandler.HandleError(c, err, http.StatusInternalServerError, "Error rendering email template")
	}
	h.logger.Info("Email template rendered successfully")
	return nil
}

// HandleRepresentativeLookup handles POST requests for fetching representatives
func (h *Handler) HandleRepresentativeLookup(c echo.Context) error {
	postalCode := c.FormValue("postal_code")
	representativeType := c.FormValue("type")

	representatives, err := h.representativeLookupService.FetchRepresentatives(postalCode)
	if err != nil {
		return h.errorHandler.HandleError(c, err, http.StatusInternalServerError, "Error fetching representatives")
	}

	filters := map[string]string{
		"type": representativeType,
	}

	filteredRepresentatives := h.representativeLookupService.FilterRepresentatives(representatives, filters)

	return c.Render(http.StatusOK, "representatives.html", map[string]interface{}{
		"Representatives": filteredRepresentatives,
	})
}
