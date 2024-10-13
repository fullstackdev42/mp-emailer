package user

import (
	"database/sql"
	"fmt"

	"github.com/jonesrussell/loggo"
)

type Repository struct {
	db     *sql.DB
	logger *loggo.Logger
}

func NewRepository(db *sql.DB, logger *loggo.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

func (r *Repository) CreateUser(username, email, passwordHash string) error {
	r.logger.Info(fmt.Sprintf("Creating user: %s", username))
	query := "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, username, email, passwordHash)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Error creating user: %s", username), err)
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (r *Repository) GetUserByUsername(username string) (*User, error) {
	var user User

	query := "SELECT id, username, password_hash FROM users WHERE username = ?"
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Error getting user: %s", username), err)
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &user, nil
}

func (r *Repository) UserExists(username, email string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE username = ? OR email = ?"
	var count int
	err := r.db.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %w", err)
	}
	return count > 0, nil
}
