package user

import (
	"time"

	"github.com/google/uuid"
)

// CreateDTO represents the data needed to create a new user
type CreateDTO struct {
	Username string `json:"username" form:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8,max=72"`
}

// LoginDTO represents the data needed for user login
type LoginDTO struct {
	Username string `json:"username" form:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" form:"password" validate:"required,min=8,max=72"`
}

// DTO represents the user data returned to the client
type DTO struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateDTO represents the data for updating a user
type UpdateDTO struct {
	Username string `json:"username" form:"username" validate:"omitempty,min=3,max=50"`
	Email    string `json:"email" form:"email" validate:"omitempty,email"`
	Password string `json:"password" form:"password" validate:"omitempty,min=8,max=72"`
}

// GetDTO represents the data for retrieving a user
type GetDTO struct {
	ID       string `json:"id" form:"id" validate:"omitempty,uuid4"`
	Username string `json:"username" form:"username" validate:"required_without=ID,omitempty,min=3,max=50"`
}

// DeleteDTO represents the data for deleting a user
type DeleteDTO struct {
	ID string `json:"id" form:"id" validate:"required,uuid4"`
}

// RegisterDTO represents the data needed for user registration
type RegisterDTO struct {
	Username        string `json:"username" form:"username" validate:"required,min=3,max=50"`
	Email           string `json:"email" form:"email" validate:"required,email"`
	Password        string `json:"password" form:"password" validate:"required,min=8,max=72"`
	PasswordConfirm string `json:"confirm_password" form:"confirm_password" validate:"required,eqfield=Password"`
}
