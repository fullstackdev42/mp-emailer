package seeders

import (
	"github.com/jonesrussell/mp-emailer/database/core"
)

type Seeder interface {
	Seed() error
}

type BaseSeeder struct {
	DB core.Interface
}
