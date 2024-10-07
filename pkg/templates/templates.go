package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
)

//go:embed *.html includes/*.html
var templateFiles embed.FS

var (
	IndexTemplate  *template.Template
	PostalTemplate *template.Template
	EmailTemplate  *template.Template
	LoginTemplate  *template.Template
)

func init() {
	var err error

	// Debug: Print all files in the embedded filesystem
	fs.WalkDir(templateFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error walking directory: %v", err)
			return err
		}
		log.Printf("Found file: %s", path)
		return nil
	})

	// Parse the includes first
	includes, err := template.ParseFS(templateFiles, "includes/*.html")
	if err != nil {
		log.Fatalf("Error parsing includes: %v", err)
	}

	// Parse each template, adding the includes to it
	IndexTemplate, err = parseTemplate(includes, "index.html")
	if err != nil {
		log.Fatalf("Error parsing index.html: %v", err)
	}

	PostalTemplate, err = parseTemplate(includes, "postal.html")
	if err != nil {
		log.Fatalf("Error parsing postal.html: %v", err)
	}

	EmailTemplate, err = parseTemplate(includes, "email.html")
	if err != nil {
		log.Fatalf("Error parsing email.html: %v", err)
	}

	LoginTemplate, err = parseTemplate(includes, "login.html")
	if err != nil {
		log.Fatalf("Error parsing login.html: %v", err)
	}
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
