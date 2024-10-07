package api

import "github.com/fullstackdev42/mp-emailer/pkg/models"

// ClientInterface defines the methods that the API client must implement.
type ClientInterface interface {
	FetchRepresentatives(postalCode string) ([]models.Representative, error)
}
