package user

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(username, email, passwordHash string) error {
	query := "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, username, email, passwordHash)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (r *Repository) GetUserByUsername(username string) (*User, error) {
	query := "SELECT id, username, email, password_hash, created_at FROM users WHERE username = ?"
	var user User
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
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
