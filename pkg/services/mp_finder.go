package services

import (
	"fmt"

	"github.com/fullstackdev42/mp-emailer/pkg/api"
	"github.com/fullstackdev42/mp-emailer/pkg/models"
	"github.com/jonesrussell/loggo"
)

type MPFinder struct {
	client *api.Client
	logger loggo.LoggerInterface
}

func NewMPFinder(client *api.Client, logger loggo.LoggerInterface) *MPFinder {
	return &MPFinder{client: client, logger: logger}
}

func (f *MPFinder) FindMP(postalCode string) (models.Representative, error) {
	representatives, err := f.client.FetchRepresentatives(postalCode)
	if err != nil {
		return models.Representative{}, err
	}

	for _, rep := range representatives {
		f.logger.Info("Checking representative", "name", rep.Name, "office", rep.ElectedOffice)
		if rep.ElectedOffice == "MP" {
			f.logger.Info("MP found", "name", rep.Name, "email", rep.Email)
			return rep, nil
		}
	}

	f.logger.Warn("No MP found for postal code", "postalCode", postalCode)
	return models.Representative{}, fmt.Errorf("no MP found for postal code %s", postalCode)
}
