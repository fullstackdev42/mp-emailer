package campaign

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/labstack/echo/v4"
)

// CampaignGET handles GET requests for campaign details
func (h *Handler) CampaignGET(c echo.Context) error {
	h.Logger.Debug("CampaignGET: Starting")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		status, msg := h.mapError(ErrInvalidCampaignID)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}
	h.Logger.Debug("CampaignGET: Parsed ID", "id", id)
	campaignParams := GetCampaignParams{ID: id}
	campaign, err := h.service.FetchCampaign(campaignParams)
	if err != nil {
		status, msg := h.mapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
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
		status, msg := h.mapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}
	h.Logger.Debug("Rendering all campaigns", "count", len(campaigns))
	data := shared.Data{
		Title:    "All Campaigns",
		PageName: "campaigns",
		Content: map[string]interface{}{
			"Campaigns": campaigns,
		},
	}
	return c.Render(http.StatusOK, "campaigns", data)
}

// CreateCampaignForm handles GET requests for the campaign creation form
func (h *Handler) CreateCampaignForm(c echo.Context) error {
	h.Logger.Debug("Handling CreateCampaignForm request")
	return c.Render(http.StatusOK, "campaign_create", shared.Data{
		Title:    "Create Campaign",
		PageName: "campaign_create",
		Content:  nil,
	})
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

// DeleteCampaign handles DELETE requests for deleting a campaign
func (h *Handler) DeleteCampaign(c echo.Context) error {
	h.Logger.Debug("Handling DeleteCampaign request")
	params := new(DeleteCampaignDTO)
	if err := c.Bind(params); err != nil {
		status, msg := h.mapError(ErrInvalidCampaignID)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}
	if err := h.service.DeleteCampaign(DeleteCampaignDTO{ID: params.ID}); err != nil {
		status, msg := h.mapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}
	h.Logger.Info("Campaign deleted successfully", "campaignID", params.ID)
	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

// EditCampaignForm handles GET requests for the campaign edit form
func (h *Handler) EditCampaignForm(c echo.Context) error {
	h.Logger.Debug("Handling EditCampaignForm request")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		status, msg := h.mapError(ErrInvalidCampaignID)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	campaign, err := h.service.FetchCampaign(GetCampaignParams{ID: id})
	if err != nil {
		status, msg := h.mapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	return c.Render(http.StatusOK, "campaign_edit", shared.Data{
		Title:    "Edit Campaign",
		PageName: "campaign_edit",
		Content: map[string]interface{}{
			"Campaign": campaign,
		},
	})
}

// EditCampaign handles POST requests for updating a campaign
func (h *Handler) EditCampaign(c echo.Context) error {
	h.Logger.Debug("Handling EditCampaign request")
	params := EditParams{}
	var err error
	params.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		status, msg := h.mapError(ErrInvalidCampaignID)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}
	params.Name = c.FormValue("name")
	params.Template = c.FormValue("template")
	if err := h.service.UpdateCampaign(&UpdateCampaignDTO{
		ID:       params.ID,
		Name:     params.Name,
		Template: params.Template,
	}); err != nil {
		status, msg := h.mapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}
	h.Logger.Info("Campaign updated successfully", "campaignID", params.ID)
	return c.Redirect(http.StatusSeeOther, "/campaign/"+strconv.Itoa(params.ID))
}

// SendCampaign handles POST requests for sending a campaign
func (h *Handler) SendCampaign(c echo.Context) error {
	h.Logger.Info("Handling campaign submit request")
	params := new(SendCampaignParams)
	if err := c.Bind(params); err != nil {
		status, msg := h.mapError(ErrInvalidCampaignData)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}
	postalCode, err := extractAndValidatePostalCode(c)
	if err != nil {
		status, msg := h.mapError(ErrInvalidPostalCode)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}
	mp, err := h.representativeLookupService.FetchRepresentatives(postalCode)
	if err != nil || len(mp) == 0 {
		status, msg := h.mapError(ErrNoRepresentatives)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}
	campaign, err := h.service.FetchCampaign(GetCampaignParams{ID: params.ID})
	if err != nil {
		status, msg := h.mapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
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
	return c.Render(http.StatusOK, "representatives", shared.Data{
		Title:    "Representatives",
		PageName: "representatives",
		Content: map[string]interface{}{
			"Representatives": filteredRepresentatives,
		},
	})
}
