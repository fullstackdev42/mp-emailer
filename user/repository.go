package user

import (
	"fmt"

	"github.com/fullstackdev42/mp-emailer/database/core"
)

// RepositoryInterface defines the contract for user repository operations
type RepositoryInterface interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
}

// Repository implements the RepositoryInterface
type Repository struct {
	db core.Interface
}

// NewRepository creates a new instance of Repository
func NewRepository(db core.Interface) RepositoryInterface {
	return &Repository{db: db}
}

// Create adds a new user to the database
func (r *Repository) Create(user *User) error {
	if err := r.db.Create(user); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

// FindByEmail retrieves a user by email
func (r *Repository) FindByEmail(email string) (*User, error) {
	var user User
	if err := r.db.FindOne(&user, "email = ?", email); err != nil {
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}
	return &user, nil
}

// FindByUsername retrieves a user by username
func (r *Repository) FindByUsername(username string) (*User, error) {
	var user User
	if err := r.db.FindOne(&user, "username = ?", username); err != nil {
		return nil, fmt.Errorf("error finding user by username: %w", err)
	}
	return &user, nil
}
