package campaign

import (
	"github.com/jonesrussell/loggo"
)

// ClientInterface defines the methods a client should implement
type ClientInterface interface {
	FetchRepresentatives(postalCode string) ([]Representative, error)
}

// DefaultClient implements ClientInterface
type DefaultClient struct {
	lookupService RepresentativeLookupServiceInterface
	logger        loggo.LoggerInterface
}

// FetchRepresentatives fetches representatives for the given postal code
func (c *DefaultClient) FetchRepresentatives(postalCode string) ([]Representative, error) {
	c.logger.Info("Fetching representatives for postal code", "postalCode", postalCode)
	return c.lookupService.FetchRepresentatives(postalCode)
}
