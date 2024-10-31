package user

import (
	"net/http"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Handler for user routes
type Handler struct {
	service         ServiceInterface
	errorHandler    shared.ErrorHandlerInterface
	flashHandler    *shared.FlashHandler
	Store           sessions.Store
	SessionName     string
	Config          *config.Config
	templateManager shared.CustomTemplateRenderer
	repo            RepositoryInterface
}

// RegisterGET handler for the register page
func (h *Handler) RegisterGET(c echo.Context) error {
	pageData := shared.Data{
		Title:   "Register",
		Content: nil,
	}
	return h.templateManager.RenderPage(c, "register", pageData, h.errorHandler)
}

// RegisterPOST handler for the register page
func (h *Handler) RegisterPOST(c echo.Context) error {
	params := new(RegisterDTO)
	if err := c.Bind(params); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
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
	return h.templateManager.RenderPage(c, "login", pageData, h.errorHandler)
}

// LoginPOST handler for the login page
func (h *Handler) LoginPOST(c echo.Context) error {
	params := new(LoginDTO)
	if err := c.Bind(params); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}

	user, err := h.repo.GetUserByUsername(params.Username)
	if err != nil || user == nil {
		return h.flashHandler.SetFlashAndSaveSession(c, "Invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password)); err != nil {
		return h.flashHandler.SetFlashAndSaveSession(c, "Invalid username or password")
	}

	sess, err := h.Store.Get(c.Request(), h.SessionName)
	if err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error getting session", http.StatusInternalServerError)
	}
	sess.Values["user_id"] = user.ID
	sess.Values["username"] = user.Username
	sess.Values["authenticated"] = true

	if err := h.flashHandler.SetFlashAndSaveSession(c, "Successfully logged in!"); err != nil {
		return h.errorHandler.HandleHTTPError(c, err, "Error saving session", http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

// LogoutGET handler for the logout page
func (h *Handler) LogoutGET(c echo.Context) error {
	if err := h.flashHandler.ClearSession(c); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/")
}

// GetUser handler for getting a user
func (h *Handler) GetUser(c echo.Context) error {
	params := &GetDTO{Username: c.Param("username")}
	user, err := h.service.GetUser(params)
	if err != nil {
		return h.templateManager.RenderPage(c, "error", shared.Data{
			Title:   "Error",
			Content: map[string]interface{}{"Message": "User not found", "Username": params.Username},
		}, h.errorHandler)
	}

	return h.templateManager.RenderPage(c, "user_details", shared.Data{
		Title:   "User Details",
		Content: map[string]interface{}{"User": user},
	}, h.errorHandler)
}
