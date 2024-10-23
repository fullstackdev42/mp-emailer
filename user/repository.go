package user

import (
	"database/sql"
	"fmt"
	"time"

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
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username = ?`
	row := r.db.SQL.QueryRow(query, username)

	var user User
	var createdAt, updatedAt sql.NullString
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	user.CreatedAt, _ = parseDateTime(createdAt.String)
	user.UpdatedAt, _ = parseDateTime(updatedAt.String)

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

// Add this helper function if it doesn't exist
func parseDateTime(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02 15:04:05", dateStr)
}
