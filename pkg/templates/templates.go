package templates

import (
	"embed"
	"fmt"
	"html/template"
)

//go:embed *.html includes/*.html
var templateFiles embed.FS

type TemplateManager struct {
	IndexTemplate  *template.Template
	PostalTemplate *template.Template
	EmailTemplate  *template.Template
	LoginTemplate  *template.Template
}

func NewTemplateManager() (*TemplateManager, error) {
	var err error
	manager := &TemplateManager{}

	// Parse the includes first
	includes, err := template.ParseFS(templateFiles, "includes/*.html")
	if err != nil {
		return nil, fmt.Errorf("error parsing includes: %v", err)
	}

	// Parse each template, adding the includes to it
	manager.IndexTemplate, err = parseTemplate(includes, "index.html")
	if err != nil {
		return nil, fmt.Errorf("error parsing index.html: %v", err)
	}

	manager.PostalTemplate, err = parseTemplate(includes, "postal.html")
	if err != nil {
		return nil, fmt.Errorf("error parsing postal.html: %v", err)
	}

	manager.EmailTemplate, err = parseTemplate(includes, "email.html")
	if err != nil {
		return nil, fmt.Errorf("error parsing email.html: %v", err)
	}

	manager.LoginTemplate, err = parseTemplate(includes, "login.html")
	if err != nil {
		return nil, fmt.Errorf("error parsing login.html: %v", err)
	}

	return manager, nil
}

func parseTemplate(includes *template.Template, filename string) (*template.Template, error) {
	content, err := templateFiles.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", filename, err)
	}

	tmpl, err := includes.Clone()
	if err != nil {
		return nil, fmt.Errorf("error cloning includes for %s: %v", filename, err)
	}

	tmpl, err = tmpl.Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %v", filename, err)
	}

	return tmpl, nil
}
