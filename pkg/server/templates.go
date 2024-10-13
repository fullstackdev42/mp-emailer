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
	tmpl, err := template.ParseFS(templateFiles, "web/templates/*.html", "web/templates/partials/*.html")
	if err != nil {
		return nil, err
	}
	return &TemplateManager{templates: tmpl}, nil
}

func (tm *TemplateManager) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	var viewData map[string]interface{}

	if data != nil {
		var ok bool
		viewData, ok = data.(map[string]interface{})
		if !ok {
			viewData = make(map[string]interface{})
			viewData["data"] = data
		}
	} else {
		viewData = make(map[string]interface{})
	}

	// Add isAuthenticated to the view data
	viewData["isAuthenticated"] = c.Get("isAuthenticated")

	return tm.templates.ExecuteTemplate(w, name, viewData)
}
