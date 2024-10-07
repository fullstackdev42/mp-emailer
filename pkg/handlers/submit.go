package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fullstackdev42/mp-emailer/pkg/api"
	"github.com/fullstackdev42/mp-emailer/pkg/models"
	"github.com/fullstackdev42/mp-emailer/pkg/services"
	"github.com/fullstackdev42/mp-emailer/pkg/templates"
	"github.com/jonesrussell/loggo"
)

func HandleSubmit(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(loggo.LoggerInterface)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postalCode := r.FormValue("postalCode")
	postalCode = strings.ToUpper(strings.ReplaceAll(postalCode, " ", ""))

	client := api.NewClient(logger)
	mpFinder := services.NewMPFinder(client, logger)

	mp, err := mpFinder.FindMP(postalCode)
	if err != nil {
		logger.Error("Error finding MP", err)
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
		logger.Error("Error executing template", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	logger.Info("Simulated email sending", "to", mp.Email)
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
