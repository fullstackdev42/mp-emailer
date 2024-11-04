package campaign

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// RegisterRoutes registers the campaign routes
func RegisterRoutes(h *Handler, e *echo.Echo, cfg *config.Config) {
	// Public routes (no authentication required)
	e.GET("/campaigns", h.GetCampaigns)
	e.GET("/campaign/:id", h.CampaignGET)
	e.POST("/campaign/:id/send", h.SendCampaign)

	// Protected routes (require authentication)
	protected := e.Group("/campaign")
	protected.Use(AuthMiddleware(cfg))
	protected.GET("/new", h.CreateCampaignForm)
	protected.POST("", h.CreateCampaign)
	protected.PUT("/:id", h.EditCampaign)
	protected.DELETE("/:id", h.DeleteCampaign)
}

// Custom middleware for protected routes
func AuthMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session := getSession(c, cfg)
			if session == nil || session.Values["user_id"] == nil {
				// Store the original requested URL in the session
				if session != nil {
					session.Values["redirect_after_login"] = c.Request().URL.String()
					_ = session.Save(c.Request(), c.Response().Writer)
				}
				return c.Redirect(http.StatusSeeOther, "/user/login")
			}
			return next(c)
		}
	}
}

func getSession(c echo.Context, cfg *config.Config) *sessions.Session {
	store := c.Get("store").(sessions.Store)
	session, _ := store.Get(c.Request(), cfg.SessionName)
	return session
}

type Handler struct {
	shared.BaseHandler
	service                     ServiceInterface
	representativeLookupService RepresentativeLookupServiceInterface
	emailService                email.Service
	client                      ClientInterface
	mapError                    func(error) (int, string)
}

// HandlerParams for dependency injection
type HandlerParams struct {
	shared.BaseHandlerParams
	fx.In
	Service                     ServiceInterface
	Logger                      loggo.LoggerInterface
	RepresentativeLookupService RepresentativeLookupServiceInterface
	EmailService                email.Service
	Client                      ClientInterface
	ErrorHandler                *shared.ErrorHandler
	TemplateRenderer            shared.TemplateRendererInterface
}

// HandlerResult is the output struct for NewHandler
type HandlerResult struct {
	fx.Out
	Handler *Handler
}

// NewHandler initializes a new Handler
func NewHandler(params HandlerParams) (HandlerResult, error) {
	handler := &Handler{
		BaseHandler:                 shared.NewBaseHandler(params.BaseHandlerParams),
		service:                     params.Service,
		representativeLookupService: params.RepresentativeLookupService,
		emailService:                params.EmailService,
		client:                      params.Client,
		mapError:                    mapErrorToHTTPStatus,
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
	return h.TemplateRenderer.Render(c.Response().Writer, "campaigns", data, c)
}

// CreateCampaignForm handles GET requests for the campaign creation form
func (h *Handler) CreateCampaignForm(c echo.Context) error {
	h.Logger.Debug("Handling CreateCampaignForm request")
	return h.TemplateRenderer.Render(c.Response().Writer, "campaign_create", shared.Data{
		Title:    "Create Campaign",
		PageName: "campaign_create",
		Content:  nil,
	}, c)
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
	return h.TemplateRenderer.Render(c.Response().Writer, "campaign_edit", shared.Data{
		Title:    "Edit Campaign",
		PageName: "campaign_edit",
		Content:  &TemplateData{Campaign: campaign},
	}, c)
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
	return h.TemplateRenderer.Render(c.Response().Writer, "representatives", shared.Data{
		Title:    "Representatives",
		PageName: "representatives",
		Content: &TemplateData{
			Representatives: filteredRepresentatives,
		},
	}, c)
}
