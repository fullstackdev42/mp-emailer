package user

import (
	"fmt"

	"github.com/jonesrussell/loggo"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserParams struct {
	Username string
	Email    string
	Password string
}

type ServiceInterface interface {
	RegisterUser(params RegisterUserParams) error
	VerifyUser(username, password string) (string, error)
}

type Service struct {
	repo   *Repository
	logger loggo.LoggerInterface
}

// Explicitly implement the ServiceInterface
var _ ServiceInterface = (*Service)(nil)

func NewService(repo *Repository, logger loggo.LoggerInterface) (ServiceInterface, error) {
	if repo == nil {
		return nil, fmt.Errorf("repository cannot be nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	return &Service{
		repo:   repo,
		logger: logger,
	}, nil
}

func (s *Service) RegisterUser(params RegisterUserParams) error {
	exists, err := s.repo.UserExists(params.Username, params.Email)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if exists {
		return fmt.Errorf("username or email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	return s.repo.CreateUser(params.Username, params.Email, string(hashedPassword))
}

func (s *Service) VerifyUser(username, password string) (string, error) {
	s.logger.Info(fmt.Sprintf("Verifying user: %s", username))
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		if err.Error() == "user not found" {
			s.logger.Warn(fmt.Sprintf("User not found: %s", username))
			return "", fmt.Errorf("invalid username or password")
		}
		s.logger.Error(fmt.Sprintf("Error getting user: %s", username), err)
		return "", fmt.Errorf("error verifying user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		s.logger.Warn(fmt.Sprintf("Invalid password for user: %s", username))
		return "", fmt.Errorf("invalid username or password")
	}

	s.logger.Info(fmt.Sprintf("User verified successfully: %s", username))
	return user.ID, nil
}
