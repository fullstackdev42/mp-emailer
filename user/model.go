package user

import (
	"time"

	"github.com/google/uuid"
)

type UserID uuid.UUID

type User struct {
	ID           UserID    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username     string    `json:"username" gorm:"uniqueIndex;not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (u *User) IsValid() bool {
	return u.Username != "" && u.Email != "" && u.PasswordHash != ""
}

// Add any other methods that might be useful for your User struct
