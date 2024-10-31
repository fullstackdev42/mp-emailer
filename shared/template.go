package shared

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type Data struct {
	IsAuthenticated bool
	Title           string
	Content         interface{}
	Error           string
	Messages        []string
	Campaigns       []CampaignData
	RequestID       string
	PageName        string
	CSRFToken       interface{}
	CurrentPath     string
}

// CampaignData represents the data structure for campaigns
type CampaignData interface {
	GetID() int
	GetName() string
	GetDescription() string
	GetTemplate() string
}

// CustomTemplateRenderer is a custom renderer for Echo
type CustomTemplateRenderer struct {
	templates *template.Template
	store     sessions.Store
}

// NewCustomTemplateRenderer creates a new CustomTemplateRenderer
func NewCustomTemplateRenderer(templates *template.Template, store sessions.Store) *CustomTemplateRenderer {
	return &CustomTemplateRenderer{
		templates: templates,
		store:     store,
	}
}

func (t *CustomTemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	session, err := t.store.Get(c.Request(), "mpe")
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	flashes := session.Flashes("messages")
	if err := session.Save(c.Request(), c.Response().Writer); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	// Initialize pageData with proper authentication state
	pageData, ok := data.(Data)
	if !ok {
		pageData = Data{
			Content: data,
			Title:   name, // default to template name if title not specified
		}
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
	pageData.PageName = name

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

// RenderPage with improved error handling
func (t *CustomTemplateRenderer) RenderPage(c echo.Context, templateName string, pageData Data, errorHandler ErrorHandlerInterface) error {
	// Ensure authentication state is set
	isAuthenticated, _ := c.Get("IsAuthenticated").(bool)
	pageData.IsAuthenticated = isAuthenticated

	if err := t.Render(c.Response(), templateName, pageData, c); err != nil {
		return errorHandler.HandleHTTPError(c, err, "Failed to render page", http.StatusInternalServerError)
	}
	return nil
}
