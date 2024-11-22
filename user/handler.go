package user

import (
	"encoding/gob"
	"time"

	"github.com/google/uuid"

	"github.com/jonesrussell/mp-emailer/shared"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/jonesrussell/mp-emailer/session"
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
	Service        ServiceInterface
	Repo           RepositoryInterface
	SessionManager session.Manager
}

type HandlerParams struct {
	fx.In
	shared.BaseHandlerParams
	Service        ServiceInterface
	Repo           RepositoryInterface
	SessionManager session.Manager
}

// NewHandler creates a new user handler
func NewHandler(params HandlerParams) *Handler {
	return &Handler{
		BaseHandler:    shared.NewBaseHandler(params.BaseHandlerParams),
		Service:        params.Service,
		Repo:           params.Repo,
		SessionManager: params.SessionManager,
	}
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

	return c.Render(shared.StatusOK, "register", data)
}

// RegisterPOST handles POST requests to register a new user
func (h *Handler) RegisterPOST(c echo.Context) error {
	params := new(RegisterDTO)
	if err := c.Bind(params); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid input", shared.StatusBadRequest)
	}

	_, err := h.Service.RegisterUser(c.Request().Context(), params)
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Failed to register user", shared.StatusInternalServerError)
	}

	if err := h.AddFlashMessage(c, "Registration successful! Please log in."); err != nil {
		h.Logger.Error("Failed to add registration flash message", err)
	}

	return c.Redirect(shared.StatusSeeOther, "/user/login")
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
	ctx := c.Request().Context()
	params := new(LoginDTO)
	if err := c.Bind(params); err != nil {
		h.Logger.Error("Failed to bind login parameters", err)
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid input", shared.StatusBadRequest)
	}

	authenticated, user, err := h.Service.AuthenticateUser(ctx, params.Username, params.Password)
	if err != nil {
		h.Logger.Error("Authentication failed", err)
		return h.ErrorHandler.HandleHTTPError(c, err, "Authentication error", shared.StatusInternalServerError)
	}

	if !authenticated || user == nil {
		return c.Render(shared.StatusUnauthorized, "login", &shared.Data{
			Title: "Login",
			Error: "Invalid username or password",
		})
	}

	// Get or create session
	sess, err := h.SessionManager.GetSession(c, h.Config.Auth.SessionName)
	if err != nil {
		h.Logger.Error("Failed to get session", err)
		return h.ErrorHandler.HandleHTTPError(c, err, "Error getting session", shared.StatusInternalServerError)
	}

	// Set session values using the new interface
	h.SessionManager.SetSessionValues(sess, user)

	if err := h.AddFlashMessage(c, "Successfully logged in!"); err != nil {
		h.Logger.Error("Failed to add login flash message", err)
	}

	return c.Redirect(shared.StatusSeeOther, "/")
}

// LogoutGET handler for the logout page
func (h *Handler) LogoutGET(c echo.Context) error {
	if err := h.SessionManager.ClearSession(c, h.Config.Auth.SessionName); err != nil {
		h.Logger.Error("Failed to clear session", err)
		return h.ErrorHandler.HandleHTTPError(c, err, "Error clearing session", shared.StatusInternalServerError)
	}

	if err := h.AddFlashMessage(c, "You have been successfully logged out"); err != nil {
		h.Logger.Error("Failed to add logout flash message", err)
	}

	return c.Redirect(shared.StatusSeeOther, "/")
}

// RequestPasswordResetPOST handles the password reset request
func (h *Handler) RequestPasswordResetPOST(c echo.Context) error {
	ctx := c.Request().Context()
	dto := new(PasswordResetDTO)
	if err := c.Bind(dto); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid request", shared.StatusBadRequest)
	}

	if err := h.Service.RequestPasswordReset(ctx, dto); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Failed to process reset request", shared.StatusInternalServerError)
	}

	if err := h.AddFlashMessage(c, "Password reset instructions have been sent to your email"); err != nil {
		h.Logger.Error("Failed to add password reset request flash message", err)
	}

	return c.Redirect(shared.StatusSeeOther, "/user/login")
}

// ResetPasswordPOST handles the password reset completion
func (h *Handler) ResetPasswordPOST(c echo.Context) error {
	ctx := c.Request().Context()
	dto := new(ResetPasswordDTO)
	if err := c.Bind(dto); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid request", shared.StatusBadRequest)
	}

	if err := h.Service.ResetPassword(ctx, dto); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Failed to reset password", shared.StatusInternalServerError)
	}

	if err := h.AddFlashMessage(c, "Your password has been reset successfully"); err != nil {
		h.Logger.Error("Failed to add password reset success flash message", err)
	}

	return c.Redirect(shared.StatusSeeOther, "/user/login")
}

func (h *Handler) RequireAuthentication() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !h.SessionManager.IsAuthenticated(c) {
				return c.Redirect(shared.StatusSeeOther, "/user/login")
			}
			return next(c)
		}
	}
}
