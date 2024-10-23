package user

import (
	"fmt"

	"github.com/fullstackdev42/mp-emailer/internal/database"
	"github.com/google/uuid"
	"github.com/jonesrussell/loggo"
)

// RepositoryInterface defines the methods that a user repository must implement
type RepositoryInterface interface {
	UserExists(username, email string) (bool, error)
	CreateUser(username, email, passwordHash string) error
	GetUserByUsername(username string) (*User, error)
	// Add any other methods that the Repository struct implements
}

// Ensure that Repository implements RepositoryInterface
var _ RepositoryInterface = (*Repository)(nil)

type Repository struct {
	db     *database.DB
	logger loggo.LoggerInterface
}

// NewRepository creates a new Repository instance
func NewRepository(db *database.DB, logger loggo.LoggerInterface) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) CreateUser(username, email, passwordHash string) error {
	user := &User{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}
	return r.db.CreateUser(user.ID, user.Username, user.Email, user.PasswordHash)
}

func (r *Repository) GetUserByUsername(username string) (*User, error) {
	var user User
	err := r.db.SQL.QueryRow("SELECT * FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &user, nil
}

func (r *Repository) UserExists(username, email string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE username = ? OR email = ?"
	var count int
	err := r.db.SQL.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Error checking user existence: %s, %s", username, email), err)
		return false, fmt.Errorf("error checking user existence: %w", err)
	}
	return count > 0, nil
}
