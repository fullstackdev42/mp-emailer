package user_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	mocksEmail "github.com/jonesrussell/mp-emailer/mocks/email"
	mocksUser "github.com/jonesrussell/mp-emailer/mocks/user"
	"github.com/jonesrussell/mp-emailer/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type ServiceTestSuite struct {
	suite.Suite
	mockRepo         *mocksUser.MockRepositoryInterface
	mockEmailService *mocksEmail.MockService
	service          user.ServiceInterface
	validate         *validator.Validate
}

func (s *ServiceTestSuite) SetupTest() {
	s.mockRepo = mocksUser.NewMockRepositoryInterface(s.T())
	s.mockEmailService = mocksEmail.NewMockService(s.T())
	s.validate = validator.New()

	s.service = user.NewService(user.ServiceParams{
		Repo:         s.mockRepo,
		Validate:     s.validate,
		EmailService: s.mockEmailService,
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
				s.mockRepo.On("FindByUsername", mock.Anything, "testuser").Return(&user.User{
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
				s.mockRepo.On("FindByUsername", mock.Anything, "nonexistent").
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
				s.mockRepo.On("FindByUsername", mock.Anything, "testuser").Return(&user.User{
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
				s.mockRepo.On("FindByUsername", mock.Anything, "testuser").
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

			token, err := s.service.LoginUser(context.Background(), tt.dto)

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

func (s *ServiceTestSuite) TestGetUser() {
	tests := []struct {
		name    string
		setup   func()
		dto     *user.GetDTO
		want    *user.DTO
		wantErr bool
	}{
		{
			name: "successful get user",
			setup: func() {
				s.mockRepo.On("FindByUsername", mock.Anything, "testuser").Return(&user.User{
					Username: "testuser",
					Email:    "test@example.com",
				}, nil)
			},
			dto: &user.GetDTO{
				Username: "testuser",
			},
			want: &user.DTO{
				ID:       uuid.New(),
				Username: "testuser",
				Email:    "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "user not found",
			setup: func() {
				s.mockRepo.On("FindByUsername", mock.Anything, "nonexistent").
					Return(nil, errors.New("user not found"))
			},
			dto: &user.GetDTO{
				Username: "nonexistent",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.setup()

			got, err := s.service.GetUser(context.Background(), tt.dto)

			if tt.wantErr {
				assert.Error(s.T(), err)
				assert.Nil(s.T(), got)
			} else {
				assert.NoError(s.T(), err)
				assert.NotNil(s.T(), got)
				assert.Equal(s.T(), tt.want.Username, got.Username)
				assert.Equal(s.T(), tt.want.Email, got.Email)
			}
		})
	}
}

func (s *ServiceTestSuite) TestRegisterUser() {
	tests := []struct {
		name    string
		setup   func()
		dto     *user.RegisterDTO
		want    *user.DTO
		wantErr bool
	}{
		{
			name: "successful registration",
			setup: func() {
				s.mockRepo.On("Create",
					mock.Anything, // First argument for context
					mock.MatchedBy(func(u *user.User) bool { // Second argument for user
						return u.Username == "newuser" &&
							u.Email == "new@example.com" &&
							len(u.PasswordHash) > 0
					}),
				).Run(func(args mock.Arguments) {
					user := args.Get(1).(*user.User)
					user.CreatedAt = time.Now()
					user.UpdatedAt = time.Now()
				}).Return(nil)
			},
			dto: &user.RegisterDTO{
				Username:        "newuser",
				Email:           "new@example.com",
				Password:        "validpassword123",
				PasswordConfirm: "validpassword123",
			},
			want: &user.DTO{
				Username: "newuser",
				Email:    "new@example.com",
			},
			wantErr: false,
		},
		{
			name:  "empty username",
			setup: func() {}, // No mock needed - validation will fail
			dto: &user.RegisterDTO{
				Username:        "",
				Email:           "test@example.com",
				Password:        "validpassword123",
				PasswordConfirm: "validpassword123",
			},
			wantErr: true,
		},
		{
			name:  "invalid email format",
			setup: func() {}, // No mock needed - validation will fail
			dto: &user.RegisterDTO{
				Username:        "testuser",
				Email:           "invalid-email",
				Password:        "validpassword123",
				PasswordConfirm: "validpassword123",
			},
			wantErr: true,
		},
		{
			name:  "password too long",
			setup: func() {}, // No mock needed - validation will fail
			dto: &user.RegisterDTO{
				Username:        "testuser",
				Email:           "test@example.com",
				Password:        strings.Repeat("a", 73),
				PasswordConfirm: strings.Repeat("a", 73),
			},
			wantErr: true,
		},
		{
			name:  "passwords don't match",
			setup: func() {}, // No mock needed - validation will fail
			dto: &user.RegisterDTO{
				Username:        "testuser",
				Email:           "test@example.com",
				Password:        "validpassword123",
				PasswordConfirm: "differentpassword123",
			},
			wantErr: true,
		},
		{
			name: "duplicate username",
			setup: func() {
				s.mockRepo.On("Create", mock.Anything,
					mock.MatchedBy(func(u *user.User) bool {
						return u.Username == "existinguser" &&
							u.Email == "test@example.com" &&
							len(u.PasswordHash) > 0
					}),
				).Return(errors.New("duplicate username"))
			},
			dto: &user.RegisterDTO{
				Username:        "existinguser",
				Email:           "test@example.com",
				Password:        "validpassword123",
				PasswordConfirm: "validpassword123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.setup()

			got, err := s.service.RegisterUser(context.Background(), tt.dto)

			if tt.wantErr {
				s.Error(err)
				s.Nil(got)
			} else {
				s.NoError(err)
				s.NotNil(got)
				s.Equal(tt.want.Username, got.Username)
				s.Equal(tt.want.Email, got.Email)
				// Verify timestamps are set
				s.False(got.CreatedAt.IsZero())
				s.False(got.UpdatedAt.IsZero())
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

// Test the logging interface methods
func (s *ServiceTestSuite) TestLoggingMethods() {
	// These methods should not panic even though they're empty
	s.Run("Info method", func() {
		assert.NotPanics(s.T(), func() {
			s.service.Info("test message")
		})
	})

	s.Run("Warn method", func() {
		assert.NotPanics(s.T(), func() {
			s.service.Warn("test warning")
		})
	})

	s.Run("Error method", func() {
		assert.NotPanics(s.T(), func() {
			s.service.Error("test error", errors.New("test error"))
		})
	})
}

func (s *ServiceTestSuite) TestRegisterUserValidation() {
	tests := []struct {
		name     string
		dto      user.RegisterDTO
		contains string
	}{
		{
			name: "username too short",
			dto: user.RegisterDTO{
				Username: "a",
				Email:    "valid@example.com",
				Password: "validpass123",
			},
			contains: "validation error: Key: 'RegisterDTO.Username' Error:Field validation for 'Username' failed on the 'min' tag",
		},
		{
			name: "invalid email format",
			dto: user.RegisterDTO{
				Username: "validuser",
				Email:    "invalid-email",
				Password: "validpass123",
			},
			contains: "validation error: Key: 'RegisterDTO.Email' Error:Field validation for 'Email' failed on the 'email' tag",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			_, err := s.service.RegisterUser(context.Background(), &tt.dto)
			s.Error(err)
			s.Contains(err.Error(), tt.contains)
		})
	}
}
