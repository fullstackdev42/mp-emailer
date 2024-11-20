package user_test

import (
	"context"
	"fmt"
	"testing"

	mockdb "github.com/jonesrussell/mp-emailer/mocks/database"
	"github.com/jonesrussell/mp-emailer/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	mockDB *mockdb.MockDatabase
	repo   user.RepositoryInterface
}

func (s *RepositoryTestSuite) SetupTest() {
	s.mockDB = mockdb.NewMockDatabase(s.T())
	s.repo = user.NewRepository(user.RepositoryParams{
		DB: s.mockDB,
	})
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) TestCreateUser() {
	testUser := &user.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	s.mockDB.EXPECT().Create(context.Background(), testUser).Return(nil)

	err := s.repo.Create(context.Background(), testUser)
	assert.NoError(s.T(), err)
}

func (s *RepositoryTestSuite) TestFindByEmail() {
	s.Run("Success", func() {
		email := "test@example.com"
		expectedUser := &user.User{
			Username: "testuser",
			Email:    email,
		}

		s.mockDB.EXPECT().FindOne(
			context.Background(),
			&user.User{},
			"email = ?",
			email,
		).Run(func(_ context.Context, dest interface{}, _ string, _ ...interface{}) {
			if userDest, ok := dest.(*user.User); ok {
				*userDest = *expectedUser
			}
		}).Return(nil)

		foundUser, err := s.repo.FindByEmail(context.Background(), email)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedUser.Email, foundUser.Email)
		assert.Equal(s.T(), expectedUser.Username, foundUser.Username)
	})

	s.Run("Not Found", func() {
		email := "notfound@example.com"
		s.mockDB.EXPECT().FindOne(
			context.Background(),
			&user.User{},
			"email = ?",
			email,
		).Return(fmt.Errorf("not found"))

		foundUser, err := s.repo.FindByEmail(context.Background(), email)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), foundUser)
	})
}
