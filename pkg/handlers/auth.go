package handlers

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	adminEmail = "admin@example.com" // Replace with your actual admin email
)

func (h *Handler) handleError(err error, statusCode int, message string) error {
	h.logger.Error(message, err)
	return echo.NewHTTPError(statusCode, message)
}

func (h *Handler) HandleLogin(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "login.html", nil)
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	valid, err := h.db.VerifyUser(username, password)
	if err != nil {
		return h.handleError(err, http.StatusInternalServerError, "Error verifying user")
	}

	if !valid {
		data := map[string]interface{}{
			"Error": "Invalid username or password",
		}
		return c.Render(http.StatusUnauthorized, "login.html", data)
	}

	// Set user in session
	sess, err := session.Get("mpe", c)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get session")
	}
	sess.Values["username"] = username
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to save session")
	}

	return c.Redirect(http.StatusFound, "/")
}

func (h *Handler) HandleLogout(c echo.Context) error {
	sess, err := session.Get("mpe", c)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get session")
	}

	sess.Values["username"] = nil
	sess.Options.MaxAge = -1 // This will delete the cookie
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to save session")
	}

	// Clear the isAuthenticated status
	c.Set("isAuthenticated", false)

	return c.Redirect(http.StatusSeeOther, "/")
}
