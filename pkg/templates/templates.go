package templates

import (
	"html/template"
	"log"
	"os"
)

var (
	IndexTemplate *template.Template
	EmailTemplate *template.Template
)

func init() {
	var err error

	// Read and parse index.html
	indexHTML, err := os.ReadFile("index.html")
	if err != nil {
		log.Fatalf("Error reading index.html: %v", err)
	}
	IndexTemplate = template.Must(template.New("index").Parse(string(indexHTML)))

	// Read and parse email.html
	emailHTML, err := os.ReadFile("email.html")
	if err != nil {
		log.Fatalf("Error reading email.html: %v", err)
	}
	EmailTemplate = template.Must(template.New("email").Parse(string(emailHTML)))
}
