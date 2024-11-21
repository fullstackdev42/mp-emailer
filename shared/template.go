package shared

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/jonesrussell/mp-emailer/config"
	"github.com/jonesrussell/mp-emailer/session"
	"github.com/labstack/echo/v4"
)

type FormData struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}

type Data struct {
	Content         interface{}
	CSRFToken       interface{}
	CurrentPath     string
	Error           string
	Form            FormData
	IsAuthenticated bool
	Messages        []string
	PageName        string
	RequestID       string
	StatusCode      int
	Title           string
}

// TemplateRendererInterface defines the interface for template rendering
type TemplateRendererInterface interface {
	echo.Renderer
}

// CustomTemplateRenderer is a custom renderer for Echo
type CustomTemplateRenderer struct {
	templates *template.Template
	manager   session.Manager
	config    *config.Config
}

// Ensure CustomTemplateRenderer implements TemplateRendererInterface
var _ TemplateRendererInterface = (*CustomTemplateRenderer)(nil)

// NewCustomTemplateRenderer creates a new template renderer
func NewCustomTemplateRenderer(t *template.Template, manager session.Manager, cfg *config.Config) TemplateRendererInterface {
	return &CustomTemplateRenderer{
		templates: t,
		manager:   manager,
		config:    cfg,
	}
}

// Render method implements echo.Renderer and handles rendering templates
func (t *CustomTemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	session, err := t.manager.GetSession(c, t.config.Auth.SessionName)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	flashes := session.Flashes("messages")
	if err := t.manager.SaveSession(c, session); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	var pageData *Data
	switch d := data.(type) {
	case Data:
		pageData = &d
	case *Data:
		pageData = d
	default:
		pageData = &Data{
			Title:    "MP Emailer",
			PageName: name,
			Content:  data.(map[string]interface{}),
		}
	}

	pageData.IsAuthenticated = false

	if authValue, ok := c.Get("IsAuthenticated").(bool); ok && authValue {
		pageData.IsAuthenticated = true
	} else if authValue, ok := session.Values["authenticated"].(bool); ok && authValue {
		pageData.IsAuthenticated = true
	}

	pageData.RequestID = c.Response().Header().Get(echo.HeaderXRequestID)

	messages := make([]string, len(flashes))
	for i, flash := range flashes {
		messages[i] = fmt.Sprint(flash)
	}
	pageData.Messages = append(pageData.Messages, messages...)
	pageData.CSRFToken = c.Get("csrf")
	pageData.CurrentPath = c.Request().URL.Path

	return t.executeTemplate(w, name, pageData)
}

func (t *CustomTemplateRenderer) executeTemplate(w io.Writer, name string, data *Data) error {
	// First render the content into the .Content field
	var contentBuffer bytes.Buffer
	if err := t.templates.ExecuteTemplate(&contentBuffer, name, data); err != nil {
		return fmt.Errorf("failed to execute content template: %w", err)
	}
	data.Content = template.HTML(contentBuffer.String())

	// Then render the full page with the layout
	if err := t.templates.ExecuteTemplate(w, "app", data); err != nil {
		fmt.Printf("Debug - Template Error: %v\nData: %+v\n", err, data)
		return fmt.Errorf("failed to execute template %s: %w", name, err)
	}
	return nil
}
