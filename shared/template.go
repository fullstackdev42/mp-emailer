package shared

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/gorilla/sessions"
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
	store     sessions.Store
	config    *config.Config
}

// Ensure CustomTemplateRenderer implements TemplateRendererInterface
var _ TemplateRendererInterface = (*CustomTemplateRenderer)(nil)

// NewCustomTemplateRenderer creates a new template renderer
func NewCustomTemplateRenderer(t *template.Template, store sessions.Store, cfg *config.Config) TemplateRendererInterface {
	return &CustomTemplateRenderer{
		templates: t,
		store:     store,
		config:    cfg,
	}
}

// Render method implements echo.Renderer and handles rendering templates
func (t *CustomTemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	session, err := t.store.Get(c.Request(), t.config.SessionName)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	flashes := session.Flashes("messages")
	if err := session.Save(c.Request(), c.Response().Writer); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	var pageData *Data
	switch d := data.(type) {
	case Data:
		pageData = &d
	case *Data:
		pageData = d
	default:
		// If data is not of type Data, create new Data struct with content
		pageData = &Data{
			Title:    "MP Emailer",
			PageName: name,
			Content:  data.(map[string]interface{}),
		}
	}

	// Set other required fields
	isAuthenticated, _ := c.Get("IsAuthenticated").(bool)
	if !isAuthenticated {
		isAuthenticated = session.Values["authenticated"] == true
	}
	pageData.IsAuthenticated = isAuthenticated
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
