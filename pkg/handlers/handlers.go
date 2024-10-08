package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/fullstackdev42/mp-emailer/pkg/api"
	"github.com/fullstackdev42/mp-emailer/pkg/database"
	"github.com/fullstackdev42/mp-emailer/pkg/models"
	"github.com/fullstackdev42/mp-emailer/pkg/services"
	"github.com/fullstackdev42/mp-emailer/pkg/templates"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	logger loggo.LoggerInterface
	client api.ClientInterface
	store  *sessions.CookieStore
	db     *database.DB
}

func NewHandler(logger loggo.LoggerInterface, client api.ClientInterface, sessionSecret string, db *database.DB) *Handler {
	store := sessions.NewCookieStore([]byte(sessionSecret))
	return &Handler{logger: logger, client: client, store: store, db: db}
}

func (h *Handler) HandleIndex(c echo.Context) error {
	h.logger.Info("Handling index request")
	return c.Render(http.StatusOK, "index.html", nil)
}

func (h *Handler) HandleSubmit(c echo.Context) error {
	h.logger.Info("Handling submit request")

	postalCode := c.FormValue("postalCode")
	postalCode = strings.ToUpper(strings.ReplaceAll(postalCode, " ", ""))

	// Server-side validation
	postalCodeRegex := regexp.MustCompile(`^[ABCEGHJ-NPRSTVXY]\d[ABCEGHJ-NPRSTV-Z]\d[ABCEGHJ-NPRSTV-Z]\d$`)
	if !postalCodeRegex.MatchString(postalCode) {
		h.logger.Warn("Invalid postal code submitted", "postalCode", postalCode)
		return c.String(http.StatusBadRequest, "Invalid postal code format")
	}

	mpFinder := services.NewMPFinder(h.client, h.logger)

	mp, err := mpFinder.FindMP(postalCode)
	if err != nil {
		h.logger.Error("Error finding MP", err)
		return c.String(http.StatusInternalServerError, "Error finding MP")
	}

	emailContent := composeEmail(mp)

	data := struct {
		Email   string
		Content string
	}{
		Email:   mp.Email,
		Content: emailContent,
	}

	return templates.EmailTemplate.Execute(c.Response().Writer, data)
}

func (h *Handler) HandleEcho(c echo.Context) error {
	type EchoRequest struct {
		Message string `json:"message"`
	}

	req := new(EchoRequest)
	if err := c.Bind(req); err != nil {
		h.logger.Error("Error binding request", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	return c.JSON(http.StatusOK, req)
}

func composeEmail(mp models.Representative) string {
	return fmt.Sprintf("Dear %s,\n\nThis is a sample email content.\n\nBest regards,\nYour constituent", mp.Name)
}

func (h *Handler) HandleLogin(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return templates.LoginTemplate.Execute(c.Response().Writer, nil)
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	// TODO: Implement proper user authentication
	if username == "admin" && password == "password" {
		session, _ := h.store.Get(c.Request(), "session")
		session.Values["user"] = username
		session.Save(c.Request(), c.Response().Writer)
		return c.Redirect(http.StatusSeeOther, "/")
	}

	return c.String(http.StatusUnauthorized, "Invalid credentials")
}

func (h *Handler) HandleLogout(c echo.Context) error {
	session, _ := h.store.Get(c.Request(), "session")
	session.Values["user"] = nil
	session.Save(c.Request(), c.Response().Writer)
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, _ := h.store.Get(c.Request(), "session")
		user := session.Values["user"]
		if user == nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}
		return next(c)
	}
}

func (h *Handler) HandleRegister(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "register.html", nil)
	}

	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	// 1. Validating the input
	if err := validateRegistrationInput(username, email, password, confirmPassword); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// 2. Checking if the username or email already exists
	exists, err := h.db.UserExists(username, email)
	if err != nil {
		h.logger.Error("Error checking user existence", err)
		return c.String(http.StatusInternalServerError, "An error occurred during registration")
	}
	if exists {
		return c.String(http.StatusBadRequest, "Username or email already exists")
	}

	// 3. Hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("Error hashing password", err)
		return c.String(http.StatusInternalServerError, "An error occurred during registration")
	}

	// 4. Storing the new user in the database
	err = h.db.CreateUser(username, email, string(hashedPassword))
	if err != nil {
		h.logger.Error("Error creating user", err)
		return c.String(http.StatusInternalServerError, "An error occurred during registration")
	}

	h.logger.Info("User registered successfully", "username", username, "email", email)

	// Redirect to login page after successful registration
	return c.Redirect(http.StatusSeeOther, "/login")
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
