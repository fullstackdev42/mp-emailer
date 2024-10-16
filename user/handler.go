package user

import (
	"net/http"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service ServiceInterface
	logger  loggo.LoggerInterface
	config  *config.Config
}

func NewHandler(
	service ServiceInterface,
	logger loggo.LoggerInterface,
	config *config.Config,
) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
		config:  config,
	}
}

const internalServerError = "Internal server error"

func (h *Handler) RegisterGET(c echo.Context) error {
	sess, err := h.getSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internalServerError)
	}

	data := map[string]interface{}{}
	if flash := sess.Flashes(); len(flash) > 0 {
		data["Error"] = flash[0]
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			h.logger.Error("Failed to save session", err)
			return echo.NewHTTPError(http.StatusInternalServerError, internalServerError)
		}
	}

	return c.Render(http.StatusOK, "register.html", data)
}

func (h *Handler) RegisterPOST(c echo.Context) error {
	sess, err := session.Get(h.config.SessionName, c)
	if err != nil {
		h.logger.Error("Failed to get session", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")

	if username == "" || password == "" || email == "" {
		h.logger.Warn("Missing required fields", "username", username, "email", email)
		sess.AddFlash("Username, password, and email are required")
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			h.logger.Error("Failed to save session", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
		return c.Redirect(http.StatusSeeOther, "/register")
	}

	if err := h.service.RegisterUser(username, email, password); err != nil {
		h.logger.Error("Failed to register user", err)
		sess.AddFlash("Failed to register user. Please try again.")
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			h.logger.Error("Failed to save session", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
		return c.Redirect(http.StatusSeeOther, "/register")
	}

	h.logger.Info("User registered successfully", "username", username)
	return c.Redirect(http.StatusSeeOther, "/login")
}

// Handler for GET requests
func (h *Handler) LoginGET(c echo.Context) error {
	h.logger.Debug("Handling login GET request", "path", c.Path())

	session, err := h.getSession(c)
	if err != nil {
		h.logger.Error("Failed to get session", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	if auth, _ := session.Values["authenticated"].(bool); auth {
		h.logger.Debug("User already authenticated, redirecting to home", "path", c.Path())
		return c.Redirect(http.StatusSeeOther, "/")
	}

	return c.Render(http.StatusOK, "login.html", nil)
}

// Handler for POST requests
func (h *Handler) LoginPOST(c echo.Context) error {
	h.logger.Debug("Handling login request", map[string]interface{}{"path": c.Path()})

	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		h.logger.Warn("Empty username or password", "username", username)
		return echo.NewHTTPError(http.StatusBadRequest, "Username and password are required")
	}

	userID, err := h.service.VerifyUser(username, password)
	if err != nil {
		h.logger.Warn("Invalid login attempt", "username", username, "error", err)
		session, err := h.getSession(c)
		if err != nil {
			h.logger.Error("Failed to get session", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
		session.AddFlash("Invalid username or password", "error")
		if err := session.Save(c.Request(), c.Response()); err != nil {
			h.logger.Error("Failed to save session", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	session, err := h.getSession(c)
	if err != nil {
		h.logger.Error("Failed to get session", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	session.Values["authenticated"] = true
	session.Values["userID"] = userID
	session.Values["username"] = username
	if err := session.Save(c.Request(), c.Response()); err != nil {
		h.logger.Error("Failed to save session", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	h.logger.Debug("User logged in successfully", map[string]interface{}{"username": username})
	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

func (h *Handler) LogoutGET(c echo.Context) error {
	h.logger.Debug("Handling logout request", map[string]interface{}{"path": c.Path()})

	session, err := h.getSession(c)
	if err != nil {
		h.logger.Error("Failed to get session", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	// Clear session
	session.Options.MaxAge = -1
	session.Values["userID"] = nil
	session.Values["username"] = nil

	if err := session.Save(c.Request(), c.Response()); err != nil {
		h.logger.Error("Failed to save session", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

// Helper method to get session
func (h *Handler) getSession(c echo.Context) (*sessions.Session, error) {
	sess, err := session.Get(h.config.SessionName, c)
	if err != nil {
		h.logger.Error("Failed to get session", err)
		return nil, err
	}
	return sess, nil
}
