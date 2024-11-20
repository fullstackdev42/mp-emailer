package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jonesrussell/mp-emailer/database"
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

// RepositoryParams defines the parameters for creating a new Repository
type RepositoryParams struct {
	DB database.Database
}

// Repository implements the RepositoryInterface
type Repository struct {
	db database.Database
}

// NewRepository creates a new instance of Repository
func NewRepository(params RepositoryParams) RepositoryInterface {
	return &Repository{
		db: params.DB,
	}
}

// Create adds a new user to the database
func (r *Repository) Create(ctx context.Context, user *User) error {
	err := r.db.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

// FindByEmail retrieves a user by email
func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.FindOne(ctx, &user, "email = ?", email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with email %s: %w", email, err)
		}
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}
	return &user, nil
}

// FindByUsername retrieves a user by username
func (r *Repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.db.FindOne(ctx, &user, "username = ?", username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with username %s: %w", username, err)
		}
		return nil, fmt.Errorf("error finding user by username: %w", err)
	}
	return &user, nil
}

// Update updates an existing user in the database
func (r *Repository) Update(ctx context.Context, user *User) error {
	err := r.db.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *Repository) FindByResetToken(ctx context.Context, token string) (*User, error) {
	var user User
	err := r.db.FindOne(ctx, &user, "reset_token = ?", token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with reset token: %w", err)
		}
		return nil, fmt.Errorf("failed to find user by reset token: %w", err)
	}
	return &user, nil
}
