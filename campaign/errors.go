package campaign

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jonesrussell/mp-emailer/middleware"
)

// Standard campaign errors
var (
	ErrCampaignNotFound    = errors.New("campaign not found")
	ErrInvalidCampaignID   = errors.New("invalid campaign ID")
	ErrInvalidCampaignData = errors.New("invalid campaign data")

	ErrUnauthorizedAccess = errors.New("unauthorized access")
	ErrUserNotFound       = errors.New("user not found in session")

	ErrDatabaseOperation = errors.New("database operation failed")
	ErrInvalidPostalCode = errors.New("invalid postal code")
	ErrNoRepresentatives = errors.New("no representatives found")
)

// mapErrorToHTTPStatus maps domain errors to HTTP status codes and messages
func mapErrorToHTTPStatus(err error) (int, string) {
	switch {
	case errors.Is(err, ErrCampaignNotFound):
		return http.StatusNotFound, "Campaign not found"
	case errors.Is(err, ErrInvalidCampaignID):
		return http.StatusBadRequest, "Invalid campaign ID"
	case errors.Is(err, ErrUnauthorizedAccess):
		return http.StatusUnauthorized, "Unauthorized access"
	case errors.Is(err, ErrInvalidCampaignData):
		return http.StatusBadRequest, "Invalid campaign data"
	case errors.Is(err, ErrInvalidPostalCode):
		return http.StatusBadRequest, "Invalid postal code"
	case errors.Is(err, ErrNoRepresentatives):
		return http.StatusNotFound, "No representatives found"
	case errors.Is(err, ErrDatabaseOperation):
		return http.StatusInternalServerError, "Internal server error"
	case errors.Is(err, validator.ValidationErrors{}):
		return http.StatusBadRequest, "Invalid input: Missing required fields"
	case errors.Is(err, middleware.ErrSessionInvalid):
		return http.StatusUnauthorized, "Invalid session"
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}
