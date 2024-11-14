package user

import (
	"encoding/gob"
	"net/http"

	"github.com/google/uuid"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	// Register UUID type with gob encoder for session serialization
	gob.Register(uuid.UUID{})
}

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
	Service         ServiceInterface
	ErrorHandler    shared.ErrorHandlerInterface
	FlashHandler    *shared.FlashHandler
	Store           sessions.Store
	SessionName     string
	Config          *config.Config
	TemplateManager shared.TemplateRendererInterface
	Repo            RepositoryInterface
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
		return h.TemplateManager.Render(c.Response().Writer, "register", data, c)
	}

	_, err := h.Service.RegisterUser(params)
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Failed to register user", http.StatusInternalServerError)
	}

	if err := h.FlashHandler.SetFlashAndSaveSession(c, "Registration successful! Please log in."); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Error saving session", http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusSeeOther, "/user/login")
}

// LoginGET handler for the login page
func (h *Handler) LoginGET(c echo.Context) error {
	pageData := shared.Data{
		Title:   "Login",
		Content: nil,
	}
	return h.TemplateManager.Render(c.Response().Writer, "login", pageData, c)
}

// LoginPOST handler for the login page
func (h *Handler) LoginPOST(c echo.Context) error {
	params := new(LoginDTO)
	h.Service.Info("Starting login attempt", "username", params.Username)

	if err := c.Bind(params); err != nil {
		h.Service.Error("Login binding error", err)
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}

	user, err := h.Repo.FindByUsername(params.Username)
	if err != nil || user == nil {
		h.Service.Info("Login failed - user not found", "username", params.Username)
		return c.Render(http.StatusUnauthorized, "login", &shared.Data{
			Title: "Login",
			Error: "Invalid username or password",
		})
	}

	h.Service.Info("User found, attempting password verification", "username", params.Username)

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password)); err != nil {
		h.Service.Info("Password verification failed", "username", params.Username)
		return c.Render(http.StatusUnauthorized, "login", &shared.Data{
			Title: "Login",
			Error: "Invalid username or password",
		})
	}

	h.Service.Info("Password verified successfully", "username", params.Username)

	sess, err := h.Store.Get(c.Request(), h.SessionName)
	if err != nil {
		h.Service.Error("Session error", err)
		return h.ErrorHandler.HandleHTTPError(c, err, "Error getting session", http.StatusInternalServerError)
	}
	sess.Values["user_id"] = user.ID
	sess.Values["username"] = user.Username
	sess.Values["authenticated"] = true

	if err := sess.Save(c.Request(), c.Response().Writer); err != nil {
		h.Service.Error("Session save error", err)
		return h.ErrorHandler.HandleHTTPError(c, err, "Error saving session", http.StatusInternalServerError)
	}

	if err := h.FlashHandler.SetFlashAndSaveSession(c, "Successfully logged in!"); err != nil {
		h.Service.Error("Flash message error", err)
		return h.ErrorHandler.HandleHTTPError(c, err, "Error saving session", http.StatusInternalServerError)
	}

	h.Service.Info("Login successful", "username", user.Username)
	return c.Redirect(http.StatusSeeOther, "/")
}

// LogoutGET handler for the logout page
func (h *Handler) LogoutGET(c echo.Context) error {
	if err := h.FlashHandler.ClearSession(c); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/")
}
