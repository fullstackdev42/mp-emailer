package handlers

import (
	"net/http"

	"github.com/fullstackdev42/mp-emailer/pkg/templates"
	"github.com/jonesrussell/loggo"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(loggo.LoggerInterface)

	err := templates.IndexTemplate.Execute(w, nil)
	if err != nil {
		logger.Error("Error executing template", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
