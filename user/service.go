package user

import (
	"fmt"
	"time"

	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/jonesrussell/loggo"
	"golang.org/x/crypto/bcrypt"
)

type ServiceInterface interface {
	RegisterUser(params *CreateDTO) error
	VerifyUser(params *LoginDTO) (string, error)
	GetUser(params *GetDTO) (*DTO, error)
}

type Service struct {
	repo   RepositoryInterface
	logger loggo.LoggerInterface
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

func (s *Service) RegisterUser(params *CreateDTO) error {
	s.logger.Info(fmt.Sprintf("Registering user: %s", params.Username))

	// Check if user exists
	exists, err := s.repo.UserExists(params)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error checking user existence for %s", params.Username), err)
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if exists {
		s.logger.Warn(fmt.Sprintf("User already exists: %s", params.Username))
		return fmt.Errorf("user already exists")
	}

	// Validate password length
	if len(params.Password) > 72 {
		s.logger.Warn(fmt.Sprintf("Password length exceeds 72 bytes for user: %s", params.Username))
		return fmt.Errorf("password length exceeds 72 bytes")
	}

	// Log the password length for debugging (do not log the actual password for security reasons)
	s.logger.Info(fmt.Sprintf("Password length for user %s: %d", params.Username, len(params.Password)))

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error hashing password for user %s", params.Username), err)
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Create a new CreateDTO with the hashed password
	createDTO := &CreateDTO{
		Username: params.Username,
		Email:    params.Email,
		Password: string(hashedPassword),
	}

	// Create the user with the hashed password
	err = s.repo.CreateUser(createDTO)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error creating user %s", params.Username), err)
		return fmt.Errorf("error creating user: %w", err)
	}

	s.logger.Info(fmt.Sprintf("User registered successfully: %s", params.Username))
	return nil
}

func (s *Service) VerifyUser(params *LoginDTO) (string, error) {
	s.logger.Info(fmt.Sprintf("Verifying user: %s", params.Username))
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
