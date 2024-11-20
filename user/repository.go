package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jonesrussell/mp-emailer/database/core"
	"gorm.io/gorm"
)

// RepositoryInterface defines the contract for user repository operations
type RepositoryInterface interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByResetToken(ctx context.Context, token string) (*User, error)
	Update(ctx context.Context, user *User) error
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
func (r *Repository) Create(ctx context.Context, user *User) error {
	if err := r.db.Create(ctx, user); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

// FindByEmail retrieves a user by email
func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.FindOne(ctx, &user, "email = ?", email); err != nil {
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}
	return &user, nil
}

// FindByUsername retrieves a user by username
func (r *Repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	if err := r.db.FindOne(ctx, &user, "username = ?", username); err != nil {
		return nil, fmt.Errorf("error finding user by username: %w", err)
	}
	return &user, nil
}

func (r *Repository) Update(ctx context.Context, user *User) error {
	if err := r.db.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *Repository) FindByResetToken(ctx context.Context, token string) (*User, error) {
	var user User
	if err := r.db.FindOne(ctx, &user, "reset_token = ?", token); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with reset token: %w", err)
		}
		return nil, fmt.Errorf("failed to find user by reset token: %w", err)
	}
	return &user, nil
}
