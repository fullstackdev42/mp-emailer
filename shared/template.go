package shared

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer interface
type TemplateRenderer interface {
	Render(w io.Writer, name string, data interface{}, c echo.Context) error
}

// TemplateRendererImpl is a custom renderer for Echo
type TemplateRendererImpl struct {
	templates *template.Template
}

// NewTemplateRenderer creates a new TemplateRendererImpl
func NewTemplateRenderer(templates *template.Template) *TemplateRendererImpl {
	return &TemplateRendererImpl{
		templates: templates,
	}
}

// Render renders a template document
func (t *TemplateRendererImpl) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	viewData := make(map[string]interface{})
	if data != nil {
		viewData["Data"] = data
	}

	// Example: Add request-specific data to viewData
	viewData["RequestID"] = c.Response().Header().Get(echo.HeaderXRequestID)

	var content bytes.Buffer
	if err := t.templates.ExecuteTemplate(&content, name, viewData); err != nil {
		c.Logger().Errorf("failed to execute template %s: %v", name, err)
		return fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	viewData["TemplateContent"] = template.HTML(content.String())
	viewData["PageName"] = name

	return t.templates.ExecuteTemplate(w, "app", viewData)
}

type TemplateManager struct {
	templates *template.Template
}

func (tm *TemplateManager) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	viewData := make(map[string]interface{})
	if data != nil {
		viewData["Data"] = data
	}

	var content bytes.Buffer
	if err := tm.templates.ExecuteTemplate(&content, name, viewData); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	viewData["TemplateContent"] = template.HTML(content.String())
	viewData["PageName"] = name // Add this line to pass the page name to the layout

	return tm.templates.ExecuteTemplate(w, "app", viewData)
}
