package user

import (
	"encoding/gob"
	"net/http"

	"github.com/google/uuid"

	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/labstack/echo/v4"
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
	shared.BaseHandler
	Service      ServiceInterface
	FlashHandler *shared.FlashHandler
	Repo         RepositoryInterface
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
		return h.TemplateRenderer.Render(c.Response().Writer, "register", data, c)
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
	return h.TemplateRenderer.Render(c.Response().Writer, "login", pageData, c)
}

// LoginPOST handler for the login page
func (h *Handler) LoginPOST(c echo.Context) error {
	params := new(LoginDTO)
	if err := c.Bind(params); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}

	authenticated, user, err := h.Service.AuthenticateUser(params.Username, params.Password)
	if err != nil || !authenticated {
		return c.Render(http.StatusUnauthorized, "login", &shared.Data{
			Title: "Login",
			Error: "Invalid username or password",
		})
	}

	if err := h.createUserSession(c, user); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Error managing session", http.StatusInternalServerError)
	}

	if err := h.FlashHandler.SetFlashAndSaveSession(c, "Successfully logged in!"); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Error saving session", http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

// createUserSession handles session creation for authenticated users
func (h *Handler) createUserSession(c echo.Context, user *User) error {
	sess, err := h.Store.Get(c.Request(), h.Config.SessionName)
	if err != nil {
		return err
	}

	sess.Values["user_id"] = user.ID
	sess.Values["username"] = user.Username
	sess.Values["authenticated"] = true

	return sess.Save(c.Request(), c.Response().Writer)
}

// LogoutGET handler for the logout page
func (h *Handler) LogoutGET(c echo.Context) error {
	if err := h.FlashHandler.ClearSession(c); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/")
}
