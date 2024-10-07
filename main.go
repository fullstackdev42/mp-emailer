package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jonesrussell/loggo"
)

var logger loggo.LoggerInterface

type PostalCodeResponse struct {
	Representatives []Representative `json:"representatives_concordance"`
}

type Representative struct {
	Name          string `json:"name"`
	ElectedOffice string `json:"elected_office"`
	Email         string `json:"email"`
}

func main() {
	var err error
	logger, err = loggo.NewLogger("mp-emailer.log", slog.LevelInfo)
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		return
	}

	postalCode := getPostalCode()
	mp, err := findMP(postalCode)
	if err != nil {
		logger.Error("Error finding MP", err)
		return
	}

	emailContent := composeEmail(mp)
	sendEmail(mp.Email, emailContent)
}

func getPostalCode() string {
	var postalCode string
	fmt.Print("Enter your postal code (e.g., A1A1A1): ")
	fmt.Scanln(&postalCode)
	return strings.ToUpper(strings.ReplaceAll(postalCode, " ", ""))
}

func findMP(postalCode string) (Representative, error) {
	url := fmt.Sprintf("https://represent.opennorth.ca/postcodes/%s/?format=json", postalCode)

	resp, err := http.Get(url)
	if err != nil {
		return Representative{}, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Representative{}, fmt.Errorf("error reading response: %w", err)
	}

	var postalCodeResp PostalCodeResponse
	err = json.Unmarshal(body, &postalCodeResp)
	if err != nil {
		return Representative{}, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	for _, rep := range postalCodeResp.Representatives {
		if rep.ElectedOffice == "MP" {
			return rep, nil
		}
	}

	return Representative{}, fmt.Errorf("no MP found for postal code %s", postalCode)
}

func composeEmail(mp Representative) string {
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

func sendEmail(to string, content string) {
	// In a real application, you would implement actual email sending logic here
	// For this example, we'll just print the email content
	logger.Info("Sending email to: %s", to)
	logger.Info("Email content:\n%s", content)
}
