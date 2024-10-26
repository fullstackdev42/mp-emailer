package user

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/fullstackdev42/mp-emailer/internal/database"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/google/uuid"
	"github.com/jonesrussell/loggo"
	"go.uber.org/fx"
)

// RepositoryInterface defines the methods that a user repository must implement
type RepositoryInterface interface {
	UserExists(params *CreateDTO) (bool, error)
	CreateUser(params *CreateDTO) error
	GetUserByUsername(username string) (*User, error)
	// Add any other methods that the Repository struct implements
}

// Ensure that Repository implements RepositoryInterface
var _ RepositoryInterface = (*Repository)(nil)

type Repository struct {
	db     *database.DB
	logger loggo.LoggerInterface
}

// RepositoryParams defines the parameters for creating a new Repository
type RepositoryParams struct {
	fx.In

	DB     *database.DB
	Logger loggo.LoggerInterface
}

// NewRepository creates a new Repository instance
func NewRepository(params RepositoryParams) RepositoryInterface {
	return &Repository{
		db:     params.DB,
		logger: params.Logger,
	}
}

// CreateUser creates a new user
func (r *Repository) CreateUser(params *CreateDTO) error {
	user := &User{
		ID:           uuid.New().String(),
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: params.Password, // Note: Hash the password in the service layer before passing it to the repository
	}

	err := r.db.CreateUser(user.ID, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Error creating user: %s, %s", params.Username, params.Email), err)
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

// GetUserByUsername gets a user by username
func (r *Repository) GetUserByUsername(username string) (*User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username = ?`
	row := r.db.SQL.QueryRow(query, username)

	var user User
	var createdAt, updatedAt sql.NullString
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		r.logger.Error(fmt.Sprintf("Error getting user by username: %s", username), err)
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	user.CreatedAt, _ = shared.ParseDateTime(createdAt.String)
	user.UpdatedAt, _ = shared.ParseDateTime(updatedAt.String)

	return &user, nil
}

// UserExists checks if a user exists
func (r *Repository) UserExists(params *CreateDTO) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE username = ? OR email = ?"
	var count int
	err := r.db.SQL.QueryRow(query, params.Username, params.Email).Scan(&count)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Error checking user existence: %s, %s", params.Username, params.Email), err)
		return false, fmt.Errorf("error checking user existence: %w", err)
	}
	return count > 0, nil
}
