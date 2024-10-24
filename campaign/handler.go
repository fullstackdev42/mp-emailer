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

// Define a DTO for returning campaign details
type DetailDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PostalCode  string `json:"postal_code"`
	Template    string `json:"template"`
	OwnerID     string `json:"owner_id"`
}

// CampaignGET handles GET requests for campaign details
func (h *Handler) CampaignGET(c echo.Context) error {
	h.logger.Debug("Handling CampaignGET request")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid campaign ID", err, "input", c.Param("id"))
		return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
			"Title":   "Bad Request",
			"Message": "Invalid campaign ID",
		})
	}

	campaign, err := h.service.FetchCampaign(id)
	if err != nil {
		if errors.Is(err, ErrCampaignNotFound) {
			h.logger.Info("Campaign not found", "campaignID", id)
			return c.Render(http.StatusNotFound, "error.html", map[string]interface{}{
				"Title":   "Not Found",
				"Message": "Campaign not found",
			})
		}
		h.logger.Error("Error fetching campaign", err, "campaignID", id)
		return c.Render(http.StatusInternalServerError, "error.html", map[string]interface{}{
			"Title":   "Internal Server Error",
			"Message": "An error occurred while fetching the campaign",
		})
	}

	dto := DetailDTO{
		ID:          campaign.ID,
		Name:        campaign.Name,
		Description: campaign.Description,
		PostalCode:  campaign.PostalCode,
		Template:    campaign.Template,
		OwnerID:     campaign.OwnerID,
	}

	h.logger.Debug("Rendering campaign details", "campaignID", id)
	return c.Render(http.StatusOK, "campaign_details.html", shared.PageData{
		Title:   "Campaign Details",
		Content: dto,
	})
}

// GetAllCampaigns handles GET requests for all campaigns
func (h *Handler) GetAllCampaigns(c echo.Context) error {
	h.logger.Debug("Handling GetAllCampaigns request")
	campaigns, err := h.service.GetAllCampaigns()
	if err != nil {
		h.logger.Error("Error fetching campaigns", err)
		return h.errorHandler.HandleHTTPError(c, err, "Error fetching campaigns", http.StatusInternalServerError)
	}
	h.logger.Debug("Rendering all campaigns", "count", len(campaigns))

	// Add this error logging
	err = c.Render(http.StatusOK, "campaigns.html", map[string]interface{}{"Campaigns": campaigns})
	if err != nil {
		h.logger.Error("Error rendering campaigns template", err)
		return h.errorHandler.HandleHTTPError(c, err, "Error rendering campaigns", http.StatusInternalServerError)
	}

	return nil
}

// CreateCampaignForm handles GET requests for the campaign creation form
func (h *Handler) CreateCampaignForm(c echo.Context) error {
	h.logger.Debug("Handling CreateCampaignForm request")
	return c.Render(http.StatusOK, "campaign_create.html", nil)
}

// CreateCampaign handles POST requests for creating a new campaign
func (h *Handler) CreateCampaign(c echo.Context) error {
	h.logger.Debug("Handling CreateCampaign request")
	dto := new(CreateCampaignDTO)
	if err := c.Bind(dto); err != nil {
		h.logger.Warn("Invalid input for campaign creation", err)
		return h.errorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}

	ownerID, err := user.GetOwnerIDFromSession(c)
	if err != nil {
		h.logger.Warn("Unauthorized access attempt", err)
		return h.errorHandler.HandleHTTPError(c, err, "Unauthorized", http.StatusUnauthorized)
	}

	dto.OwnerID = ownerID

	if err := h.service.CreateCampaign(dto); err != nil {
		h.logger.Error("Error creating campaign", err)
		return h.errorHandler.HandleHTTPError(c, err, "Error creating campaign", http.StatusInternalServerError)
	}

	h.logger.Info("Campaign created successfully", "ownerID", ownerID)
	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

// DeleteCampaign handles DELETE requests for deleting a campaign
func (h *Handler) DeleteCampaign(c echo.Context) error {
	h.logger.Debug("Handling DeleteCampaign request")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid campaign ID for deletion", err, "input", c.Param("id"))
		return h.errorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}
	if err := h.service.DeleteCampaign(id); err != nil {
		h.logger.Error("Error deleting campaign", err, "campaignID", id)
		return h.errorHandler.HandleHTTPError(c, err, "Error deleting campaign", http.StatusInternalServerError)
	}
	h.logger.Info("Campaign deleted successfully", "campaignID", id)
	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

// EditCampaignForm handles GET requests for the campaign edit form
func (h *Handler) EditCampaignForm(c echo.Context) error {
	h.logger.Debug("Handling EditCampaignForm request")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid campaign ID for edit form", err, "input", c.Param("id"))
		return h.errorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}
	campaign, err := h.service.FetchCampaign(id)
	if err != nil {
		h.logger.Error("Error fetching campaign for edit", err, "campaignID", id)
		return h.errorHandler.HandleHTTPError(c, err, "Error fetching campaign", http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "campaign_edit.html", map[string]interface{}{
		"Campaign": campaign,
	})
}

type EditParams struct {
	ID       int
	Name     string
	Template string
}

// EditCampaign handles POST requests for updating a campaign
func (h *Handler) EditCampaign(c echo.Context) error {
	h.logger.Debug("Handling EditCampaign request")
	params := EditParams{}
	var err error
	params.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid campaign ID for edit", err, "input", c.Param("id"))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}
	params.Name = c.FormValue("name")
	params.Template = c.FormValue("template")

	if err := h.service.UpdateCampaign(&Campaign{
		ID:       params.ID,
		Name:     params.Name,
		Template: params.Template,
	}); err != nil {
		h.logger.Error("Error updating campaign", err, "campaignID", params.ID)
		return h.errorHandler.HandleHTTPError(c, err, "Error updating campaign", http.StatusInternalServerError)
	}
	h.logger.Info("Campaign updated successfully", "campaignID", params.ID)
	return c.Redirect(http.StatusSeeOther, "/campaigns/"+strconv.Itoa(params.ID))
}

// SendCampaign handles POST requests for sending a campaign
func (h *Handler) SendCampaign(c echo.Context) error {
	h.logger.Info("Handling campaign submit request")
	postalCode, err := extractAndValidatePostalCode(c)
	if err != nil {
		h.logger.Warn("Invalid postal code submitted", err)
		return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
			"Error": "Invalid postal code",
		})
	}
	mp, err := h.representativeLookupService.FetchRepresentatives(postalCode)
	if err != nil {
		h.logger.Error("Error finding MP", err, "postalCode", postalCode)
		return h.errorHandler.HandleHTTPError(c, err, "Error finding MP", http.StatusInternalServerError)
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid campaign ID for send", err, "input", c.Param("id"))
		return h.errorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}
	campaign, err := h.service.FetchCampaign(id)
	if err != nil {
		h.logger.Error("Error fetching campaign for send", err, "campaignID", id)
		return h.errorHandler.HandleHTTPError(c, err, "Error fetching campaign", http.StatusInternalServerError)
	}
	userData := extractUserData(c)

	if len(mp) == 0 {
		h.logger.Warn("No representatives found", "postalCode", postalCode)
		return h.errorHandler.HandleHTTPError(c, errors.New("no representatives found"), "No representatives found", http.StatusNotFound)
	}
	representative := mp[0]

	emailContent := h.service.ComposeEmail(representative, campaign, userData)
	h.logger.Info("Email composed successfully", "campaignID", id, "representative", representative.Email)
	return h.RenderEmailTemplate(c, representative.Email, emailContent)
}

// HandleError renders an error page
func (h *Handler) HandleError(c echo.Context, err error, statusCode int, message string) error {
	return c.Render(statusCode, "error.html", map[string]interface{}{"Error": message, "Details": err.Error()})
}

// RenderEmailTemplate renders the email template
func (h *Handler) RenderEmailTemplate(c echo.Context, email, content string) error {
	h.logger.Debug("Rendering email template", "recipientEmail", email)
	data := struct {
		Email   string
		Content template.HTML
	}{
		Email:   email,
		Content: template.HTML(content),
	}

	err := c.Render(http.StatusOK, "email.html", map[string]interface{}{"Data": data})
	if err != nil {
		h.logger.Error("Error rendering email template", err)
		return h.errorHandler.HandleHTTPError(c, err, "Error rendering email template", http.StatusInternalServerError)
	}
	h.logger.Info("Email template rendered successfully")
	return nil
}

// HandleRepresentativeLookup handles POST requests for fetching representatives
func (h *Handler) HandleRepresentativeLookup(c echo.Context) error {
	h.logger.Debug("Handling representative lookup request")
	postalCode := c.FormValue("postal_code")
	representativeType := c.FormValue("type")

	representatives, err := h.representativeLookupService.FetchRepresentatives(postalCode)
	if err != nil {
		h.logger.Error("Error fetching representatives", err, "postalCode", postalCode)
		return h.errorHandler.HandleHTTPError(c, err, "Error fetching representatives", http.StatusInternalServerError)
	}

	filters := map[string]string{
		"type": representativeType,
	}

	filteredRepresentatives := h.representativeLookupService.FilterRepresentatives(representatives, filters)

	h.logger.Info("Representatives lookup successful", "count", len(filteredRepresentatives), "postalCode", postalCode, "type", representativeType)
	return c.Render(http.StatusOK, "representatives.html", map[string]interface{}{
		"Representatives": filteredRepresentatives,
	})
}
