package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jonesrussell/mp-emailer/email"
	"github.com/jonesrussell/mp-emailer/shared"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

// ServiceInterface defines the interface for user services
type ServiceInterface interface {
	shared.LoggableService
	GetUser(params *GetDTO) (*DTO, error)
	RegisterUser(params *RegisterDTO) (*DTO, error)
	LoginUser(params *LoginDTO) (string, error)
	AuthenticateUser(username, password string) (bool, *User, error)
	RequestPasswordReset(dto *PasswordResetDTO) error
	ResetPassword(dto *ResetPasswordDTO) error
}

// Service is the implementation of the UserServiceInterface
type Service struct {
	repo         RepositoryInterface
	validate     *validator.Validate
	emailService email.ServiceInterface
}

// Explicitly implement the ServiceInterface
var _ ServiceInterface = (*Service)(nil)

// ServiceParams for dependency injection
type ServiceParams struct {
	fx.In
	Repo         RepositoryInterface
	Validate     *validator.Validate
	EmailService email.ServiceInterface
}

// NewService creates a new user service
func NewService(params ServiceParams) ServiceInterface {
	return &Service{
		repo:         params.Repo,
		validate:     params.Validate,
		emailService: params.EmailService,
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
	if len(params.Password) < 8 {
		return "", fmt.Errorf("password too short")
	}

	if strings.Count(params.Username, "@") > 1 {
		return "", fmt.Errorf("invalid username format")
	}

	authenticated, _, err := s.AuthenticateUser(params.Username, params.Password)
	if err != nil || !authenticated {
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

// AuthenticateUser validates the username and password combination
func (s *Service) AuthenticateUser(username, password string) (bool, *User, error) {
	if username == "" || password == "" {
		return false, nil, fmt.Errorf("username and password required")
	}

	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return false, nil, fmt.Errorf("authentication failed")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return false, nil, fmt.Errorf("authentication failed")
	}

	return true, user, nil
}

type PasswordResetDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordDTO struct {
	Token           string `json:"token" validate:"required"`
	Password        string `json:"password" validate:"required,min=8,max=72"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
}

// RequestPasswordReset initiates the password reset process
func (s *Service) RequestPasswordReset(dto *PasswordResetDTO) error {
	if err := s.validate.Struct(dto); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	user, err := s.repo.FindByEmail(dto.Email)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Generate reset token
	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	// Save reset token to user
	user.ResetToken = token
	user.ResetTokenExpiresAt = expiresAt
	if err := s.repo.Update(user); err != nil {
		return fmt.Errorf("failed to save reset token: %w", err)
	}

	// Send reset email
	if err := s.emailService.SendPasswordReset(user.Email, token); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

// ResetPassword completes the password reset process
func (s *Service) ResetPassword(dto *ResetPasswordDTO) error {
	if err := s.validate.Struct(dto); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	user, err := s.repo.FindByResetToken(dto.Token)
	if err != nil {
		return fmt.Errorf("invalid or expired token: %w", err)
	}

	if user.ResetTokenExpiresAt.Before(time.Now()) {
		return fmt.Errorf("reset token has expired")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password and clear reset token
	user.PasswordHash = string(hashedPassword)
	user.ResetToken = ""
	user.ResetTokenExpiresAt = time.Time{}

	if err := s.repo.Update(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
