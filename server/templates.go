package server

import (
	"bytes"
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
