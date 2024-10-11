package templates

import (
	"embed"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type TemplateManager struct {
	templates *template.Template
}

func NewTemplateManager(templateFiles embed.FS) (*TemplateManager, error) {
	tmpl, err := template.ParseFS(templateFiles, "web/public/*.html", "web/public/partials/*.html")
	if err != nil {
		return nil, err
	}
	return &TemplateManager{templates: tmpl}, nil
}

func (tm *TemplateManager) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return tm.templates.ExecuteTemplate(w, name, data)
}
