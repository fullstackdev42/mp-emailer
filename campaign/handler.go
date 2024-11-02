package campaign

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	shared.BaseHandler
	service                     ServiceInterface
	representativeLookupService RepresentativeLookupServiceInterface
	emailService                email.Service
	client                      ClientInterface
}

// NewHandler initializes a new Handler
func NewHandler(params HandlerParams) (HandlerResult, error) {
	handler := &Handler{
		BaseHandler:                 shared.NewBaseHandler(params.BaseHandlerParams),
		service:                     params.Service,
		representativeLookupService: params.RepresentativeLookupService,
		emailService:                params.EmailService,
		client:                      params.Client,
	}
	return HandlerResult{Handler: handler}, nil
}

// TemplateData provides a consistent structure for all template rendering
type TemplateData struct {
	Campaign        *Campaign
	Campaigns       []Campaign
	Email           string
	Content         template.HTML
	Error           error
	Representatives []Representative
}

// CampaignGET handles GET requests for campaign details
func (h *Handler) CampaignGET(c echo.Context) error {
	h.Logger.Debug("CampaignGET: Starting")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}
	h.Logger.Debug("CampaignGET: Parsed ID", "id", id)
	campaignParams := GetCampaignParams{ID: id}
	campaign, err := h.service.FetchCampaign(campaignParams)
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Failed to fetch campaign", http.StatusInternalServerError)
	}
	h.Logger.Debug("CampaignGET: Campaign fetched successfully", "id", id)

	data := map[string]interface{}{
		"Title":    "Campaign Details",
		"PageName": "campaign",
		"Campaign": campaign,
	}

	return c.Render(http.StatusOK, "campaign", data)
}

// GetCampaigns handles GET requests for all campaigns
func (h *Handler) GetCampaigns(c echo.Context) error {
	h.Logger.Debug("Handling GetCampaigns request")
	campaigns, err := h.service.GetCampaigns()
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Error fetching campaigns", http.StatusInternalServerError)
	}
	h.Logger.Debug("Rendering all campaigns", "count", len(campaigns))
	return h.renderTemplate(c, "campaigns", &TemplateData{
		Campaigns: campaigns,
	})
}

// CreateCampaignForm handles GET requests for the campaign creation form
func (h *Handler) CreateCampaignForm(c echo.Context) error {
	h.Logger.Debug("Handling CreateCampaignForm request")
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
	h.Logger.Debug("Handling CreateCampaign request")
	params := new(CreateCampaignParams)
	if err := c.Bind(params); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}
	ownerID, err := user.GetOwnerIDFromSession(c)
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Unauthorized", http.StatusUnauthorized)
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
		return h.ErrorHandler.HandleHTTPError(c, err, "Error creating campaign", http.StatusInternalServerError)
	}
	h.Logger.Info("Campaign created successfully", "campaignID", campaign.ID, "ownerID", ownerID)
	return c.Redirect(http.StatusSeeOther, "/campaign/"+strconv.Itoa(campaign.ID))
}

// DeleteCampaignParams defines the parameters for deleting a campaign
type DeleteCampaignParams struct {
	ID int `param:"id"`
}

// DeleteCampaign handles DELETE requests for deleting a campaign
func (h *Handler) DeleteCampaign(c echo.Context) error {
	h.Logger.Debug("Handling DeleteCampaign request")
	params := new(DeleteCampaignParams)
	if err := c.Bind(params); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}
	if err := h.service.DeleteCampaign(DeleteCampaignParams{ID: params.ID}); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Error deleting campaign", http.StatusInternalServerError)
	}
	h.Logger.Info("Campaign deleted successfully", "campaignID", params.ID)
	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

// EditCampaignForm handles GET requests for the campaign edit form
func (h *Handler) EditCampaignForm(c echo.Context) error {
	h.Logger.Debug("Handling EditCampaignForm request")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}
	campaign, err := h.service.FetchCampaign(GetCampaignParams{ID: id})
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Error fetching campaign", http.StatusInternalServerError)
	}
	return h.renderTemplate(c, "campaign_edit", &TemplateData{Campaign: campaign})
}

// EditParams defines the parameters for editing a campaign
type EditParams struct {
	ID       int
	Name     string
	Template string
}

// EditCampaign handles POST requests for updating a campaign
func (h *Handler) EditCampaign(c echo.Context) error {
	h.Logger.Debug("Handling EditCampaign request")
	params := EditParams{}
	var err error
	params.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}
	params.Name = c.FormValue("name")
	params.Template = c.FormValue("template")
	if err := h.service.UpdateCampaign(&UpdateCampaignDTO{
		ID:       params.ID,
		Name:     params.Name,
		Template: params.Template,
	}); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Error updating campaign", http.StatusInternalServerError)
	}
	h.Logger.Info("Campaign updated successfully", "campaignID", params.ID)
	return c.Redirect(http.StatusSeeOther, "/campaign/"+strconv.Itoa(params.ID))
}

// SendCampaignParams defines the parameters for sending a campaign
type SendCampaignParams struct {
	ID         int    `param:"id"`
	PostalCode string `form:"postal_code"`
}

// SendCampaign handles POST requests for sending a campaign
func (h *Handler) SendCampaign(c echo.Context) error {
	h.Logger.Info("Handling campaign submit request")
	params := new(SendCampaignParams)
	if err := c.Bind(params); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}
	postalCode, err := extractAndValidatePostalCode(c)
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid postal code", http.StatusBadRequest)
	}
	mp, err := h.representativeLookupService.FetchRepresentatives(postalCode)
	if err != nil || len(mp) == 0 {
		return h.ErrorHandler.HandleHTTPError(c, err, "No representatives found", http.StatusNotFound)
	}
	campaign, err := h.service.FetchCampaign(GetCampaignParams{ID: params.ID})
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid campaign ID", http.StatusBadRequest)
	}
	userData := extractUserData(c)
	representative := mp[0]
	emailContent := h.service.ComposeEmail(ComposeEmailParams{
		MP:       representative,
		Campaign: campaign,
		UserData: userData,
	})
	h.Logger.Info("Email composed successfully", "campaignID", params.ID, "representative", representative.Email)
	return h.RenderEmailTemplate(c, representative.Email, emailContent)
}

// RenderEmailTemplate renders the email template
func (h *Handler) RenderEmailTemplate(c echo.Context, email string, content string) error {
	h.Logger.Debug("Rendering email template", "recipientEmail", email)

	// Convert content to template.HTML once
	htmlContent := template.HTML(content)

	data := map[string]interface{}{
		"Title":      "Email Preview",
		"PageName":   "email",
		"Email":      email,
		"Content":    htmlContent,
		"RawContent": content, // Add raw content for mailto link
	}

	return c.Render(http.StatusOK, "email", data)
}

// renderTemplate renders the specified template with strongly typed data
func (h *Handler) renderTemplate(c echo.Context, name string, data *TemplateData) error {
	if err := h.TemplateRenderer.Render(c.Response().Writer, name, data, c); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Failed to render "+name, http.StatusInternalServerError)
	}
	h.Logger.Debug(name + " template rendered successfully")
	return nil
}

// HandleRepresentativeLookup handles POST requests for fetching representatives
func (h *Handler) HandleRepresentativeLookup(c echo.Context) error {
	h.Logger.Debug("Handling representative lookup request")
	postalCode := c.FormValue("postal_code")
	representativeType := c.FormValue("type")
	representatives, err := h.representativeLookupService.FetchRepresentatives(postalCode)
	if err != nil {
		h.Logger.Error("Error fetching representatives", err, "postalCode", postalCode)
		return h.ErrorHandler.HandleHTTPError(c, err, "Error fetching representatives", http.StatusInternalServerError)
	}
	filters := map[string]string{"type": representativeType}
	filteredRepresentatives := h.representativeLookupService.FilterRepresentatives(representatives, filters)
	h.Logger.Info("Representatives lookup successful", "count", len(filteredRepresentatives), "postalCode", postalCode, "type", representativeType)
	return h.renderTemplate(c, "representatives", &TemplateData{
		Representatives: filteredRepresentatives,
	})
}
