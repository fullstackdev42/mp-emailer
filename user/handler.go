package user

import (
	"encoding/gob"
	"net/http"

	"github.com/google/uuid"

	"github.com/jonesrussell/mp-emailer/shared"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/jonesrussell/mp-emailer/session"
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

	return c.Render(http.StatusOK, "register", data)
}

// RegisterPOST handles POST requests to register a new user
func (h *Handler) RegisterPOST(c echo.Context) error {
	params := new(RegisterDTO)
	if err := c.Bind(params); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}

	_, err := h.Service.RegisterUser(c.Request().Context(), params)
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Failed to register user", http.StatusInternalServerError)
	}

	sess, err := h.SessionManager.GetSession(c, h.Config.Auth.SessionName)
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Error getting session", http.StatusInternalServerError)
	}

	sess.AddFlash("Registration successful! Please log in.")
	if err := h.SessionManager.SaveSession(c, sess); err != nil {
		h.Logger.Error("Failed to save session", err)
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
	h.Logger.Debug("Processing login request")

	params := new(LoginDTO)
	if err := c.Bind(params); err != nil {
		h.Logger.Error("Failed to bind login parameters", err)
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid input", http.StatusBadRequest)
	}

	h.Logger.Debug("Attempting user authentication", "username", params.Username)

	authenticated, user, err := h.Service.AuthenticateUser(c.Request().Context(), params.Username, params.Password)
	if err != nil {
		h.Logger.Error("Authentication failed", err)
		return h.ErrorHandler.HandleHTTPError(c, err, "Authentication error", http.StatusInternalServerError)
	}

	if !authenticated || user == nil {
		h.Logger.Debug("Invalid login attempt", "username", params.Username)
		return c.Render(http.StatusUnauthorized, "login", &shared.Data{
			Title: "Login",
			Error: "Invalid username or password",
		})
	}

	h.Logger.Debug("User authenticated successfully", "username", params.Username, "userID", user.ID)

	// Get or create session
	sess, err := h.SessionManager.GetSession(c, h.Config.Auth.SessionName)
	if err != nil {
		h.Logger.Error("Failed to get session", err)
		return h.ErrorHandler.HandleHTTPError(c, err, "Error getting session", http.StatusInternalServerError)
	}

	// Set session values using the new interface
	h.SessionManager.SetSessionValues(sess, user)

	// Add flash message to session
	sess.AddFlash("Successfully logged in!")

	// Save session
	if err := h.SessionManager.SaveSession(c, sess); err != nil {
		h.Logger.Error("Failed to save session", err)
		return h.ErrorHandler.HandleHTTPError(c, err, "Error saving session", http.StatusInternalServerError)
	}

	h.Logger.Debug("Login process completed successfully", "username", params.Username)

	return c.Redirect(http.StatusSeeOther, "/")
}

// LogoutGET handler for the logout page
func (h *Handler) LogoutGET(c echo.Context) error {
	// Clear the session
	if err := h.SessionManager.ClearSession(c, h.Config.Auth.SessionName); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Error clearing session", http.StatusInternalServerError)
	}

	// Set authenticated state to false
	if err := h.SessionManager.SetAuthenticated(c, false); err != nil {
		h.Logger.Error("Failed to set authenticated state", err)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

// RequestPasswordResetPOST handles the password reset request
func (h *Handler) RequestPasswordResetPOST(c echo.Context) error {
	dto := new(PasswordResetDTO)
	if err := c.Bind(dto); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid request", http.StatusBadRequest)
	}

	if err := h.Service.RequestPasswordReset(c.Request().Context(), dto); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Failed to process reset request", http.StatusInternalServerError)
	}

	sess, err := h.SessionManager.GetSession(c, h.Config.Auth.SessionName)
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Error getting session", http.StatusInternalServerError)
	}

	sess.AddFlash("Password reset instructions have been sent to your email")
	if err := h.SessionManager.SaveSession(c, sess); err != nil {
		h.Logger.Error("Failed to save session", err)
	}

	return c.Redirect(http.StatusSeeOther, "/user/login")
}

// ResetPasswordPOST handles the password reset completion
func (h *Handler) ResetPasswordPOST(c echo.Context) error {
	dto := new(ResetPasswordDTO)
	if err := c.Bind(dto); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Invalid request", http.StatusBadRequest)
	}

	if err := h.Service.ResetPassword(c.Request().Context(), dto); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Failed to reset password", http.StatusInternalServerError)
	}

	sess, err := h.SessionManager.GetSession(c, h.Config.Auth.SessionName)
	if err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Error getting session", http.StatusInternalServerError)
	}

	sess.AddFlash("Your password has been reset successfully")
	if err := h.SessionManager.SaveSession(c, sess); err != nil {
		h.Logger.Error("Failed to save session", err)
	}

	return c.Redirect(http.StatusSeeOther, "/user/login")
}

func (h *Handler) RequireAuthentication() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !h.SessionManager.IsAuthenticated(c) {
				return c.Redirect(http.StatusSeeOther, "/user/login")
			}
			return next(c)
		}
	}
}
