package templates

import (
	"embed"
	"html/template"
	"log"
)

//go:embed *.html
var templateFiles embed.FS

var (
	IndexTemplate  *template.Template
	PostalTemplate *template.Template
	EmailTemplate  *template.Template
	LoginTemplate  *template.Template
)

func init() {
	var err error

	IndexTemplate, err = template.ParseFS(templateFiles, "index.html")
	if err != nil {
		log.Fatalf("Error parsing index.html: %v", err)
	}

	PostalTemplate, err = template.ParseFS(templateFiles, "postal.html")
	if err != nil {
		log.Fatalf("Error parsing postal.html: %v", err)
	}

	EmailTemplate, err = template.ParseFS(templateFiles, "email.html")
	if err != nil {
		log.Fatalf("Error parsing email.html: %v", err)
	}

	LoginTemplate, err = template.ParseFS(templateFiles, "login.html")
	if err != nil {
		log.Fatalf("Error parsing login.html: %v", err)
	}
}
