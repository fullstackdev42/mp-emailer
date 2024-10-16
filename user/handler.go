package user

import (
	"net/http"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	repo        Repository
	Logger      loggo.LoggerInterface
	Store       sessions.Store
	SessionName string
}

func NewHandler(
	repo Repository,
	logger loggo.LoggerInterface,
	store sessions.Store,
	config *config.Config,
) *Handler {
	return &Handler{
		repo:        repo,
		Logger:      logger,
		Store:       store,
		SessionName: config.SessionName,
	}
}

func (h *Handler) RegisterGET(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", nil)
}

func (h *Handler) RegisterPOST(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	exists, err := h.repo.UserExists(username, email)
	if err != nil {
		h.Logger.Error("Error checking user existence", err)
		return c.Render(http.StatusInternalServerError, "error.html", map[string]interface{}{
			"Message": "An error occurred while processing your request",
		})
	}
	if exists {
		return c.Render(http.StatusConflict, "register.html", map[string]interface{}{
			"Error": "Username or email already exists",
		})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		h.Logger.Error("Error hashing password", err)
		return c.Render(http.StatusInternalServerError, "error.html", map[string]interface{}{
			"Message": "An error occurred while processing your request",
		})
	}

	if err := h.repo.CreateUser(username, email, string(hashedPassword)); err != nil {
		h.Logger.Error("Error creating user", err)
		return c.Render(http.StatusInternalServerError, "error.html", map[string]interface{}{
			"Message": "An error occurred while creating your account",
		})
	}

	// Create a new session
	sess, err := h.Store.Get(c.Request(), h.SessionName)
	if err != nil {
		h.Logger.Error("Error getting session", err)
		return c.Render(http.StatusInternalServerError, "error.html", map[string]interface{}{
			"Message": "An error occurred while processing your request",
		})
	}
	// Get the newly created user
	user, err := h.repo.GetUserByUsername(username)
	if err != nil {
		h.Logger.Error("Error fetching user after creation", err)
		return c.Render(http.StatusInternalServerError, "error.html", map[string]interface{}{
			"Message": "An error occurred while processing your request",
		})
	}

	// Set user information in the session
	sess.Values["user_id"] = user.ID
	sess.Values["username"] = user.Username
	sess.Values["authenticated"] = true

	// Save the session
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		h.Logger.Error("Error saving session", err)
		return c.Render(http.StatusInternalServerError, "error.html", map[string]interface{}{
			"Message": "An error occurred while processing your request",
		})
	}

	// Redirect to the home page
	return c.Redirect(http.StatusSeeOther, "/")
}

// Handler for GET requests
func (h *Handler) LoginGET(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

// Handler for POST requests
func (h *Handler) LoginPOST(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.repo.GetUserByUsername(username)
	if err != nil {
		h.Logger.Warn("Login failed: user not found", "username", username)
		return c.Render(http.StatusUnauthorized, "login.html", map[string]interface{}{
			"Error": "Invalid username or password",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		h.Logger.Warn("Login failed: incorrect password", "username", username)
		return c.Render(http.StatusUnauthorized, "login.html", map[string]interface{}{
			"Error": "Invalid username or password",
		})
	}

	// Create a new session
	sess, err := h.Store.Get(c.Request(), h.SessionName)
	if err != nil {
		h.Logger.Error("Error getting session", err)
		return c.Render(http.StatusInternalServerError, "error.html", map[string]interface{}{
			"Message": "An error occurred while processing your request",
		})
	}

	// Set user information in the session
	sess.Values["user_id"] = user.ID
	sess.Values["username"] = user.Username
	sess.Values["authenticated"] = true

	// Save the session
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		h.Logger.Error("Error saving session", err)
		return c.Render(http.StatusInternalServerError, "error.html", map[string]interface{}{
			"Message": "An error occurred while processing your request",
		})
	}

	h.Logger.Info("User logged in successfully", "username", username)

	// Redirect to the home page or dashboard
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) LogoutGET(c echo.Context) error {
	sess, err := h.Store.Get(c.Request(), h.SessionName)
	if err != nil {
		h.Logger.Error("Error getting session", err)
		return c.Redirect(http.StatusSeeOther, "/")
	}

	// Clear session values
	sess.Values = make(map[interface{}]interface{})

	// Set MaxAge to -1 to delete the cookie
	sess.Options.MaxAge = -1

	// Save the session (this will delete it)
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		h.Logger.Error("Error saving session", err)
		return c.Redirect(http.StatusSeeOther, "/")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) CreateUser(c echo.Context) error {
	return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
		"Message": "Invalid request payload",
	})
}

func (h *Handler) GetUser(c echo.Context) error {
	username := c.Param("username")
	user, err := h.repo.GetUserByUsername(username)
	if err != nil {
		h.Logger.Warn("User not found", "username", username)
		return c.Render(http.StatusNotFound, "error.html", map[string]interface{}{
			"Message":  "User not found",
			"Username": username,
		})
	}

	return c.Render(http.StatusOK, "user_details.html", map[string]interface{}{
		"User": user,
	})
}
