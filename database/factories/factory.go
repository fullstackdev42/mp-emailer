package factories

import "github.com/fullstackdev42/mp-emailer/database"

type Factory interface {
	Make() interface{}
	MakeMany(count int) []interface{}
}

type BaseFactory struct {
	DBInterface database.Interface
}
