package handlers

import (
	"errors"
	"net/http"

	"github.com/fullstackdev42/mp-emailer/pkg/database"
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
	h.logger.Debug("HandleLogin called with method: " + c.Request().Method)

	if c.Request().Method == http.MethodGet {
		h.logger.Debug("Rendering login page")
		return c.Render(http.StatusOK, "login.html", nil)
	}

	username := c.FormValue("username")
	password := c.FormValue("password")
	h.logger.Debug("Login attempt for username: " + username)

	// Retrieve the database connection from the context
	db := c.Get("db").(*database.DB)
	if db == nil {
		h.logger.Error("Database connection not found in context", errors.New("database connection not available"))
		return h.handleError(nil, http.StatusInternalServerError, "Database connection not available")
	}

	userID, err := db.VerifyUser(username, password)
	if err != nil {
		h.logger.Error("User verification failed", err)
		data := map[string]interface{}{
			"Error": err.Error(),
		}
		return c.Render(http.StatusUnauthorized, "login.html", data)
	}
	h.logger.Debug("User verified successfully. UserID: " + userID)

	// Set user in session
	sess, err := session.Get("mpe", c)
	if err != nil {
		h.logger.Error("Failed to get session", err)
		return c.String(http.StatusInternalServerError, "Failed to get session")
	}

	sess.Values["userID"] = userID
	sess.Values["username"] = username
	h.logger.Debug("Session values set. UserID: " + userID + ", Username: " + username)

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		h.logger.Error("Failed to save session", err)
		return c.String(http.StatusInternalServerError, "Failed to save session")
	}
	h.logger.Debug("Session saved successfully")

	h.logger.Debug("Redirecting to home page")
	return c.Redirect(http.StatusFound, "/")
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
