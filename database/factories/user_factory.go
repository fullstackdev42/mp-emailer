package factories

import (
	"github.com/go-faker/faker/v4"
	"github.com/jonesrussell/mp-emailer/database/core"
	"github.com/jonesrussell/mp-emailer/user"
)

type UserFactory struct {
	BaseFactory
}

func NewUserFactory(db core.Interface) *UserFactory {
	return &UserFactory{BaseFactory{DBInterface: db}}
}

func (f *UserFactory) Make() interface{} {
	user := &user.User{
		Username:     faker.Username(),
		Email:        faker.Email(),
		PasswordHash: "$2a$10$7U0oMJZ0qtKcrJPI0otrXOTczXRfHdYD64JZ6oB2QTluNMSF9zmE6", // "password"
	}
	return user
}

func (f *UserFactory) MakeMany(count int) []interface{} {
	var users []interface{}
	for i := 0; i < count; i++ {
		users = append(users, f.Make())
	}
	return users
}
