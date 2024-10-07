package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
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
	logger, err = loggo.NewLogger("mp-emailer.log", loggo.LevelInfo)
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/submit", handleSubmit)

	logger.Info("Starting server on :8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		logger.Error("Error starting server", err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>MP Emailer</title>
</head>
<body>
    <h1>MP Emailer</h1>
    <form action="/submit" method="post">
        <label for="postalCode">Enter your postal code:</label>
        <input type="text" id="postalCode" name="postalCode" required>
        <input type="submit" value="Find MP">
    </form>
</body>
</html>
`
	t, err := template.New("index").Parse(tmpl)
	if err != nil {
		logger.Error("Error parsing template", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		logger.Error("Error executing template", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postalCode := r.FormValue("postalCode")
	postalCode = strings.ToUpper(strings.ReplaceAll(postalCode, " ", ""))

	mp, err := findMP(postalCode)
	if err != nil {
		logger.Error("Error finding MP", err)
		http.Error(w, "Error finding MP", http.StatusInternalServerError)
		return
	}

	emailContent := composeEmail(mp)

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>MP Email</title>
</head>
<body>
    <h1>Email to MP</h1>
    <p>To: {{.Email}}</p>
    <pre>{{.Content}}</pre>
</body>
</html>
`
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		logger.Error("Error parsing template", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Email   string
		Content string
	}{
		Email:   mp.Email,
		Content: emailContent,
	}

	err = t.Execute(w, data)
	if err != nil {
		logger.Error("Error executing template", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	// Log the email sending (in a real app, you'd actually send the email here)
	logger.Info("Simulated email sending", "to", mp.Email)
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
			logger.Info("MP found", "name", rep.Name, "email", rep.Email)
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
