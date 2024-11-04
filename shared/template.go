package shared

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type Data struct {
	IsAuthenticated bool
	Title           string
	Content         interface{}
	Error           string
	Messages        []string
	PageName        string
	RequestID       string
	CSRFToken       interface{}
	CurrentPath     string
	StatusCode      int
}

// TemplateRendererInterface defines the interface for template rendering
type TemplateRendererInterface interface {
	echo.Renderer
}

// CustomTemplateRenderer is a custom renderer for Echo
type CustomTemplateRenderer struct {
	templates *template.Template
	store     sessions.Store
}

// Ensure CustomTemplateRenderer implements TemplateRendererInterface
var _ TemplateRendererInterface = (*CustomTemplateRenderer)(nil)

// NewCustomTemplateRenderer creates a new template renderer
func NewCustomTemplateRenderer(t *template.Template, store sessions.Store) TemplateRendererInterface {
	return &CustomTemplateRenderer{
		templates: t,
		store:     store,
	}
}

// Render method implements echo.Renderer and handles rendering templates
func (t *CustomTemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	session, err := t.store.Get(c.Request(), "mpe")
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	flashes := session.Flashes("messages")
	if err := session.Save(c.Request(), c.Response().Writer); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	// Create a new pageData with safe defaults
	pageData := Data{
		Title:    "MP Emailer", // Default title
		PageName: name,
		Content:  data, // Store the original data as Content
	}

	// If data is a map, try to extract known fields
	if m, ok := data.(map[string]interface{}); ok {
		// Safely extract title and page name if they exist
		if title, ok := m["Title"].(string); ok {
			pageData.Title = title
		}
		if pageName, ok := m["PageName"].(string); ok {
			pageData.PageName = pageName
		}
		pageData.Content = m // Store the entire map as Content
	}

	// Get authentication state from context and session
	isAuthenticated, _ := c.Get("IsAuthenticated").(bool)
	if !isAuthenticated {
		// Double check session if context value is false
		isAuthenticated = session.Values["authenticated"] == true
	}
	pageData.IsAuthenticated = isAuthenticated

	// Set other required fields
	pageData.RequestID = c.Response().Header().Get(echo.HeaderXRequestID)

	// Convert interface{} flashes to strings
	messages := make([]string, len(flashes))
	for i, flash := range flashes {
		messages[i] = fmt.Sprint(flash)
	}
	pageData.Messages = append(pageData.Messages, messages...)
	pageData.CSRFToken = c.Get("csrf")
	pageData.CurrentPath = c.Request().URL.Path

	return t.executeTemplate(w, name, pageData)
}

func (t *CustomTemplateRenderer) executeTemplate(w io.Writer, name string, data Data) error {
	var content bytes.Buffer
	if err := t.templates.ExecuteTemplate(&content, name, data); err != nil {
		fmt.Printf("Debug - Template Error: %v\nData: %+v\n", err, data)
		return fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	data.Content = template.HTML(content.String())

	return t.templates.ExecuteTemplate(w, "app", data)
}
