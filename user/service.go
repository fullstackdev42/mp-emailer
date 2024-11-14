package user

import (
	"fmt"
	"strings"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

// ServiceInterface defines the interface for user services
type ServiceInterface interface {
	shared.LoggableService
	GetUser(params *GetDTO) (*DTO, error)
	RegisterUser(params *RegisterDTO) (*DTO, error)
	LoginUser(params *LoginDTO) (string, error)
}

// Service is the implementation of the UserServiceInterface
type Service struct {
	repo     RepositoryInterface
	validate *validator.Validate
	cfg      *config.Config
}

// Explicitly implement the ServiceInterface
var _ ServiceInterface = (*Service)(nil)

// ServiceParams for dependency injection
type ServiceParams struct {
	fx.In
	Repo     RepositoryInterface
	Validate *validator.Validate
	Cfg      *config.Config
}

// NewService creates a new user service
func NewService(params ServiceParams) ServiceInterface {
	return &Service{
		repo:     params.Repo,
		validate: params.Validate,
		cfg:      params.Cfg,
	}
}

// RegisterUser registers a new user and returns the user DTO
func (s *Service) RegisterUser(params *RegisterDTO) (*DTO, error) {
	// Validate the DTO
	if err := s.validate.Struct(params); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Validate password length
	if len(params.Password) > 72 {
		return nil, fmt.Errorf("password length exceeds 72 bytes")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	user := &User{
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return &DTO{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// LoginUser logs in a user and returns a JWT token
func (s *Service) LoginUser(params *LoginDTO) (string, error) {
	if params.Username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}
	if params.Password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	if strings.Count(params.Username, "@") > 1 {
		return "", fmt.Errorf("invalid username format")
	}

	if len(params.Password) < 8 {
		return "", fmt.Errorf("password too short")
	}

	user, err := s.repo.FindByUsername(params.Username)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password))
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	return params.Username, nil
}

// GetUser retrieves a user by their username
func (s *Service) GetUser(params *GetDTO) (*DTO, error) {
	user, err := s.repo.FindByUsername(params.Username)
	if err != nil {
		return nil, fmt.Errorf("error querying user: %w", err)
	}
	return &DTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// Info logs an info message with the given parameters
func (s *Service) Info(_ string, _ ...interface{}) {
	// Empty implementation as logging is handled by the decorator
}

// Warn logs a warning message with the given parameters
func (s *Service) Warn(_ string, _ ...interface{}) {
	// Empty implementation as logging is handled by the decorator
}

// Error logs an error message with the given parameters
func (s *Service) Error(_ string, _ error, _ ...interface{}) {
	// Empty implementation as logging is handled by the decorator
}
