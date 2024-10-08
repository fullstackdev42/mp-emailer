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

type TemplateData struct {
	IsAuthenticated bool
	// Add other fields as needed
}

func NewHandler(logger loggo.LoggerInterface, client api.ClientInterface, sessionSecret string, db *database.DB) *Handler {
	store := sessions.NewCookieStore([]byte(sessionSecret))
	return &Handler{logger: logger, client: client, store: store, db: db}
}

func (h *Handler) getSession(c echo.Context) (*sessions.Session, error) {
	return h.store.Get(c.Request(), "session")
}

func (h *Handler) saveSession(session *sessions.Session, c echo.Context) error {
	return session.Save(c.Request(), c.Response().Writer)
}

func (h *Handler) handleError(err error, statusCode int, message string) error {
	h.logger.Error(message, err)
	return echo.NewHTTPError(statusCode, message)
}

func (h *Handler) HandleIndex(c echo.Context) error {
	data := TemplateData{
		IsAuthenticated: c.Get("isAuthenticated").(bool),
	}

	return c.Render(http.StatusOK, "index.html", data)
}

func (h *Handler) HandleSubmit(c echo.Context) error {
	h.logger.Info("Handling submit request")

	postalCode := c.FormValue("postalCode")
	postalCode = strings.ToUpper(strings.ReplaceAll(postalCode, " ", ""))

	// Server-side validation
	postalCodeRegex := regexp.MustCompile(`^[ABCEGHJ-NPRSTVXY]\d[ABCEGHJ-NPRSTV-Z]\d[ABCEGHJ-NPRSTV-Z]\d$`)
	if !postalCodeRegex.MatchString(postalCode) {
		h.logger.Warn("Invalid postal code submitted", "postalCode", postalCode)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid postal code format")
	}

	mpFinder := services.NewMPFinder(h.client, h.logger)

	mp, err := mpFinder.FindMP(postalCode)
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error finding MP")
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
		return h.handleError(err, http.StatusBadRequest, "Error binding request")
	}

	return c.JSON(http.StatusOK, req)
}

func composeEmail(mp models.Representative) string {
	return fmt.Sprintf("Dear %s,\n\nThis is a sample email content.\n\nBest regards,\nYour constituent", mp.Name)
}

func (h *Handler) HandleLogin(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "login.html", nil)
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	valid, err := h.db.VerifyUser(username, password)
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error verifying user")
	}

	if !valid {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	session, err := h.getSession(c)
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Failed to get session")
	}

	session.Values["user"] = username
	if err := h.saveSession(session, c); err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Failed to save session")
	}
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) HandleLogout(c echo.Context) error {
	session, err := h.getSession(c)
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Failed to get session")
	}

	session.Values["user"] = nil
	if err := h.saveSession(session, c); err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Failed to save session during logout")
	}
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, _ := h.store.Get(c.Request(), "session")
		user := session.Values["user"]
		c.Set("isAuthenticated", user != nil)
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 2. Checking if the username or email already exists
	exists, err := h.db.UserExists(username, email)
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error checking user existence")
	}
	if exists {
		return echo.NewHTTPError(http.StatusConflict, "Username or email already exists")
	}

	// 3. Hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error hashing password")
	}

	// 4. Storing the new user in the database
	err = h.db.CreateUser(username, email, string(hashedPassword))
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error creating user")
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
