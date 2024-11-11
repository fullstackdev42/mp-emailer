package campaign

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/fullstackdev42/mp-emailer/middleware"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// CampaignGET handles GET requests for campaign details
func (h *Handler) CampaignGET(c echo.Context) error {
	h.Logger.Debug("CampaignGET: Starting")
	id := c.Param("id")
	campaignID, err := uuid.Parse(id)
	if err != nil {
		status, msg := h.mapError(ErrInvalidCampaignID)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}
	h.Logger.Debug("CampaignGET: Parsed ID", "id", campaignID)
	campaignParams := GetCampaignParams{ID: campaignID}
	campaign, err := h.service.FetchCampaign(campaignParams)
	if err != nil {
		status, msg := h.mapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}
	h.Logger.Debug("CampaignGET: Campaign fetched successfully", "id", campaignID)

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
		Title:    "Campaigns",
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
	h.Logger.Debug("CreateCampaign: Starting")

	// Get the middleware manager from context
	manager, ok := c.Get("middleware_manager").(*middleware.Manager)
	if !ok {
		return h.ErrorHandler.HandleHTTPError(c, errors.New("middleware manager not found"), "Internal server error", http.StatusInternalServerError)
	}

	userID, err := manager.GetOwnerIDFromSession(c)
	if err != nil {
		h.Logger.Error("CreateCampaign: Failed to get owner ID from session", err)
		status, msg := h.mapError(ErrUnauthorizedAccess)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	params := &CreateCampaignParams{
		Name:        strings.TrimSpace(c.FormValue("name")),
		Description: strings.TrimSpace(c.FormValue("description")),
		Template:    strings.TrimSpace(c.FormValue("template")),
		OwnerID:     uuid.Must(uuid.Parse(userID)), // Convert string to UUID
	}

	// Enhanced validation with specific error messages
	var validationErrors []string
	if params.Name == "" {
		validationErrors = append(validationErrors, "Name is required")
	}
	if params.Description == "" {
		validationErrors = append(validationErrors, "Description is required")
	}
	if params.Template == "" {
		validationErrors = append(validationErrors, "Template is required")
	}

	if len(validationErrors) > 0 {
		h.Logger.Error("CreateCampaign: Validation failed", nil,
			"errors", strings.Join(validationErrors, ", "),
			"name", params.Name,
			"description", params.Description)

		// Return to form with error messages
		return c.Render(http.StatusBadRequest, "campaign_create", shared.Data{
			Title:    "Create Campaign",
			PageName: "campaign_create",
			Content: map[string]interface{}{
				"Errors":     validationErrors,
				"FormValues": params, // Preserve form values
			},
		})
	}

	h.Logger.Debug("CreateCampaign: Creating campaign", "ownerID", userID)

	dto := &CreateCampaignDTO{
		Name:        params.Name,
		Description: params.Description,
		Template:    params.Template,
		OwnerID:     params.OwnerID,
	}

	campaign, err := h.service.CreateCampaign(dto)
	if err != nil {
		h.Logger.Error("CreateCampaign: Failed to create campaign", err,
			"ownerID", userID,
			"name", params.Name)
		status, msg := h.mapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	h.Logger.Info("CreateCampaign: Campaign created successfully",
		"campaignID", campaign.ID,
		"ownerID", userID)

	return c.Redirect(http.StatusSeeOther, "/campaign/"+campaign.ID.String())
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
	id := c.Param("id")
	campaignID, err := uuid.Parse(id)
	if err != nil {
		status, msg := h.mapError(ErrInvalidCampaignID)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	campaign, err := h.service.FetchCampaign(GetCampaignParams{ID: campaignID})
	if err != nil {
		if err == ErrCampaignNotFound {
			return h.ErrorHandler.HandleHTTPError(c, err, "Campaign not found", http.StatusNotFound)
		}
		status, msg := h.mapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	// Get the middleware manager from context
	manager, ok := c.Get("middleware_manager").(*middleware.Manager)
	if !ok {
		return h.ErrorHandler.HandleHTTPError(c, errors.New("middleware manager not found"), "Internal server error", http.StatusInternalServerError)
	}

	ownerID, err := manager.GetOwnerIDFromSession(c)
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Unauthorized", http.StatusUnauthorized)
	}

	if campaign.OwnerID.String() != ownerID {
		return h.ErrorHandler.HandleHTTPError(c, errors.New("unauthorized"), "Unauthorized", http.StatusUnauthorized)
	}

	return c.Render(http.StatusOK, "campaign_edit", shared.Data{
		Title:    "Edit Campaign",
		PageName: "campaign_edit",
		Content: map[string]interface{}{
			"Campaign": campaign,
		},
	})
}

// EditCampaign handles PUT/POST requests for updating a campaign
func (h *Handler) EditCampaign(c echo.Context) error {
	h.Logger.Debug("Handling EditCampaign request")

	// Get campaign ID from URL parameter
	params := EditParams{}
	id := c.Param("id")
	campaignID, err := uuid.Parse(id)
	if err != nil {
		status, msg := h.mapError(ErrInvalidCampaignID)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}
	params.ID = campaignID

	// Get form values
	params.Name = c.FormValue("name")
	params.Template = c.FormValue("template")

	// Update campaign
	if err := h.service.UpdateCampaign(&UpdateCampaignDTO{
		ID:       params.ID,
		Name:     params.Name,
		Template: params.Template,
	}); err != nil {
		status, msg := h.mapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	// Add success flash message
	session, err := h.Store.Get(c.Request(), h.Config.SessionName)
	if err != nil {
		h.Logger.Error("Failed to get session", err)
		return h.ErrorHandler.HandleHTTPError(c, err, "Session error", http.StatusInternalServerError)
	}

	session.AddFlash("Campaign updated successfully", "messages")
	if err := session.Save(c.Request(), c.Response().Writer); err != nil {
		h.Logger.Error("Failed to save session", err)
		return h.ErrorHandler.HandleHTTPError(c, err, "Session error", http.StatusInternalServerError)
	}

	h.Logger.Info("Campaign updated successfully", "campaignID", params.ID)

	// Redirect to campaign details page
	return c.Redirect(http.StatusSeeOther, "/campaign/"+params.ID.String())
}

// ComposeEmail handles the initial postal code submission and email composition
func (h *Handler) ComposeEmail(c echo.Context) error {
	h.Logger.Info("Handling email composition request")

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

	h.Logger.Info("Email composed successfully",
		"campaignID", params.ID,
		"representative", representative.Email)

	return h.RenderEmailTemplate(c, representative.Email, emailContent)
}

// SendCampaign handles the actual email sending
func (h *Handler) SendCampaign(c echo.Context) error {
	h.Logger.Info("Handling email send request")

	email := c.FormValue("email")
	content := template.HTML(c.FormValue("content"))

	if email == "" || content == "" {
		h.Logger.Error("Missing required fields", nil,
			"email", email != "",
			"hasContent", content != "")
		return h.ErrorHandler.HandleHTTPError(c,
			ErrInvalidCampaignData,
			"Email and content are required",
			http.StatusBadRequest)
	}

	htmlContent := string(content)
	err := h.emailService.SendEmail(email, "Campaign", htmlContent, true)
	if err != nil {
		h.Logger.Error("Failed to send email", err,
			"recipient", email)
		status, msg := h.mapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	h.Logger.Info("Email sent successfully",
		"recipient", email,
		"campaignID", c.Param("id"))

	// Add success flash message
	session, err := h.Store.Get(c.Request(), h.Config.SessionName)
	if err == nil {
		session.AddFlash("Email sent successfully!", "messages")
		_ = session.Save(c.Request(), c.Response().Writer)
	}

	return c.Redirect(http.StatusSeeOther, "/campaign/"+c.Param("id"))
}

// RenderEmailTemplate renders the email template
func (h *Handler) RenderEmailTemplate(c echo.Context, email string, content string) error {
	h.Logger.Debug("Rendering email template", "recipientEmail", email)

	campaignID := c.Param("id")

	data := shared.Data{
		Title:    "Email Preview",
		PageName: "email",
		Content: map[string]interface{}{
			"Email":      email,
			"Content":    template.HTML(content),
			"CampaignID": campaignID,
		},
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
