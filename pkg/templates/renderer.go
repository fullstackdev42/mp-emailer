package templates

import (
	"embed"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	var viewData map[string]interface{}

	if data != nil {
		var ok bool
		viewData, ok = data.(map[string]interface{})
		if !ok {
			viewData = make(map[string]interface{})
		}
	} else {
		viewData = make(map[string]interface{})
	}

	viewData["isAuthenticated"] = c.Get("isAuthenticated")

	return t.templates.ExecuteTemplate(w, name, viewData)
}

func NewRenderer(templateFiles embed.FS) *TemplateRenderer {
	return &TemplateRenderer{
		templates: template.Must(template.ParseFS(templateFiles, "*.html", "includes/*.html")),
	}
}
