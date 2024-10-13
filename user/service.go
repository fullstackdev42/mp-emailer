package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
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
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", fmt.Errorf("error getting user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	return fmt.Sprintf("%d", user.ID), nil
}
