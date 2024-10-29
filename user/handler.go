package user

import (
	"net/http"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Handler for user routes
type Handler struct {
	repo            RepositoryInterface
	service         ServiceInterface
	Logger          loggo.LoggerInterface
	Store           sessions.Store
	SessionName     string
	Config          *config.Config
	errorHandler    *shared.ErrorHandler
	templateManager shared.CustomTemplateRenderer
}

// LoginUserParams for logging in a user
type LoginUserParams struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

// RegisterGET handler for the register page
func (h *Handler) RegisterGET(c echo.Context) error {
	h.Logger.Debug("RegisterGET handler invoked",
		"method", c.Request().Method,
		"uri", c.Request().RequestURI)

	pageData := shared.PageData{
		Title:   "Register",
		Content: nil,
	}

	err := h.templateManager.Render(c.Response(), "register", pageData, c)
	if err != nil {
		h.Logger.Error("Failed to render register template", err)
		return h.errorHandler.HandleHTTPError(c, err, "Failed to render page", http.StatusInternalServerError)
	}

	return nil
}

// RegisterPOST handler for the register page
func (h *Handler) RegisterPOST(c echo.Context) error {
	// Parse form values
	params := new(RegisterDTO)
	if err := c.Bind(params); err != nil {
		return h.errorHandler.HandleHTTPError(c, err,
			"Invalid input",
			http.StatusBadRequest)
	}

	// Log the attempt
	h.Logger.Info("Attempting to register user", "username", params.Username)

	// Call service
	_, err := h.service.RegisterUser(params)
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err,
			"Failed to register user",
			http.StatusInternalServerError)
	}

	// Set success message in context
	sess, err := h.Store.Get(c.Request(), h.SessionName)
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error getting session", http.StatusInternalServerError)
	}

	h.Logger.Debug("Setting flash message",
		"session_name", h.SessionName,
		"message", "Registration successful! Please log in.")

	sess.AddFlash("Registration successful! Please log in.", "messages")
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		h.Logger.Error("Failed to save session", err)
		return h.errorHandler.HandleHTTPError(c, err, "Error saving session", http.StatusInternalServerError)
	}

	// Redirect on success
	return c.Redirect(http.StatusSeeOther, "/user/login")
}

// LoginGET handler for the login page
func (h *Handler) LoginGET(c echo.Context) error {
	h.Logger.Debug("LoginGET handler invoked", "method", c.Request().Method, "uri", c.Request().RequestURI)
	pageData := shared.PageData{
		Title:   "Login",
		Content: nil,
	}
	return h.templateManager.Render(c.Response(), "login", pageData, c)
}

// LoginPOST handler for the login page
func (h *Handler) LoginPOST(c echo.Context) error {
	params := new(LoginDTO)
	if err := c.Bind(params); err != nil {
		h.Logger.Error("Failed to bind login params", err)
		return h.errorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}

	h.Logger.Info("Login attempt", "username", params.Username)
	user, err := h.repo.GetUserByUsername(params.Username)
	if err != nil {
		h.Logger.Error("Failed to get user by username", err, "username", params.Username)
	}
	if err != nil || user == nil {
		h.Logger.Info("User not found or error occurred", "username", params.Username)
		sess, err := h.Store.Get(c.Request(), h.SessionName)
		if err != nil {
			return h.errorHandler.HandleHTTPError(c, err, "Error getting session", http.StatusInternalServerError)
		}
		sess.AddFlash("Invalid username or password", "messages")
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return h.errorHandler.HandleHTTPError(c, err, "Error saving session", http.StatusInternalServerError)
		}
		return c.Redirect(http.StatusSeeOther, "/user/login")
	}

	h.Logger.Debug("Comparing passwords",
		"username", params.Username,
		"stored_hash_length", len(user.PasswordHash),
		"provided_password_length", len(params.Password))

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password)); err != nil {
		h.Logger.Error("Password comparison failed", err,
			"username", params.Username)
		sess, err := h.Store.Get(c.Request(), h.SessionName)
		if err != nil {
			return h.errorHandler.HandleHTTPError(c, err, "Error getting session", http.StatusInternalServerError)
		}
		sess.AddFlash("Invalid username or password", "messages")
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return h.errorHandler.HandleHTTPError(c, err, "Error saving session", http.StatusInternalServerError)
		}
		return c.Redirect(http.StatusSeeOther, "/user/login")
	}

	h.Logger.Info("Password verified", "username", params.Username)

	// Create a new session
	sess, err := h.Store.Get(c.Request(), h.SessionName)
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error getting session", http.StatusInternalServerError)
	}

	// Set user information in the session
	sess.Values["user_id"] = user.ID
	sess.Values["username"] = user.Username
	sess.Values["authenticated"] = true

	// Add success flash message
	sess.AddFlash("Successfully logged in!", "messages")

	// Save the session
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error saving session", http.StatusInternalServerError)
	}

	h.Logger.Info("Session saved successfully", "username", params.Username)
	return c.Redirect(http.StatusSeeOther, "/")
}

// LogoutGET handler for the logout page
func (h *Handler) LogoutGET(c echo.Context) error {
	if err := h.clearSession(c); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "An error occurred during logout", http.StatusInternalServerError)
	}
	return c.Redirect(http.StatusSeeOther, "/")
}

// GetUser handler for getting a user
func (h *Handler) GetUser(c echo.Context) error {
	params := &GetDTO{Username: c.Param("username")}
	user, err := h.service.GetUser(params)
	if err != nil {
		return h.templateManager.Render(c.Response(), "error", map[string]interface{}{"Message": "User not found", "Username": params.Username}, c)
	}
	return h.templateManager.Render(c.Response(), "user_details", map[string]interface{}{"User": user}, c)
}

func (h *Handler) clearSession(c echo.Context) error {
	sess, err := h.Store.Get(c.Request(), h.SessionName)
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error getting session", http.StatusInternalServerError)
	}
	sess.Values = make(map[interface{}]interface{})
	sess.Options.MaxAge = -1
	return sess.Save(c.Request(), c.Response())
}
