package seeders

import (
	"context"
	"fmt"

	"github.com/jonesrussell/mp-emailer/database/factories"
)

type UserSeeder struct {
	BaseSeeder
}

func (s *UserSeeder) Seed() error {
	factory := factories.NewUserFactory(s.DB)
	users := factory.MakeMany(5)

	for _, user := range users {
		if err := s.DB.Create(context.Background(), user); err != nil {
			return fmt.Errorf("failed to seed user: %w", err)
		}
	}
	return nil
}
