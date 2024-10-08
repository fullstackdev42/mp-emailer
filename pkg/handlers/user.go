package handlers

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) HandleRegister(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "register.html", nil)
	}

	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	// Validate input
	if err := validateRegistrationInput(username, email, password, confirmPassword); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Register user
	if err := h.registerUser(username, email, password); err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error registering user")
	}

	// Redirect to login page after successful registration
	return c.Redirect(http.StatusSeeOther, "/login")
}

func (h *Handler) registerUser(username, email, password string) error {
	// Check if the username or email already exists
	exists, err := h.db.UserExists(username, email)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if exists {
		return fmt.Errorf("username or email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Store the new user in the database
	if err := h.db.CreateUser(username, email, string(hashedPassword)); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	// Send email to admin
	if err := h.SendAdminNotification(username, email); err != nil {
		h.logger.Error("Failed to send admin notification email", err)
		// Note: We're not returning here as we don't want to interrupt the user flow
		// if the email fails to send
	}

	h.logger.Info("User registered successfully", "username", username, "email", email)
	return nil
}

func (h *Handler) SendAdminNotification(username, email string) error {
	subject := "New User Registration"
	body := fmt.Sprintf("A new user has registered:\nUsername: %s\nEmail: %s", username, email)
	return h.emailService.SendEmail(adminEmail, subject, body)
}

func validateRegistrationInput(username, email, password, confirmPassword string) error {
	if username == "" || email == "" || password == "" || confirmPassword == "" {
		return fmt.Errorf("all fields are required")
	}

	if password != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}
