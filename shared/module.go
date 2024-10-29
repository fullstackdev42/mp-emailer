package shared

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/gorilla/sessions"
	"go.uber.org/fx"
)

// Module defines the shared module
//
//nolint:gochecknoglobals
var Module = fx.Options(
	fx.Provide(
		NewErrorHandler,
		NewFlashHandler,
		ProvideTemplates,
	),
)

// ProvideTemplates creates and configures the template renderer
func ProvideTemplates(store sessions.Store) (*CustomTemplateRenderer, error) {
	tmpl := template.New("").Funcs(template.FuncMap{
		"hasPrefix": strings.HasPrefix,
		"safeHTML":  func(s string) template.HTML { return template.HTML(s) },
		"safeURL":   func(s string) template.URL { return template.URL(s) },
	})

	pattern := filepath.Join("web", "templates", "**", "*.gohtml")
	templates, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob templates: %w", err)
	}

	if len(templates) == 0 {
		return nil, fmt.Errorf("no templates found in %s", pattern)
	}

	tmpl, err = tmpl.ParseFiles(templates...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return NewCustomTemplateRenderer(tmpl, store), nil
}
