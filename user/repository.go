package user

import (
	"github.com/jonesrussell/loggo"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// RepositoryInterface defines the methods for user repository
type RepositoryInterface interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
}

// Repository implements the user repository interface
type Repository struct {
	db     *gorm.DB
	logger loggo.LoggerInterface
}

// CreateUser creates a new user in the database
func (r *Repository) CreateUser(user *User) error {
	r.logger.Debug("Creating user", "user", user)
	return r.db.Create(user).Error
}

// GetUserByUsername retrieves a user by username
func (r *Repository) GetUserByUsername(username string) (*User, error) {
	var user User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		r.logger.Error("Error fetching user", err, "username", username)
		return nil, err
	}
	return &user, nil
}

// RepositoryParams for dependency injection
type RepositoryParams struct {
	fx.In
	DB     *gorm.DB
	Logger loggo.LoggerInterface
}

// RepositoryResult is the output struct for NewRepository
type RepositoryResult struct {
	fx.Out
	Repository RepositoryInterface `group:"repositories"`
}

// NewRepository creates a new user repository
func NewRepository(params RepositoryParams) RepositoryResult {
	return RepositoryResult{
		Repository: &Repository{
			db:     params.DB,
			logger: params.Logger,
		},
	}
}
