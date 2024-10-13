package campaign

import (
	"log"
)

type DefaultClient struct{}

type ClientInterface interface {
	FetchRepresentatives(postalCode string) ([]Representative, error)
}

// Implement the methods required by ClientInterface

func (c *DefaultClient) FetchRepresentatives(address string) ([]Representative, error) {
	// Use the address parameter in a log statement to avoid the unused parameter warning
	log.Printf("Fetching representatives for address: %s", address)
	// Implement the logic to fetch representatives here
	// For now, return an empty slice and nil error as a placeholder
	return []Representative{}, nil
}
