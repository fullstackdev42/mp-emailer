package campaign

import (
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
	templateRenderer            shared.TemplateRenderer
}

// CampaignGET handles GET requests for campaign details
func (h *Handler) CampaignGET(c echo.Context) error {
	h.logger.Debug("CampaignGET: Starting")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}
	h.logger.Debug("CampaignGET: Parsed ID", "id", id)
	campaignParams := GetCampaignParams{ID: id}
	campaign, err := h.service.FetchCampaign(campaignParams)
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Failed to fetch campaign", http.StatusInternalServerError)
	}
	h.logger.Debug("CampaignGET: Campaign fetched successfully", "id", id)
	return h.renderTemplate(c, "campaign", map[string]interface{}{"Campaign": campaign})
}

// GetAllCampaigns handles GET requests for all campaigns
func (h *Handler) GetAllCampaigns(c echo.Context) error {
	h.logger.Debug("Handling GetAllCampaigns request")
	campaigns, err := h.service.GetAllCampaigns()
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error fetching campaigns", http.StatusInternalServerError)
	}
	isAuthenticated, _ := c.Get("isAuthenticated").(bool)
	h.logger.Debug("Rendering all campaigns", "count", len(campaigns))
	return h.renderTemplate(c, "campaigns", map[string]interface{}{
		"Campaigns":       campaigns,
		"IsAuthenticated": isAuthenticated,
	})
}

// CreateCampaignForm handles GET requests for the campaign creation form
func (h *Handler) CreateCampaignForm(c echo.Context) error {
	h.logger.Debug("Handling CreateCampaignForm request")
	return h.renderTemplate(c, "campaign_create", nil)
}

// CreateCampaignParams defines the parameters for creating a campaign
type CreateCampaignParams struct {
	Name        string `form:"name"`
	Description string `form:"description"`
	Template    string `form:"template"`
	OwnerID     string // This will be set from the session
}

// CreateCampaign handles POST requests for creating a new campaign
func (h *Handler) CreateCampaign(c echo.Context) error {
	h.logger.Debug("Handling CreateCampaign request")
	params := new(CreateCampaignParams)
	if err := c.Bind(params); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}
	ownerID, err := user.GetOwnerIDFromSession(c)
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Unauthorized", http.StatusUnauthorized)
	}
	params.OwnerID = ownerID
	dto := &CreateCampaignDTO{
		Name:        params.Name,
		Description: params.Description,
		Template:    params.Template,
		OwnerID:     params.OwnerID,
	}
	campaign, err := h.service.CreateCampaign(dto)
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error creating campaign", http.StatusInternalServerError)
	}
	h.logger.Info("Campaign created successfully", "campaignID", campaign.ID, "ownerID", ownerID)
	return c.Redirect(http.StatusSeeOther, "/campaign")
}

// DeleteCampaignParams defines the parameters for deleting a campaign
type DeleteCampaignParams struct {
	ID int `param:"id"`
}

// DeleteCampaign handles DELETE requests for deleting a campaign
func (h *Handler) DeleteCampaign(c echo.Context) error {
	h.logger.Debug("Handling DeleteCampaign request")
	params := new(DeleteCampaignParams)
	if err := c.Bind(params); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}
	if err := h.service.DeleteCampaign(DeleteCampaignParams{ID: params.ID}); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error deleting campaign", http.StatusInternalServerError)
	}
	h.logger.Info("Campaign deleted successfully", "campaignID", params.ID)
	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

// EditCampaignForm handles GET requests for the campaign edit form
func (h *Handler) EditCampaignForm(c echo.Context) error {
	h.logger.Debug("Handling EditCampaignForm request")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}
	campaign, err := h.service.FetchCampaign(GetCampaignParams{ID: id})
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error fetching campaign", http.StatusInternalServerError)
	}
	return h.renderTemplate(c, "campaign_edit", map[string]interface{}{"Campaign": campaign})
}

// EditParams defines the parameters for editing a campaign
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
		return h.errorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}
	params.Name = c.FormValue("name")
	params.Template = c.FormValue("template")
	if err := h.service.UpdateCampaign(&UpdateCampaignDTO{
		ID:       params.ID,
		Name:     params.Name,
		Template: params.Template,
	}); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error updating campaign", http.StatusInternalServerError)
	}
	h.logger.Info("Campaign updated successfully", "campaignID", params.ID)
	return c.Redirect(http.StatusSeeOther, "/campaigns/"+strconv.Itoa(params.ID))
}

// SendCampaignParams defines the parameters for sending a campaign
type SendCampaignParams struct {
	ID         int    `param:"id"`
	PostalCode string `form:"postal_code"`
}

// SendCampaign handles POST requests for sending a campaign
func (h *Handler) SendCampaign(c echo.Context) error {
	h.logger.Info("Handling campaign submit request")
	params := new(SendCampaignParams)
	if err := c.Bind(params); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}
	postalCode, err := extractAndValidatePostalCode(c)
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid postal code", http.StatusBadRequest)
	}
	mp, err := h.representativeLookupService.FetchRepresentatives(postalCode)
	if err != nil || len(mp) == 0 {
		return h.errorHandler.HandleHTTPError(c, err, "No representatives found", http.StatusNotFound)
	}
	campaign, err := h.service.FetchCampaign(GetCampaignParams{ID: params.ID})
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}
	userData := extractUserData(c)
	representative := mp[0]
	emailContent := h.service.ComposeEmail(ComposeEmailParams{
		MP:       representative,
		Campaign: campaign,
		UserData: userData,
	})
	h.logger.Info("Email composed successfully", "campaignID", params.ID, "representative", representative.Email)
	return h.RenderEmailTemplate(c, representative.Email, emailContent)
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
	return h.renderTemplate(c, "email", data)
}

// renderTemplate renders the specified template with the given data
func (h *Handler) renderTemplate(c echo.Context, name string, data interface{}) error {
	if err := h.templateRenderer.Render(c.Response().Writer, name, data, c); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Failed to render "+name, http.StatusInternalServerError)
	}
	h.logger.Debug(name + " template rendered successfully")
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
	filters := map[string]string{"type": representativeType}
	filteredRepresentatives := h.representativeLookupService.FilterRepresentatives(representatives, filters)
	h.logger.Info("Representatives lookup successful", "count", len(filteredRepresentatives), "postalCode", postalCode, "type", representativeType)
	return h.renderTemplate(c, "representatives", map[string]interface{}{
		"Representatives": filteredRepresentatives,
	})
}
