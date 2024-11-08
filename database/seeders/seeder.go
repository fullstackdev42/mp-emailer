package seeders

import (
	"fmt"

	"github.com/fullstackdev42/mp-emailer/database"
)

type Seeder interface {
	Seed() error
}

type BaseSeeder struct {
	DB database.Interface
}

func RunSeeders(seeders ...Seeder) error {
	for _, seeder := range seeders {
		if err := seeder.Seed(); err != nil {
			return fmt.Errorf("failed to run seeder: %w", err)
		}
	}
	return nil
}