package user_test

import (
	"fmt"
	"testing"

	mockDB "github.com/jonesrussell/mp-emailer/mocks/core"
	"github.com/jonesrussell/mp-emailer/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRepository(t *testing.T) {
	mockDB := mockDB.NewMockInterface(t)
	repo := user.NewRepository(mockDB)

	t.Run("Create User", func(t *testing.T) {
		user := &user.User{
			Username: "testuser",
			Email:    "test@example.com",
		}

		mockDB.On("Create", user).Return(nil)

		err := repo.Create(user)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("FindByEmail Success", func(t *testing.T) {
		email := "test@example.com"
		expectedUser := &user.User{
			Username: "testuser",
			Email:    email,
		}

		mockDB.On("FindOne", mock.AnythingOfType("*user.User"), "email = ?", email).
			Run(func(args mock.Arguments) {
				arg := args.Get(0).(*user.User)
				*arg = *expectedUser
			}).
			Return(nil)

		foundUser, err := repo.FindByEmail(email)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, foundUser)
		mockDB.AssertExpectations(t)
	})

	t.Run("FindByEmail Not Found", func(t *testing.T) {
		email := "notfound@example.com"
		mockDB.On("FindOne", mock.AnythingOfType("*user.User"), "email = ?", email).
			Return(fmt.Errorf("not found"))

		foundUser, err := repo.FindByEmail(email)
		assert.Error(t, err)
		assert.Nil(t, foundUser)
		mockDB.AssertExpectations(t)
	})
}
