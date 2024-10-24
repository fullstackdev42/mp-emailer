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

// CustomTemplateRenderer is a custom renderer for Echo
type CustomTemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer creates a new TemplateRendererImpl
func NewTemplateRenderer(templates *template.Template) *CustomTemplateRenderer {
	return &CustomTemplateRenderer{
		templates: templates,
	}
}

// Render renders a template document
func (t *CustomTemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	viewData := map[string]interface{}{
		"Data":      data,
		"RequestID": c.Response().Header().Get(echo.HeaderXRequestID),
		"PageName":  name,
	}

	var content bytes.Buffer
	if err := t.templates.ExecuteTemplate(&content, name, viewData); err != nil {
		c.Logger().Errorf("failed to execute template %s: %v", name, err)
		return fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	viewData["TemplateContent"] = template.HTML(content.String())

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
