package user

import (
	"fmt"
	"time"

	"github.com/fullstackdev42/mp-emailer/shared"
	"golang.org/x/crypto/bcrypt"
)

type ServiceInterface interface {
	RegisterUser(params *RegisterDTO) (*DTO, error)
	LoginUser(params *LoginDTO) (string, error)
	GetUser(params *GetDTO) (*DTO, error)
}

type Service struct {
	repo   RepositoryInterface
	logger shared.ServiceInterface
	config *Config
}

type Config struct {
	JWTSecret string
	JWTExpiry time.Duration
}

// Explicitly implement the ServiceInterface
var _ ServiceInterface = (*Service)(nil)

// RegisterUserServiceParams for registering a user
type RegisterUserServiceParams struct {
	Username string
	Email    string
	Password string
}

func (s *Service) RegisterUser(params *RegisterDTO) (*DTO, error) {
	s.logger.Info(fmt.Sprintf("Registering user: %s", params.Username))

	// Validate password length
	if len(params.Password) > 72 {
		s.logger.Warn(fmt.Sprintf("Password length exceeds 72 bytes for user: %s", params.Username), nil)
		return nil, fmt.Errorf("password length exceeds 72 bytes")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error hashing password for user %s", params.Username), err)
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// Create the user with the hashed password
	user, err := s.repo.CreateUser(&CreateDTO{
		Username: params.Username,
		Email:    params.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error creating user %s", params.Username), err)
		return nil, err
	}

	s.logger.Info(fmt.Sprintf("User registered successfully: %s", params.Username))
	return &DTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *Service) LoginUser(params *LoginDTO) (string, error) {
	s.logger.Info(fmt.Sprintf("Logging in user: %s", params.Username))
	// Check if user exists and password is correct
	user, err := s.repo.GetUserByUsername(params.Username)
	if err != nil {
		s.logger.Warn(fmt.Sprintf("Failed to find user: %s", params.Username), err)
		return "", fmt.Errorf("invalid username or password")
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password))
	if err != nil {
		s.logger.Warn(fmt.Sprintf("Invalid password for user: %s", params.Username), err)
		return "", fmt.Errorf("invalid username or password")
	}

	// Generate JWT token
	token, err := shared.GenerateToken(params.Username, s.config.JWTSecret, int(s.config.JWTExpiry.Minutes()))
	if err != nil {
		s.logger.Error("Failed to generate token", err)
		return "", fmt.Errorf("error generating token")
	}

	s.logger.Info(fmt.Sprintf("User verified successfully: %s", params.Username))
	return token, nil
}

func (s *Service) GetUser(params *GetDTO) (*DTO, error) {
	user, err := s.repo.GetUserByUsername(params.Username)
	if err != nil {
		return nil, err
	}
	return &DTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
