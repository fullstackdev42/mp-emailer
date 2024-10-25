package user

import (
	"errors"
	"net/http"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	repo            RepositoryInterface
	service         ServiceInterface
	Logger          loggo.LoggerInterface
	Store           sessions.Store
	SessionName     string
	Config          *config.Config
	errorHandler    *shared.ErrorHandler
	templateManager shared.TemplateRenderer
}

func (h *Handler) RegisterGET(c echo.Context) error {
	return h.templateManager.Render(c.Response(), "register", nil, c)
}

func (h *Handler) RegisterPOST(c echo.Context) error {
	if h.repo == nil || h.service == nil {
		h.Logger.Error("Repository or Service is not initialized", errors.New("repository or service is not initialized"))
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	// Parse form values
	params := RegisterUserParams{
		Username: c.FormValue("username"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	// Check if user exists
	exists, err := h.repo.UserExists(params.Username, params.Email)
	if err != nil {
		h.Logger.Error("Error checking user existence", err)
		return h.errorHandler.HandleHTTPError(c, err, "Error checking user existence", http.StatusInternalServerError)
	}
	if exists {
		h.Logger.Warn("User already exists")
		return c.String(http.StatusBadRequest, "User already exists")
	}

	// Register the user
	if err := h.service.RegisterUser(params); err != nil {
		h.Logger.Error("Failed to register user", err)
		return h.errorHandler.HandleHTTPError(c, err, "Failed to register user", http.StatusInternalServerError)
	}

	// Redirect on success
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) LoginGET(c echo.Context) error {
	h.Logger.Debug("LoginGET handler invoked", "method", c.Request().Method, "uri", c.Request().RequestURI)
	return h.templateManager.Render(c.Response(), "login", nil, c)
}

func (h *Handler) LoginPOST(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	h.Logger.Info("Login attempt", "username", username)
	user, err := h.repo.GetUserByUsername(username)
	if err != nil || user == nil {
		h.Logger.Warn("Login failed: user not found", "username", username, "error", err)
		return h.templateManager.Render(c.Response(), "login", map[string]interface{}{"Error": "Invalid username or password"}, c)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		h.Logger.Warn("Login failed: incorrect password", "username", username, "error", err)
		return h.templateManager.Render(c.Response(), "login", map[string]interface{}{"Error": "Invalid username or password"}, c)
	}

	h.Logger.Info("Password verified", "username", username)

	// Create a new session
	sess, err := h.Store.Get(c.Request(), h.SessionName)
	if err != nil {
		h.Logger.Error("Error getting session", err)
		return h.templateManager.Render(c.Response(), "error", map[string]interface{}{"Message": "An error occurred while processing your request"}, c)
	}

	// Set user information in the session
	sess.Values["user_id"] = user.ID
	sess.Values["username"] = user.Username
	sess.Values["authenticated"] = true

	// Save the session
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		h.Logger.Error("Error saving session", err)
		return h.templateManager.Render(c.Response(), "error", map[string]interface{}{"Message": "An error occurred while processing your request"}, c)
	}

	h.Logger.Info("Session saved successfully", "username", username)

	// Redirect to the home page or dashboard
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) LogoutGET(c echo.Context) error {
	sess, err := h.Store.Get(c.Request(), h.SessionName)
	if err != nil {
		h.Logger.Error("Error getting session", err)
		return c.Redirect(http.StatusSeeOther, "/")
	}

	// Clear session values and delete cookie
	sess.Values = make(map[interface{}]interface{})
	sess.Options.MaxAge = -1
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		h.Logger.Error("Error saving session", err)
		return c.Redirect(http.StatusSeeOther, "/")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) GetUser(c echo.Context) error {
	username := c.Param("username")
	user, err := h.repo.GetUserByUsername(username)
	if err != nil {
		h.Logger.Warn("User not found", "username", username)
		return h.templateManager.Render(c.Response(), "error", map[string]interface{}{"Message": "User not found", "Username": username}, c)
	}
	return h.templateManager.Render(c.Response(), "user_details", map[string]interface{}{"User": user}, c)
}

func (h *Handler) RequireAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := h.Store.Get(c.Request(), h.SessionName)
		if err != nil || sess.Values["authenticated"] != true {
			return c.Redirect(http.StatusSeeOther, "/user/login")
		}
		return next(c)
	}
}
