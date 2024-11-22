package campaign

import (
	"fmt"

	"github.com/jonesrussell/mp-emailer/logger"
	"go.uber.org/fx"
)

// MPFinderParams is the parameter object for the MPFinder service.
type MPFinderParams struct {
	fx.In

	Client ClientInterface
	Logger logger.Interface
}

// MPFinder is a service that finds Members of Parliament (MPs) based on postal codes.
type MPFinder struct {
	client ClientInterface
	Logger logger.Interface
}

// NewMPFinder is the constructor for the MPFinder service.
func NewMPFinder(params MPFinderParams) *MPFinder {
	return &MPFinder{
		client: params.Client,
		Logger: params.Logger,
	}
}

// FindMP finds the MP for a given postal code.
func (f *MPFinder) FindMP(postalCode string) (Representative, error) {
	if f.client == nil {
		return Representative{}, fmt.Errorf("API client is not initialized")
	}

	representatives, err := f.client.FetchRepresentatives(postalCode)
	if err != nil {
		return Representative{}, fmt.Errorf("error fetching representatives for postal code %s: %w", postalCode, err)
	}

	const mpOffice = "MP"
	for _, rep := range representatives {
		f.Logger.Info("Checking representative", "name", rep.Name, "office", rep.ElectedOffice)
		if rep.ElectedOffice == mpOffice {
			f.Logger.Info("MP found", "name", rep.Name, "email", rep.Email)
			return rep, nil
		}
	}

	f.Logger.Warn("No MP found for postal code", "postalCode", postalCode)
	return Representative{}, fmt.Errorf("no MP found for postal code %s", postalCode)
}
