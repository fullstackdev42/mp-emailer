package server

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer interface for rendering templates
type TemplateRenderer interface {
	Render(w io.Writer, name string, data interface{}, c echo.Context) error
}

type TemplateManager struct {
	templates *template.Template
}

// NewTemplateManager initializes and parses templates
func NewTemplateManager(templateFiles embed.FS) (*TemplateManager, error) {
	tm := &TemplateManager{}

	// Parse all templates
	tmpl, err := template.New("").ParseFS(templateFiles, "web/templates/**/*.gohtml")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	// Ensure the "app" template exists
	if tmpl.Lookup("app") == nil {
		return nil, fmt.Errorf("layout template 'app' not found")
	}

	tm.templates = tmpl
	return tm, nil
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
