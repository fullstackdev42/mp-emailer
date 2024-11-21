package campaign

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jonesrussell/mp-emailer/email"
	"github.com/jonesrussell/mp-emailer/session"
	"github.com/jonesrussell/mp-emailer/shared"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type Handler struct {
	shared.BaseHandler
	service                     ServiceInterface
	representativeLookupService RepresentativeLookupServiceInterface
	emailService                email.Service
	client                      ClientInterface
}

// HandlerParams for dependency injection
type HandlerParams struct {
	shared.BaseHandlerParams
	fx.In
	Service                     ServiceInterface
	RepresentativeLookupService RepresentativeLookupServiceInterface
	EmailService                email.Service
	Client                      ClientInterface
}

// HandlerResult is the output struct for NewHandler
type HandlerResult struct {
	fx.Out
	Handler *Handler
}

// NewHandler initializes a new Handler
func NewHandler(params HandlerParams) (HandlerResult, error) {
	base := shared.NewBaseHandler(params.BaseHandlerParams)
	base.MapError = mapErrorToHTTPStatus

	handler := &Handler{
		BaseHandler:                 base,
		service:                     params.Service,
		representativeLookupService: params.RepresentativeLookupService,
		emailService:                params.EmailService,
		client:                      params.Client,
	}
	return HandlerResult{Handler: handler}, nil
}

// CampaignGET handles GET requests for campaign details
func (h *Handler) CampaignGET(c echo.Context) error {
	h.Logger.Debug("CampaignGET: Starting")
	id := c.Param("id")
	campaignID, err := uuid.Parse(id)
	if err != nil {
		status, msg := h.MapError(ErrInvalidCampaignID)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	campaign, err := h.service.FetchCampaign(c.Request().Context(), GetCampaignParams{ID: campaignID})
	if err != nil {
		status, msg := h.MapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	// Get userID and check authentication
	userID, err := h.GetUserIDFromSession(c)
	isAuthenticated := err == nil && userID != ""

	// Optional: Add ownership check
	if isAuthenticated && campaign.OwnerID.String() != userID {
		return h.ErrorHandler.HandleHTTPError(c, ErrUnauthorizedAccess, "Unauthorized", http.StatusUnauthorized)
	}

	data := shared.Data{
		Title:           "Campaign Details",
		PageName:        "campaign",
		IsAuthenticated: isAuthenticated,
		Content: map[string]interface{}{
			"Campaign": campaign,
		},
	}

	return c.Render(http.StatusOK, "campaign", data)
}

// GetCampaigns handles GET requests for all campaigns
func (h *Handler) GetCampaigns(c echo.Context) error {
	h.Logger.Debug("Handling GetCampaigns request")
	campaigns, err := h.service.GetCampaigns(c.Request().Context())
	if err != nil {
		status, msg := h.MapError(err)
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

	userID, err := h.GetUserIDFromSession(c)
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Unauthorized", http.StatusUnauthorized)
	}

	params := &CreateCampaignParams{
		Name:        strings.TrimSpace(c.FormValue("name")),
		Description: strings.TrimSpace(c.FormValue("description")),
		Template:    strings.TrimSpace(c.FormValue("template")),
		OwnerID:     uuid.Must(uuid.Parse(userID)),
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

		return c.Render(http.StatusBadRequest, "campaign_create", shared.Data{
			Title:    "Create Campaign",
			PageName: "campaign_create",
			Content: map[string]interface{}{
				"Errors":     validationErrors,
				"FormValues": params,
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

	// Add context from echo.Context
	campaign, err := h.service.CreateCampaign(c.Request().Context(), dto)
	if err != nil {
		h.Logger.Error("CreateCampaign: Failed to create campaign", err,
			"ownerID", userID,
			"name", params.Name)
		status, msg := h.MapError(err)
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

	userID, err := h.GetUserIDFromSession(c)
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Unauthorized", http.StatusUnauthorized)
	}

	id := c.Param("id")
	campaignID, err := uuid.Parse(id)
	if err != nil {
		status, msg := h.MapError(ErrInvalidCampaignID)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	// Verify ownership
	campaign, err := h.service.FetchCampaign(c.Request().Context(), GetCampaignParams{ID: campaignID})
	if err != nil {
		status, msg := h.MapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	if campaign.OwnerID.String() != userID {
		return h.ErrorHandler.HandleHTTPError(c, ErrUnauthorizedAccess, "Unauthorized", http.StatusUnauthorized)
	}

	if err := h.service.DeleteCampaign(c.Request().Context(), DeleteCampaignDTO{ID: campaignID}); err != nil {
		status, msg := h.MapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	h.Logger.Info("Campaign deleted successfully", "campaignID", campaignID)
	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

// EditCampaignForm handles GET requests for the campaign edit form
func (h *Handler) EditCampaignForm(c echo.Context) error {
	h.Logger.Debug("Handling EditCampaignForm request")
	id := c.Param("id")
	campaignID, err := uuid.Parse(id)
	if err != nil {
		status, msg := h.MapError(ErrInvalidCampaignID)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	campaign, err := h.service.FetchCampaign(c.Request().Context(), GetCampaignParams{ID: campaignID})
	if err != nil {
		if err == ErrCampaignNotFound {
			return h.ErrorHandler.HandleHTTPError(c, err, "Campaign not found", http.StatusNotFound)
		}
		status, msg := h.MapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	userID, err := h.GetUserIDFromSession(c)
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Unauthorized", http.StatusUnauthorized)
	}

	if campaign.OwnerID.String() != userID {
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

	userID, err := h.GetUserIDFromSession(c)
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Unauthorized", http.StatusUnauthorized)
	}

	// Get campaign ID from URL parameter
	id := c.Param("id")
	campaignID, err := uuid.Parse(id)
	if err != nil {
		status, msg := h.MapError(ErrInvalidCampaignID)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	// Verify ownership
	campaign, err := h.service.FetchCampaign(c.Request().Context(), GetCampaignParams{ID: campaignID})
	if err != nil {
		status, msg := h.MapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	if campaign.OwnerID.String() != userID {
		return h.ErrorHandler.HandleHTTPError(c, ErrUnauthorizedAccess, "Unauthorized", http.StatusUnauthorized)
	}

	params := EditParams{
		ID:       campaignID,
		Name:     c.FormValue("name"),
		Template: c.FormValue("template"),
	}

	if err := h.service.UpdateCampaign(c.Request().Context(), &UpdateCampaignDTO{
		ID:       params.ID,
		Name:     params.Name,
		Template: params.Template,
	}); err != nil {
		status, msg := h.MapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	if err := h.addFlashMessage(c, "Campaign updated successfully"); err != nil {
		h.Logger.Error("Failed to add flash message", err)
	}

	h.Logger.Info("Campaign updated successfully", "campaignID", params.ID)
	return c.Redirect(http.StatusSeeOther, "/campaign/"+params.ID.String())
}

// ComposeEmail handles the initial postal code submission and email composition
func (h *Handler) ComposeEmail(c echo.Context) error {
	h.Logger.Info("Handling email composition request")

	params := new(SendCampaignParams)
	if err := c.Bind(params); err != nil {
		status, msg := h.MapError(ErrInvalidCampaignData)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	postalCode, err := extractAndValidatePostalCode(c)
	if err != nil {
		status, msg := h.MapError(ErrInvalidPostalCode)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	mp, err := h.representativeLookupService.FetchRepresentatives(postalCode)
	if err != nil || len(mp) == 0 {
		status, msg := h.MapError(ErrNoRepresentatives)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	campaign, err := h.service.FetchCampaign(c.Request().Context(), GetCampaignParams{ID: params.ID})
	if err != nil {
		status, msg := h.MapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	userData := extractUserData(c)
	representative := mp[0]
	emailContent, err := h.service.ComposeEmail(c.Request().Context(), ComposeEmailParams{
		MP:       representative,
		Campaign: campaign,
		UserData: userData,
	})
	if err != nil {
		status, msg := h.MapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

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
		status, msg := h.MapError(err)
		return h.ErrorHandler.HandleHTTPError(c, err, msg, status)
	}

	h.Logger.Info("Email sent successfully",
		"recipient", email,
		"campaignID", c.Param("id"))

	// Add success flash message
	session, err := h.Store.Get(c.Request(), h.Config.Auth.SessionName)
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

// GetSessionManager retrieves the session manager from context
func (h *Handler) GetSessionManager(c echo.Context) (session.Manager, error) {
	sessionManager, ok := c.Get("session_manager").(session.Manager)
	if !ok {
		return nil, errors.New("session manager not found")
	}
	return sessionManager, nil
}

// GetUserIDFromSession retrieves and validates the user ID from session
func (h *Handler) GetUserIDFromSession(c echo.Context) (string, error) {
	sessionManager, err := h.GetSessionManager(c)
	if err != nil {
		return "", err
	}

	sess, err := sessionManager.GetSession(c, h.Config.Auth.SessionName)
	if err != nil {
		h.Logger.Error("Failed to get session", err)
		return "", err
	}

	userID, ok := sessionManager.GetSessionValue(sess, "user_id").(string)
	if !ok || userID == "" {
		return "", ErrUnauthorizedAccess
	}

	return userID, nil
}

// addFlashMessage adds a flash message to the session
func (h *Handler) addFlashMessage(c echo.Context, message string) error {
	session, err := h.Store.Get(c.Request(), h.Config.Auth.SessionName)
	if err != nil {
		return err
	}

	session.AddFlash(message, "messages")
	return session.Save(c.Request(), c.Response().Writer)
}
