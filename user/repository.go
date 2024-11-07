package user

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RepositoryInterface defines the methods that a user repository must implement
type RepositoryInterface interface {
	UserExists(params *CreateDTO) (bool, error)
	CreateUser(params *CreateDTO) (*User, error)
	GetUserByUsername(username string) (*User, error)
}

// Ensure that Repository implements RepositoryInterface
var _ RepositoryInterface = (*Repository)(nil)

type Repository struct {
	db *gorm.DB
}

func (r *Repository) CreateUser(params *CreateDTO) (*User, error) {
	user := &User{
		ID:           uuid.New().String(),
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: params.Password,
	}

	if err := r.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	return user, nil
}

func (r *Repository) UserExists(params *CreateDTO) (bool, error) {
	var count int64
	err := r.db.Model(&User{}).
		Where("username = ? OR email = ?", params.Username, params.Email).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %w", err)
	}
	return count > 0, nil
}

func (r *Repository) GetUserByUsername(username string) (*User, error) {
	var user User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}
	return &user, nil
}
