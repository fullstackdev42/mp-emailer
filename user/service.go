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
	repo   RepositoryInterface
	logger loggo.LoggerInterface
}

// Explicitly implement the ServiceInterface
var _ ServiceInterface = (*Service)(nil)

func (s *Service) RegisterUser(params RegisterUserParams) error {
	s.logger.Info(fmt.Sprintf("Registering user: %s", params.Username))

	// Check if user exists
	exists, err := s.repo.UserExists(params.Username, params.Email)
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

	// Create the user with the hashed password
	err = s.repo.CreateUser(params.Username, params.Email, string(hashedPassword))
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error creating user %s", params.Username), err)
		return fmt.Errorf("error creating user: %w", err)
	}

	s.logger.Info(fmt.Sprintf("User registered successfully: %s", params.Username))
	return nil
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
