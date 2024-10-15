package user

import (
	"net/http"

	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service ServiceInterface
	logger  loggo.LoggerInterface
}

func NewHandler(service ServiceInterface, logger loggo.LoggerInterface) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) HandleRegister(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "register.html", nil)
	}

	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	err := h.service.RegisterUser(username, email, password)
	if err != nil {
		h.logger.Error("Registration error", err)
		return c.Render(http.StatusBadRequest, "register.html", map[string]interface{}{
			"Error": err.Error(),
		})
	}

	return c.Redirect(http.StatusSeeOther, "/login")
}

func (h *Handler) HandleLogin(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "login.html", nil)
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	h.logger.Debug("HandleLogin called with method: POST")
	h.logger.Debug("Login attempt for username: " + username)

	userID, err := h.service.VerifyUser(username, password)
	if err != nil {
		h.logger.Warn("Login failed for user: " + username)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
	}

	sess, err := session.Get("mpe", c)
	if err != nil {
		h.logger.Error("Failed to get session", err)
		return c.String(http.StatusInternalServerError, "Failed to get session")
	}

	sess.Values["userID"] = userID
	sess.Values["username"] = username
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		h.logger.Error("Failed to save session", err)
		return c.String(http.StatusInternalServerError, "Failed to save session")
	}

	return c.Redirect(http.StatusSeeOther, "/campaigns")
}

func (h *Handler) HandleLogout(c echo.Context) error {
	sess, err := session.Get("mpe", c)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get session")
	}

	sess.Values["userID"] = nil
	sess.Values["username"] = nil
	sess.Options.MaxAge = -1 // This will delete the cookie
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to save session")
	}

	// Clear the isAuthenticated status
	c.Set("isAuthenticated", false)

	return c.Redirect(http.StatusSeeOther, "/")
}
