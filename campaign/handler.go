package campaign

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/fullstackdev42/mp-emailer/email"
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
	}
}

func (h *Handler) GetCampaign(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id")) // Convert string to int
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}
	campaign, err := h.service.FetchCampaign(id) // Use the converted int
	if err != nil {
		return h.HandleError(c, err, http.StatusInternalServerError, "Error fetching campaign")
	}
	return c.JSON(http.StatusOK, campaign)
}

func (h *Handler) GetAllCampaigns(c echo.Context) error {
	campaigns, err := h.service.GetAllCampaigns()
	if err != nil {
		return h.HandleError(c, err, http.StatusInternalServerError, "Error fetching campaigns")
	}
	return c.Render(http.StatusOK, "campaigns.html", map[string]interface{}{"Campaigns": campaigns})
}

func (h *Handler) CreateCampaignForm(c echo.Context) error {
	return c.Render(http.StatusOK, "campaign_create.html", nil)
}

func (h *Handler) CreateCampaign(c echo.Context) error {
	name := c.FormValue("name")
	template := c.FormValue("template")
	ownerID, err := user.GetOwnerIDFromSession(c)
	if err != nil {
		return h.HandleError(c, err, http.StatusUnauthorized, "Unauthorized")
	}
	campaign := &Campaign{
		Name:     name,
		Template: template,
		OwnerID:  ownerID,
	}
	if err := h.service.CreateCampaign(campaign); err != nil {
		return h.HandleError(c, err, http.StatusInternalServerError, "Error creating campaign")
	}
	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

func (h *Handler) DeleteCampaign(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID")
	}
	if err := h.service.DeleteCampaign(id); err != nil {
		return h.HandleError(c, err, http.StatusInternalServerError, "Error deleting campaign")
	}
	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

func (h *Handler) EditCampaignForm(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id")) // Convert string to int
	if err != nil {
		return h.HandleError(c, err, http.StatusBadRequest, "Invalid campaign ID")
	}
	campaign, err := h.service.FetchCampaign(id) // Use the converted int
	if err != nil {
		return h.HandleError(c, err, http.StatusInternalServerError, "Error fetching campaign")
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
		return h.HandleError(c, err, http.StatusInternalServerError, "Error updating campaign")
	}
	return c.Redirect(http.StatusSeeOther, "/campaigns/"+strconv.Itoa(id))
}

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
		return h.HandleError(c, err, http.StatusInternalServerError, "Error finding MP")
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.HandleError(c, err, http.StatusBadRequest, "Invalid campaign ID")
	}
	campaign, err := h.service.FetchCampaign(id)
	if err != nil {
		return h.HandleError(c, err, http.StatusInternalServerError, "Error fetching campaign")
	}
	userData := extractUserData(c)

	// Assuming mp is a slice of Representatives, we'll use the first one
	if len(mp) == 0 {
		return h.HandleError(c, errors.New("no representatives found"), http.StatusNotFound, "No representatives found")
	}
	representative := mp[0]

	emailContent := h.service.ComposeEmail(representative, campaign, userData)
	return h.RenderEmailTemplate(c, representative.Email, emailContent)
}

func (h *Handler) HandleError(c echo.Context, err error, statusCode int, message string) error {
	h.logger.Error(message, err)
	return c.Render(statusCode, "error.html", map[string]interface{}{"Error": message, "Details": err.Error()})
}

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
		h.logger.Error("Error rendering email template", err)
		return h.HandleError(c, err, http.StatusInternalServerError, "Error rendering email template")
	}
	h.logger.Info("Email template rendered successfully")
	return nil
}

func (h *Handler) HandleRepresentativeLookup(c echo.Context) error {
	postalCode := c.FormValue("postal_code")
	representativeType := c.FormValue("type")

	representatives, err := h.representativeLookupService.FetchRepresentatives(postalCode)
	if err != nil {
		return h.HandleError(c, err, http.StatusInternalServerError, "Error fetching representatives")
	}

	filters := map[string]string{
		"type": representativeType,
	}

	filteredRepresentatives := h.representativeLookupService.FilterRepresentatives(representatives, filters)

	return c.Render(http.StatusOK, "representatives.html", map[string]interface{}{
		"Representatives": filteredRepresentatives,
	})
}
