package user

import (
	"fmt"

	"github.com/jonesrussell/loggo"
	"golang.org/x/crypto/bcrypt"
)

type ServiceInterface interface {
	RegisterUser(username, email, password string) error
	VerifyUser(username, password string) (string, error)
}

type Service struct {
	repo   *Repository
	logger *loggo.Logger
}

func NewService(repo *Repository, logger *loggo.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) RegisterUser(username, email, password string) error {
	exists, err := s.repo.UserExists(username, email)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if exists {
		return fmt.Errorf("username or email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	return s.repo.CreateUser(username, email, string(hashedPassword))
}

func (s *Service) VerifyUser(username, password string) (string, error) {
	s.logger.Info(fmt.Sprintf("Verifying user: %s", username))

	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error getting user: %s", username), err)
		return "", fmt.Errorf("error getting user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		s.logger.Warn(fmt.Sprintf("Invalid password for user: %s", username))
		return "", fmt.Errorf("invalid username or password")
	}

	s.logger.Info(fmt.Sprintf("User verified successfully: %s", username))
	return fmt.Sprintf("%d", user.ID), nil
}
