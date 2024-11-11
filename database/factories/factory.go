package factories

import "github.com/fullstackdev42/mp-emailer/database/core"

type Factory interface {
	Make() interface{}
	MakeMany(count int) []interface{}
}

type BaseFactory struct {
	DBInterface core.Interface
}
