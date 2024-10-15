package campaign

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// ExtractUserData extracts user data from the form and returns it as a map.
func (h *Handler) ExtractUserData(c echo.Context) map[string]string {
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

// RenderEmailTemplate renders the email template with provided data.
func (h *Handler) RenderEmailTemplate(c echo.Context, email, content string) error {
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
		return h.HandleError(c, err, http.StatusInternalServerError, "Error rendering email template")
	}
	h.logger.Info("Email template rendered successfully")
	return nil
}

// FetchCampaign fetches a campaign by ID.
func (h *Handler) FetchCampaign(id string) (*Campaign, error) {
	campaign, err := h.service.GetCampaignByID(id)
	if err != nil {
		h.logger.Error("Error fetching campaign", err)
		return nil, err // We'll handle the error in the calling function
	}
	return campaign, nil
}

// HandleError handles errors by logging and rendering an error page.
func (h *Handler) HandleError(c echo.Context, err error, statusCode int, message string) error {
	h.logger.Error(message, err)
	return c.Render(statusCode, "error.html", map[string]interface{}{
		"Error":   message,
		"Details": err.Error(),
	})
}

// ComposeEmail composes an email based on campaign template and user data.
func (h *Handler) ComposeEmail(mp Representative, campaign *Campaign, userData map[string]string) string {
	emailTemplate := campaign.Template
	for key, value := range userData {
		placeholder := fmt.Sprintf("{{%s}}", key)
		emailTemplate = strings.ReplaceAll(emailTemplate, placeholder, value)
	}
	emailTemplate = strings.ReplaceAll(emailTemplate, "{{MP's Name}}", mp.Name)
	emailTemplate = strings.ReplaceAll(emailTemplate, "{{MPEmail}}", mp.Email)
	emailTemplate = strings.ReplaceAll(emailTemplate, "{{Date}}", time.Now().Format("2006-01-02"))
	return emailTemplate
}

// findMP finds the MP based on postal code.
func (h *Handler) findMP(postalCode string) (Representative, error) {
	mpFinder := NewMPFinder(h.client, h.logger)
	mp, err := mpFinder.FindMP(postalCode)
	if err != nil {
		h.logger.Error("Error finding MP", err)
		return Representative{}, err // We'll handle the error in the calling function
	}
	return mp, nil
}
