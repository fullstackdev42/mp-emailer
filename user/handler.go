package user

import (
	"net/http"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRoutes registers the user routes
func RegisterRoutes(h *Handler, e *echo.Echo) {
	userGroup := e.Group("/user")
	userGroup.GET("/register", h.RegisterGET)
	userGroup.POST("/register", h.RegisterPOST)
	userGroup.GET("/login", h.LoginGET)
	userGroup.POST("/login", h.LoginPOST)
	userGroup.GET("/logout", h.LogoutGET)
}

// Handler for user routes
type Handler struct {
	service         ServiceInterface
	errorHandler    shared.ErrorHandlerInterface
	flashHandler    *shared.FlashHandler
	Store           sessions.Store
	SessionName     string
	Config          *config.Config
	templateManager shared.TemplateRendererInterface
	repo            RepositoryInterface
}

// RegisterGET handler for the register page
func (h *Handler) RegisterGET(c echo.Context) error {
	data := &shared.Data{
		Title:    "Register",
		PageName: "register",
		Form: shared.FormData{
			Username: "",
			Email:    "",
		},
	}

	return c.Render(http.StatusOK, "register", data)
}

// RegisterPOST handles POST requests to register a new user
func (h *Handler) RegisterPOST(c echo.Context) error {
	params := new(RegisterDTO)
	if err := c.Bind(params); err != nil {
		data := &shared.Data{
			Title:    "Register",
			PageName: "register",
			Content: map[string]interface{}{
				"Username":        params.Username,
				"Email":           params.Email,
				"Password":        params.Password,
				"PasswordConfirm": params.PasswordConfirm,
			},
		}
		return h.templateManager.Render(c.Response().Writer, "register", data, c)
	}

	_, err := h.service.RegisterUser(params)
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Failed to register user", http.StatusInternalServerError)
	}

	if err := h.flashHandler.SetFlashAndSaveSession(c, "Registration successful! Please log in."); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error saving session", http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusSeeOther, "/user/login")
}

// LoginGET handler for the login page
func (h *Handler) LoginGET(c echo.Context) error {
	pageData := shared.Data{
		Title:   "Login",
		Content: nil,
	}
	return h.templateManager.Render(c.Response().Writer, "login", pageData, c)
}

// LoginPOST handler for the login page
func (h *Handler) LoginPOST(c echo.Context) error {
	params := new(LoginDTO)
	if err := c.Bind(params); err != nil {
		h.service.Error("Login binding error", err)
		return h.errorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}

	user, err := h.repo.FindByUsername(params.Username)
	if err != nil || user == nil {
		h.service.Info("Login failed - user not found", "username", params.Username)
		return h.flashHandler.SetFlashAndSaveSession(c, "Invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password)); err != nil {
		h.service.Info("Login failed - invalid password", "username", params.Username)
		return h.flashHandler.SetFlashAndSaveSession(c, "Invalid username or password")
	}

	sess, err := h.Store.Get(c.Request(), h.SessionName)
	if err != nil {
		h.service.Error("Session error", err)
		return h.errorHandler.HandleHTTPError(c, err, "Error getting session", http.StatusInternalServerError)
	}
	sess.Values["user_id"] = user.ID
	sess.Values["username"] = user.Username
	sess.Values["authenticated"] = true

	if err := sess.Save(c.Request(), c.Response().Writer); err != nil {
		h.service.Error("Session save error", err)
		return h.errorHandler.HandleHTTPError(c, err, "Error saving session", http.StatusInternalServerError)
	}

	if err := h.flashHandler.SetFlashAndSaveSession(c, "Successfully logged in!"); err != nil {
		h.service.Error("Flash message error", err)
		return h.errorHandler.HandleHTTPError(c, err, "Error saving session", http.StatusInternalServerError)
	}

	h.service.Info("Login successful", "username", user.Username)
	return c.Redirect(http.StatusSeeOther, "/")
}

// LogoutGET handler for the logout page
func (h *Handler) LogoutGET(c echo.Context) error {
	if err := h.flashHandler.ClearSession(c); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/")
}
