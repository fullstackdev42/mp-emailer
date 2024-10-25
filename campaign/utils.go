package campaign

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

// extractUserData extracts user data from the context
func extractUserData(c echo.Context) map[string]string {
	return map[string]string{
		"First Name":    c.FormValue("first_name"),
		"Last Name":     c.FormValue("last_name"),
		"Address 1":     c.FormValue("address_1"),
		"City":          c.FormValue("city"),
		"Province":      c.FormValue("province"),
		"Postal Code":   c.FormValue("postal_code"),
		"Email Address": c.FormValue("email"),
	}
}

// validatePostalCode validates the postal code
func validatePostalCode(postalCode string) (string, error) {
	if postalCode == "" {
		return "", fmt.Errorf("postal code is required")
	}
	postalCode = strings.ToUpper(strings.ReplaceAll(postalCode, " ", ""))
	postalCodeRegex := regexp.MustCompile(`^[ABCEGHJ-NPRSTVXY]\d[ABCEGHJ-NPRSTV-Z]\d[ABCEGHJ-NPRSTV-Z]\d$`)
	if !postalCodeRegex.MatchString(postalCode) {
		return "", fmt.Errorf("invalid postal code format")
	}
	return postalCode, nil
}

// extractAndValidatePostalCode extracts and validates the postal code
func extractAndValidatePostalCode(c echo.Context) (string, error) {
	postalCode := c.FormValue("postal_code")
	validatedPostalCode, err := validatePostalCode(postalCode)
	if err != nil {
		return "", fmt.Errorf("invalid postal code: %w", err)
	}
	return validatedPostalCode, nil
}
