package user_test

import (
	"testing"

	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserBeforeCreate(t *testing.T) {
	u := &user.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	err := u.BeforeCreate(&gorm.DB{})

	assert.NoError(t, err)
	assert.NotEmpty(t, u.ID)
	assert.NotZero(t, u.ID)
}
