package campaign

import "github.com/jonesrussell/mp-emailer/logger"

// ClientInterface defines the methods a client should implement
type ClientInterface interface {
	FetchRepresentatives(postalCode string) ([]Representative, error)
}

// DefaultClient implements ClientInterface
type DefaultClient struct {
	lookupService RepresentativeLookupServiceInterface
	Logger        logger.Interface
}

// FetchRepresentatives fetches representatives for the given postal code
func (c *DefaultClient) FetchRepresentatives(postalCode string) ([]Representative, error) {
	c.Logger.Info("Fetching representatives for postal code", "postalCode", postalCode)
	return c.lookupService.FetchRepresentatives(postalCode)
}
