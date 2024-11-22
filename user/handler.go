package user

import (
	"encoding/gob"
	"time"

	"github.com/google/uuid"
	"github.com/jonesrussell/mp-emailer/shared"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func init() {
	// Register types with gob encoder for session serialization
	gob.Register(uuid.UUID{})
	gob.Register(time.Time{})
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
	gob.Register(&shared.Data{})
}

// RegisterRoutes registers the user routes
func RegisterRoutes(h *Handler, e *echo.Echo) {
	userGroup := e.Group("/user")
	userGroup.GET("/register", h.RegisterGET)
	userGroup.POST("/register", h.RegisterPOST)
	userGroup.GET("/login", h.LoginGET)
	userGroup.POST("/login", h.LoginPOST)
	userGroup.GET("/logout", h.LogoutGET)
	userGroup.POST("/request-password-reset", h.RequestPasswordResetPOST)
	userGroup.POST("/reset-password", h.ResetPasswordPOST)
}

// Handler for user routes
type Handler struct {
	shared.BaseHandler
	Service ServiceInterface
	Repo    RepositoryInterface
}

type HandlerParams struct {
	fx.In
	shared.BaseHandlerParams
	Service ServiceInterface
	Repo    RepositoryInterface
}

// NewHandler creates a new user handler
func NewHandler(params HandlerParams) *Handler {
	return &Handler{
		BaseHandler: shared.NewBaseHandler(params.BaseHandlerParams),
		Service:     params.Service,
		Repo:        params.Repo,
	}
}

// RegisterGET handler for the register page
func (h *Handler) RegisterGET(c echo.Context) error {
	data := h.prepareRegisterPageData()
	return c.Render(shared.StatusOK, "register", data)
}

func (h *Handler) prepareRegisterPageData() *shared.Data {
	return &shared.Data{
		Title:    "Register",
		PageName: "register",
		Form: shared.FormData{
			Username: "",
			Email:    "",
		},
	}
}

// RegisterPOST handles POST requests to register a new user
func (h *Handler) RegisterPOST(c echo.Context) error {
	params := new(RegisterDTO)
	if err := c.Bind(params); err != nil {
		return h.handleError(c, err, "Invalid input", shared.StatusBadRequest)
	}

	if err := h.registerUser(c, params); err != nil {
		return h.handleError(c, err, "Failed to register user", shared.StatusInternalServerError)
	}

	h.addFlashMessage(c, "Registration successful! Please log in.")
	return c.Redirect(shared.StatusSeeOther, "/user/login")
}

func (h *Handler) registerUser(c echo.Context, params *RegisterDTO) error {
	_, err := h.Service.RegisterUser(c.Request().Context(), params)
	return err
}

func (h *Handler) handleError(c echo.Context, err error, message string, statusCode int) error {
	h.Logger.Error(message, err)
	return h.ErrorHandler.HandleHTTPError(c, err, message, statusCode)
}

func (h *Handler) addFlashMessage(c echo.Context, message string) {
	if err := h.AddFlashMessage(c, message); err != nil {
		h.Logger.Error("Failed to add flash message", err)
	}
}

// LoginGET handler for the login page
func (h *Handler) LoginGET(c echo.Context) error {
	pageData := h.prepareLoginPageData()
	return c.Render(shared.StatusOK, "login", pageData)
}

func (h *Handler) prepareLoginPageData() *shared.Data {
	return &shared.Data{
		Title:   "Login",
		Content: nil,
	}
}

// LoginPOST handler for the login page
func (h *Handler) LoginPOST(c echo.Context) error {
	params := new(LoginDTO)
	if err := c.Bind(params); err != nil {
		return h.handleError(c, err, "Invalid input", shared.StatusBadRequest)
	}

	user, err := h.Service.AuthenticateUser(c.Request().Context(), params.Username, params.Password)
	if err != nil {
		return h.renderLoginError(c, "Invalid username or password")
	}

	if err := h.createUserSession(c, user); err != nil {
		return err
	}

	h.addFlashMessage(c, "Successfully logged in!")
	return c.Redirect(shared.StatusSeeOther, "/")
}

func (h *Handler) renderLoginError(c echo.Context, message string) error {
	return c.Render(shared.StatusUnauthorized, "login", &shared.Data{
		Title: "Login",
		Error: message,
	})
}

func (h *Handler) createUserSession(c echo.Context, user *User) error {
	sess, err := h.GetSession(c)
	if err != nil {
		return h.handleError(c, err, "Error getting session", shared.StatusInternalServerError)
	}

	h.SetSessionValues(sess, user)

	if err := h.SaveSession(c, sess); err != nil {
		return h.handleError(c, err, "Error saving session", shared.StatusInternalServerError)
	}

	return nil
}

// LogoutGET handler for the logout page
func (h *Handler) LogoutGET(c echo.Context) error {
	if err := h.ClearSession(c); err != nil {
		return h.handleError(c, err, "Error clearing session", shared.StatusInternalServerError)
	}

	return c.Redirect(shared.StatusSeeOther, "/")
}

// RequestPasswordResetPOST handles the password reset request
func (h *Handler) RequestPasswordResetPOST(c echo.Context) error {
	ctx := c.Request().Context()
	dto := new(PasswordResetDTO)
	if err := c.Bind(dto); err != nil {
		return h.handleError(c, err, "Invalid request", shared.StatusBadRequest)
	}

	if err := h.Service.RequestPasswordReset(ctx, dto); err != nil {
		return h.handleError(c, err, "Failed to process reset request", shared.StatusInternalServerError)
	}

	h.addFlashMessage(c, "Password reset instructions have been sent to your email")
	return c.Redirect(shared.StatusSeeOther, "/user/login")
}

// ResetPasswordPOST handles the password reset completion
func (h *Handler) ResetPasswordPOST(c echo.Context) error {
	ctx := c.Request().Context()
	dto := new(ResetPasswordDTO)
	if err := c.Bind(dto); err != nil {
		return h.handleError(c, err, "Invalid request", shared.StatusBadRequest)
	}

	if err := h.Service.ResetPassword(ctx, dto); err != nil {
		return h.handleError(c, err, "Failed to reset password", shared.StatusInternalServerError)
	}

	h.addFlashMessage(c, "Your password has been reset successfully")
	return c.Redirect(shared.StatusSeeOther, "/user/login")
}

// RequireAuthentication middleware
func (h *Handler) RequireAuthentication() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !h.IsAuthenticated(c) {
				return c.Redirect(shared.StatusSeeOther, "/user/login")
			}
			return next(c)
		}
	}
}
