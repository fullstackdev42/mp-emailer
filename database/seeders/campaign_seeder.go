package seeders

import (
	"fmt"

	"github.com/fullstackdev42/mp-emailer/database"
	"github.com/fullstackdev42/mp-emailer/database/factories"
)

type CampaignSeeder struct {
	BaseSeeder
	UserID string
}

func NewCampaignSeeder(db database.Interface, userID string) *CampaignSeeder {
	return &CampaignSeeder{
		BaseSeeder: BaseSeeder{DB: db},
		UserID:     userID,
	}
}

func (s *CampaignSeeder) Seed() error {
	factory := factories.NewCampaignFactory(s.DB, s.UserID)
	campaigns := factory.MakeMany(3) // Create 3 sample campaigns

	for _, campaign := range campaigns {
		if err := s.DB.Create(campaign); err != nil {
			return fmt.Errorf("failed to seed campaign: %w", err)
		}
	}
	return nil
}