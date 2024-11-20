package seeders

import (
	"context"
	"fmt"

	"github.com/jonesrussell/mp-emailer/database/factories"
)

type CampaignSeeder struct {
	BaseSeeder
	UserID string
}

func (s *CampaignSeeder) Seed() error {
	factory := factories.NewCampaignFactory(s.DB, s.UserID)
	campaigns := factory.MakeMany(3) // Create 3 sample campaigns
	for _, campaign := range campaigns {
		if err := s.DB.Create(context.Background(), campaign); err != nil {
			return fmt.Errorf("failed to seed campaign: %w", err)
		}
	}
	return nil
}
