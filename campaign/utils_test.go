package campaign

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestExtractUserData(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Set form values
	req.Form = map[string][]string{
		"first_name":  {"John"},
		"last_name":   {"Doe"},
		"address_1":   {"123 Main St"},
		"city":        {"Anytown"},
		"province":    {"ON"},
		"postal_code": {"A1A1A1"},
		"email":       {"john.doe@example.com"},
	}

	expected := map[string]string{
		"First Name":    "John",
		"Last Name":     "Doe",
		"Address 1":     "123 Main St",
		"City":          "Anytown",
		"Province":      "ON",
		"Postal Code":   "A1A1A1",
		"Email Address": "john.doe@example.com",
	}

	result := extractUserData(c)
	assert.Equal(t, expected, result)
}

func TestValidatePostalCode(t *testing.T) {
	tests := []struct {
		postalCode string
		expected   string
		expectErr  bool
	}{
		{"A1A 1A1", "A1A1A1", false},
		{"B2B2B2", "B2B2B2", false},
		{"", "", true},
		{"12345", "", true},
		{"Z9Z9Z9", "", true}, // Invalid starting character
	}

	for _, test := range tests {
		result, err := validatePostalCode(test.postalCode)
		if test.expectErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestExtractAndValidatePostalCode(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Valid postal code
	req.Form = url.Values{"postal_code": {"A1A 1A1"}}
	result, err := extractAndValidatePostalCode(c)
	assert.NoError(t, err)
	assert.Equal(t, "A1A1A1", result)
	// Invalid postal code
	req = httptest.NewRequest(echo.POST, "/", strings.NewReader("postal_code=12345"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	c = e.NewContext(req, rec)
	result, err = extractAndValidatePostalCode(c)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}
