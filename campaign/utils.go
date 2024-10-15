package campaign

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

func ValidatePostalCode(postalCode string) (string, error) {
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

func ExtractAndValidatePostalCode(c echo.Context) (string, error) {
	postalCode := c.FormValue("postal_code")
	validatedPostalCode, err := ValidatePostalCode(postalCode)
	if err != nil {
		return "", fmt.Errorf("invalid postal code: %w", err)
	}
	return validatedPostalCode, nil
}
