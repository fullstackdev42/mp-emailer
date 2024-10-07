package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/fullstackdev42/mp-emailer/pkg/api"
	"github.com/fullstackdev42/mp-emailer/pkg/models"
	"github.com/fullstackdev42/mp-emailer/pkg/services"
	"github.com/fullstackdev42/mp-emailer/pkg/templates"
)

func (h *Handler) HandleSubmit(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling submit request")

	postalCode := r.FormValue("postalCode")
	postalCode = strings.ToUpper(strings.ReplaceAll(postalCode, " ", ""))

	// Server-side validation
	postalCodeRegex := regexp.MustCompile(`^[ABCEGHJ-NPRSTVXY]\d[ABCEGHJ-NPRSTV-Z]\d[ABCEGHJ-NPRSTV-Z]\d$`)
	if !postalCodeRegex.MatchString(postalCode) {
		h.logger.Warn("Invalid postal code submitted", "postalCode", postalCode)
		http.Error(w, "Invalid postal code format", http.StatusBadRequest)
		return
	}

	client := api.NewClient(h.logger)
	mpFinder := services.NewMPFinder(client, h.logger)

	mp, err := mpFinder.FindMP(postalCode)
	if err != nil {
		h.logger.Error("Error finding MP", err)
		http.Error(w, "Error finding MP", http.StatusInternalServerError)
		return
	}

	emailContent := composeEmail(mp)

	data := struct {
		Email   string
		Content string
	}{
		Email:   mp.Email,
		Content: emailContent,
	}

	err = templates.EmailTemplate.Execute(w, data)
	if err != nil {
		h.logger.Error("Error executing template", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Simulated email sending", "to", mp.Email)
}

func composeEmail(mp models.Representative) string {
	template := `
Dear %s,

I am writing to express my concerns about [ISSUE].

[BODY OF THE EMAIL]

Thank you for your time and consideration.

Sincerely,
[YOUR NAME]
`
	return fmt.Sprintf(template, mp.Name)
}
