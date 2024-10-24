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
	return h.templateManager.Render(c.Response(), "register.gohtml", nil, c)
}

func (h *Handler) RegisterPOST(c echo.Context) error {
	// Parse form values
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Create params object
	params := RegisterUserParams{Username: username, Email: email, Password: password}

	// Check if user exists
	exists, err := h.repo.UserExists(username, email)
	if err != nil {
		h.Logger.Error("Error checking user existence", err)
		return h.errorHandler.HandleHTTPError(c, err, "Error checking user existence", http.StatusInternalServerError)
	}
	if exists {
		h.Logger.Warn("User already exists")
		return c.String(http.StatusBadRequest, "User already exists")
	}

	// Register the user
	err = h.service.RegisterUser(params)
	if err != nil {
		h.Logger.Error("Failed to register user", err)
		return h.errorHandler.HandleHTTPError(c, err, "Failed to register user", http.StatusInternalServerError)
	}

	// Create user in repo
	err = h.repo.CreateUser(username, email, password)
	if err != nil {
		h.Logger.Error("Error creating user", err)
		return h.errorHandler.HandleHTTPError(c, err, "Error creating user", http.StatusInternalServerError)
	}

	// Get user by username
	_, err = h.repo.GetUserByUsername(username)
	if err != nil {
		h.Logger.Error("Error getting user by username", err)
		return h.errorHandler.HandleHTTPError(c, err, "Error getting user by username", http.StatusInternalServerError)
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
	if err != nil {
		h.Logger.Warn("Login failed: user not found", "username", username, "error", err)
		return c.Render(http.StatusUnauthorized, "login.gohtml", map[string]interface{}{"Error": "Invalid username or password"})
	}

	h.Logger.Info("User found", "username", username, "user_id", user.ID)
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		h.Logger.Warn("Login failed: incorrect password", "username", username, "error", err)
		return c.Render(http.StatusUnauthorized, "login.gohtml", map[string]interface{}{"Error": "Invalid username or password"})
	}

	h.Logger.Info("Password verified", "username", username)

	// Create a new session
	sess, err := h.Store.Get(c.Request(), h.SessionName)
	if err != nil {
		h.Logger.Error("Error getting session", err)
		return c.Render(http.StatusInternalServerError, "error.gohtml", map[string]interface{}{"Message": "An error occurred while processing your request"})
	}

	h.Logger.Info("Session created", "session_name", h.SessionName)

	// Set user information in the session
	sess.Values["user_id"] = user.ID
	sess.Values["username"] = user.Username
	sess.Values["authenticated"] = true

	// Save the session
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		h.Logger.Error("Error saving session", err)
		return c.Render(http.StatusInternalServerError, "error.gohtml", map[string]interface{}{"Message": "An error occurred while processing your request"})
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
	return c.Render(http.StatusBadRequest, "error.gohtml", map[string]interface{}{"Message": "Invalid request payload"})
}

func (h *Handler) GetUser(c echo.Context) error {
	username := c.Param("username")
	user, err := h.repo.GetUserByUsername(username)
	if err != nil {
		h.Logger.Warn("User not found", "username", username)
		return c.Render(http.StatusNotFound, "error.gohtml", map[string]interface{}{"Message": "User not found", "Username": username})
	}
	return c.Render(http.StatusOK, "user_details.gohtml", map[string]interface{}{"User": user})
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
