package campaign

import (
	"github.com/jonesrussell/loggo"
)

type ClientInterface interface {
	FetchRepresentatives(postalCode string) ([]Representative, error)
}

type DefaultClient struct {
	lookupService *RepresentativeLookupService
	logger        loggo.LoggerInterface
}

func NewDefaultClient(logger loggo.LoggerInterface) ClientInterface {
	return &DefaultClient{
		lookupService: NewRepresentativeLookupService(logger),
		logger:        logger,
	}
}

func (c *DefaultClient) FetchRepresentatives(postalCode string) ([]Representative, error) {
	c.logger.Info("Fetching representatives for postal code", "postalCode", postalCode)
	return c.lookupService.FetchRepresentatives(postalCode)
}
