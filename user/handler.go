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

func (h *Handler) withSession(f func(c echo.Context, sess *sessions.Session) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := h.getSession(c)
		if err != nil {
			return h.handleError("Failed to get session", err)
		}
		return f(c, sess)
	}
}

func (h *Handler) RegisterGET(c echo.Context) error {
	return h.withSession(func(c echo.Context, sess *sessions.Session) error {
		data := map[string]interface{}{}
		if flashMessages := getSessionFlashes(sess); len(flashMessages) > 0 {
			data["Error"] = flashMessages[0]
			if err := saveSession(sess, c.Request(), c.Response(), h.logger); err != nil {
				return h.handleError("Failed to save session", err)
			}
		}
		return c.Render(http.StatusOK, "register.html", data)
	})(c)
}

func (h *Handler) RegisterPOST(c echo.Context) error {
	sess, err := h.getSession(c)
	if err != nil {
		return h.handleError("Failed to get session", err)
	}

	username, password, email := c.FormValue("username"), c.FormValue("password"), c.FormValue("email")
	if missingRequiredFields(username, password, email) {
		return h.handleMissingFields(c, sess)
	}

	if err := h.service.RegisterUser(username, email, password); err != nil {
		return h.handleUserRegistrationFailure(c, sess, err)
	}

	h.logger.Info("User registered successfully", "username", username)
	return c.Redirect(http.StatusSeeOther, "/login")
}

// Handler for GET requests
func (h *Handler) LoginGET(c echo.Context) error {
	return h.withSession(func(c echo.Context, sess *sessions.Session) error {
		if auth, _ := sess.Values["authenticated"].(bool); auth {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		data := map[string]interface{}{}
		if flashMessages := sess.Flashes(); len(flashMessages) > 0 {
			data["Error"] = flashMessages[0]
			if err := sess.Save(c.Request(), c.Response()); err != nil {
				return h.handleError("Failed to save session", err)
			}
		}

		return c.Render(http.StatusOK, "login.html", data)
	})(c)
}

// Handler for POST requests
func (h *Handler) LoginPOST(c echo.Context) error {
	return h.withSession(func(c echo.Context, sess *sessions.Session) error {
		h.logger.Debug("Handling login request", map[string]interface{}{"path": c.Path()})

		username := c.FormValue("username")
		password := c.FormValue("password")

		if username == "" || password == "" {
			h.logger.Warn("Empty username or password", "username", username)
			sess.AddFlash("Username and password are required")
			if err := sess.Save(c.Request(), c.Response()); err != nil {
				return h.handleError("Failed to save session", err)
			}
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		userID, err := h.service.VerifyUser(username, password)
		if err != nil {
			h.logger.Warn("Invalid login attempt", "username", username, "error", err)
			sess.AddFlash("Invalid username or password")
			if err := sess.Save(c.Request(), c.Response()); err != nil {
				return h.handleError("Failed to save session", err)
			}
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		sess.Values["authenticated"] = true
		sess.Values["userID"] = userID
		sess.Values["username"] = username
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return h.handleError("Failed to save session", err)
		}

		h.logger.Debug("User logged in successfully", map[string]interface{}{"username": username})
		return c.Redirect(http.StatusSeeOther, "/campaigns")
	})(c)
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

func (h *Handler) handleMissingFields(c echo.Context, sess *sessions.Session) error {
	h.logger.Warn("Missing required fields")
	sess.AddFlash("Username, password, and email are required")
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return h.handleError("Failed to save session", err)
	}
	return c.Redirect(http.StatusSeeOther, "/register")
}

func (h *Handler) handleUserRegistrationFailure(c echo.Context, sess *sessions.Session, err error) error {
	h.logger.Error("Failed to register user", err)
	sess.AddFlash("Failed to register user. Please try again.")
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return h.handleError("Failed to save session", err)
	}
	return c.Redirect(http.StatusSeeOther, "/register")
}

func (h *Handler) handleError(message string, err error) error {
	h.logger.Error(message, err)
	return echo.NewHTTPError(http.StatusInternalServerError, internalServerError)
}

func getSessionFlashes(sess *sessions.Session) []interface{} {
	return sess.Flashes()
}

func saveSession(sess *sessions.Session, req *http.Request, res http.ResponseWriter, logger loggo.LoggerInterface) error {
	if err := sess.Save(req, res); err != nil {
		logger.Error("Failed to save session", err)
		return err
	}
	return nil
}

func missingRequiredFields(username, password, email string) bool {
	return username == "" || password == "" || email == ""
}

// Helper method to get session
func (h *Handler) getSession(c echo.Context) (*sessions.Session, error) {
	sess, err := session.Get(h.config.SessionName, c)
	if err != nil {
		return nil, err
	}
	return sess, nil
}
