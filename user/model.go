package user

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ID uuid.UUID

// Scan implements the sql.Scanner interface.
func (id *ID) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal ID value: %v", value)
	}
	parsed, err := uuid.ParseBytes(bytes)
	if err != nil {
		return err
	}
	*id = ID(parsed)
	return nil
}

// Value implements the driver.Valuer interface.
func (id ID) Value() (driver.Value, error) {
	return uuid.UUID(id).MarshalBinary()
}

// User is the model for a user
type User struct {
	ID           string    `json:"id" gorm:"primaryKey;type:char(36)"`
	Username     string    `json:"username" gorm:"uniqueIndex;not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// IsValid checks if the user is valid
func (u User) IsValid() bool {
	return u.Username != "" && u.Email != "" && u.PasswordHash != ""
}
