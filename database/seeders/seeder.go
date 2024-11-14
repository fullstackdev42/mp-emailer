package seeders

import (
	"github.com/fullstackdev42/mp-emailer/database/core"
)

type Seeder interface {
	Seed() error
}

type BaseSeeder struct {
	DB core.Interface
}
