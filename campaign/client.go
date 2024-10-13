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
	lookupService *RepresentativeLookupService
	logger        *loggo.Logger
}

// NewDefaultClient creates a new DefaultClient
func NewDefaultClient(logger *loggo.Logger) ClientInterface {
	return &DefaultClient{
		lookupService: NewRepresentativeLookupService(logger),
		logger:        logger,
	}
}

// FetchRepresentatives fetches representatives for the given postal code
func (c *DefaultClient) FetchRepresentatives(postalCode string) ([]Representative, error) {
	c.logger.Info("Fetching representatives for postal code", "postalCode", postalCode)
	return c.lookupService.FetchRepresentatives(postalCode)
}
