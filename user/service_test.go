package user_test

import (
	"errors"
	"testing"

	mocksUser "github.com/fullstackdev42/mp-emailer/mocks/user"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type ServiceTestSuite struct {
	suite.Suite
	mockRepo *mocksUser.MockRepositoryInterface
	service  user.ServiceInterface
}

func (s *ServiceTestSuite) SetupTest() {
	s.mockRepo = mocksUser.NewMockRepositoryInterface(s.T())
	s.service = user.NewService(user.ServiceParams{
		Repo: s.mockRepo,
	})
}

func (s *ServiceTestSuite) TearDownTest() {
	if s.mockRepo != nil {
		s.mockRepo.AssertExpectations(s.T())
	}
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestLoginUser() {
	tests := []struct {
		name    string
		setup   func()
		dto     *user.LoginDTO
		wantErr bool
	}{
		{
			name: "successful login",
			setup: func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
				s.mockRepo.On("FindByUsername", "testuser").Return(&user.User{
					Username:     "testuser",
					PasswordHash: string(hashedPassword),
				}, nil)
			},
			dto: &user.LoginDTO{
				Username: "testuser",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "user not found",
			setup: func() {
				s.mockRepo.On("FindByUsername", "nonexistent").
					Return(nil, errors.New("user not found"))
			},
			dto: &user.LoginDTO{
				Username: "nonexistent",
				Password: "anypassword",
			},
			wantErr: true,
		},
		{
			name: "incorrect password",
			setup: func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
				s.mockRepo.On("FindByUsername", "testuser").Return(&user.User{
					Username:     "testuser",
					PasswordHash: string(hashedPassword),
				}, nil)
			},
			dto: &user.LoginDTO{
				Username: "testuser",
				Password: "wrongpass",
			},
			wantErr: true,
		},
		{
			name:  "empty username",
			setup: func() {},
			dto: &user.LoginDTO{
				Username: "",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name:  "empty password",
			setup: func() {},
			dto: &user.LoginDTO{
				Username: "testuser",
				Password: "",
			},
			wantErr: true,
		},
		{
			name: "invalid username format",
			setup: func() {
				// No mock setup needed - should fail at validation
			},
			dto: &user.LoginDTO{
				Username: "user@with@invalid@chars",
				Password: "validpassword123",
			},
			wantErr: true,
		},
		{
			name: "password too short",
			setup: func() {
				// No mock setup needed - should fail at validation
			},
			dto: &user.LoginDTO{
				Username: "testuser",
				Password: "short",
			},
			wantErr: true,
		},
		{
			name: "repository error",
			setup: func() {
				s.mockRepo.On("FindByUsername", "testuser").
					Return(nil, errors.New("database connection error"))
			},
			dto: &user.LoginDTO{
				Username: "testuser",
				Password: "password",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest() // Reset mock for each test case
			tt.setup()

			token, err := s.service.LoginUser(tt.dto)

			if tt.wantErr {
				assert.Error(s.T(), err)
				assert.Empty(s.T(), token)
			} else {
				assert.NoError(s.T(), err)
				assert.NotEmpty(s.T(), token)
			}
		})
	}
}

func TestPasswordHashing(t *testing.T) {
	password := "mypassword123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	assert.NoError(t, err)
	assert.NotEqual(t, password, string(hashedPassword))

	// Verify correct password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	assert.NoError(t, err)

	// Verify incorrect password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("wrongpassword"))
	assert.Error(t, err)
}
