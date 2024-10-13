package campaign

import (
	"github.com/jonesrussell/loggo"
)

type ClientInterface interface {
	FetchRepresentatives(postalCode string) ([]Representative, error)
}

type DefaultClient struct {
	lookupService *RepresentativeLookupService
	logger        *loggo.Logger
}

func NewDefaultClient(logger *loggo.Logger) ClientInterface {
	return &DefaultClient{
		lookupService: NewRepresentativeLookupService(logger),
		logger:        logger,
	}
}

func (c *DefaultClient) FetchRepresentatives(postalCode string) ([]Representative, error) {
	c.logger.Info("Fetching representatives for postal code", "postalCode", postalCode)
	return c.lookupService.FetchRepresentatives(postalCode)
}
