package user

import "time"

// CreateDTO represents the data needed to create a new user
type CreateDTO struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

// LoginDTO represents the data needed for user login
type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// DTO represents the user data returned to the client
type DTO struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateDTO represents the data for updating a user
type UpdateDTO struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

// GetDTO represents the data for retrieving a user
type GetDTO struct {
	ID       string `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
}

// DeleteDTO represents the data for deleting a user
type DeleteDTO struct {
	ID string `json:"id" form:"id"`
}

// RegisterDTO represents the data needed for user registration
type RegisterDTO struct {
	Username        string `json:"username" form:"username"`
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	PasswordConfirm string `json:"confirm_password" form:"confirm_password"`
}
