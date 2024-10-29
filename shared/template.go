package shared

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

// PageData represents the common data structure for page rendering
type PageData struct {
	IsAuthenticated bool
	Title           string
	Content         interface{}
	Error           string
	Messages        []string
	Campaigns       []Campaign // Add this line
}

// Campaign represents a campaign
type Campaign struct {
	ID          int
	Name        string
	Description string
	Template    string
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

// TemplateData represents all data passed to templates
type TemplateData struct {
	Data        PageData
	RequestID   string
	PageName    string
	Messages    []interface{}
	CSRFToken   interface{}
	CurrentPath string
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

	pageData, ok := data.(PageData)
	if !ok {
		pageData = PageData{Content: data, Title: name} // default to template name if title not specified
	}

	templateData := TemplateData{
		Data:        pageData,
		RequestID:   c.Response().Header().Get(echo.HeaderXRequestID),
		PageName:    name,
		Messages:    flashes,
		CSRFToken:   c.Get("csrf"),
		CurrentPath: c.Request().URL.Path,
	}

	return t.executeTemplate(w, name, templateData)
}

func (t *CustomTemplateRenderer) executeTemplate(w io.Writer, name string, data TemplateData) error {
	var content bytes.Buffer
	if err := t.templates.ExecuteTemplate(&content, name, data); err != nil {
		fmt.Printf("Debug - Template Error: %v\nData: %+v\n", err, data)
		return fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	layoutData := data
	layoutData.Data = PageData{Content: template.HTML(content.String()), Title: data.Data.Title}

	return t.templates.ExecuteTemplate(w, "app", layoutData)
}

// RenderPage with improved error handling
func (t *CustomTemplateRenderer) RenderPage(c echo.Context, templateName string, pageData PageData, logger loggo.LoggerInterface, errorHandler *ErrorHandler) error {
	if err := t.Render(c.Response(), templateName, pageData, c); err != nil {
		logger.Error("Failed to render template", err)
		return errorHandler.HandleHTTPError(c, err, "Failed to render page", http.StatusInternalServerError)
	}
	return nil
}
