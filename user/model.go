package user

import (
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User is the model for a user
type User struct {
	shared.BaseModel
	Username     string `gorm:"uniqueIndex;not null" json:"username"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database
func (u *User) BeforeCreate(_ *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
