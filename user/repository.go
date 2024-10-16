package user

import (
	"database/sql"
	"fmt"

	"github.com/fullstackdev42/mp-emailer/internal/database"
	"github.com/jonesrussell/loggo"
)

// Repository defines the interface for user-related database operations
type Repository interface {
	CreateUser(username, email, passwordHash string) error
	GetUserByUsername(username string) (*User, error)
	UserExists(username, email string) (bool, error)
}

type repositoryImpl struct {
	db     *database.DB
	logger loggo.LoggerInterface
}

// NewRepository creates a new Repository instance
func NewRepository(db *database.DB, logger loggo.LoggerInterface) Repository {
	return &repositoryImpl{db: db, logger: logger}
}

func (r *repositoryImpl) CreateUser(username, email, passwordHash string) error {
	r.logger.Info(fmt.Sprintf("Creating user: %s", username))
	query := "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)"
	_, err := r.db.SQL.Exec(query, username, email, passwordHash)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Error creating user: %s", username), err)
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (r *repositoryImpl) GetUserByUsername(username string) (*User, error) {
	var user User
	query := "SELECT id, username, password_hash FROM users WHERE username = ?"
	err := r.db.SQL.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn(fmt.Sprintf("User not found: %s", username))
			return nil, fmt.Errorf("user not found")
		}
		r.logger.Error(fmt.Sprintf("Error getting user: %s", username), err)
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &user, nil
}

func (r *repositoryImpl) UserExists(username, email string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE username = ? OR email = ?"
	var count int
	err := r.db.SQL.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Error checking user existence: %s, %s", username, email), err)
		return false, fmt.Errorf("error checking user existence: %w", err)
	}
	return count > 0, nil
}
