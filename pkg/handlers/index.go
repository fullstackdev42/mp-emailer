package handlers

import (
	"net/http"

	"github.com/fullstackdev42/mp-emailer/pkg/templates"
)

func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling index request")

	err := templates.IndexTemplate.Execute(w, nil)
	if err != nil {
		h.logger.Error("Error executing template", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
