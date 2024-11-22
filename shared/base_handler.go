package shared

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/mp-emailer/config"
	"github.com/jonesrussell/mp-emailer/session"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// HTTP Status codes
const (
	StatusOK                  = http.StatusOK
	StatusSeeOther            = http.StatusSeeOther
	StatusBadRequest          = http.StatusBadRequest
	StatusUnauthorized        = http.StatusUnauthorized
	StatusInternalServerError = http.StatusInternalServerError
)

// BaseHandlerParams defines the input parameters for BaseHandler
type BaseHandlerParams struct {
	fx.In

	Logger           loggo.LoggerInterface
	ErrorHandler     ErrorHandlerInterface
	Config           *config.Config
	TemplateRenderer TemplateRendererInterface
	SessionManager   session.Manager `optional:"true"`
}

type BaseHandler struct {
	Logger           loggo.LoggerInterface
	ErrorHandler     ErrorHandlerInterface
	Config           *config.Config
	TemplateRenderer TemplateRendererInterface
	SessionManager   session.Manager
	MapError         func(error) (int, string)
}

func NewBaseHandler(params BaseHandlerParams) BaseHandler {
	return BaseHandler{
		Logger:           params.Logger,
		ErrorHandler:     params.ErrorHandler,
		Config:           params.Config,
		TemplateRenderer: params.TemplateRenderer,
		SessionManager:   params.SessionManager,
		MapError:         DefaultErrorMapper,
	}
}

// DefaultErrorMapper provides default error mapping logic
func DefaultErrorMapper(_ error) (int, string) {
	return http.StatusInternalServerError, "Internal Server Error"
}

// Error implements HandlerLoggable interface
func (h *BaseHandler) Error(msg string, err error, keyvals ...interface{}) {
	h.Logger.Error(msg, err, keyvals...)
}

// Info implements HandlerLoggable interface
func (h *BaseHandler) Info(msg string, keyvals ...interface{}) {
	h.Logger.Info(msg, keyvals...)
}

// Warn implements HandlerLoggable interface
func (h *BaseHandler) Warn(msg string, keyvals ...interface{}) {
	h.Logger.Warn(msg, keyvals...)
}

// AddFlashMessage adds a flash message to the session
func (h *BaseHandler) AddFlashMessage(c echo.Context, message string) error {
	h.Logger.Debug("Adding flash message", "message", message)

	if h.SessionManager == nil {
		return session.ErrSessionNotFound
	}

	sess, err := h.SessionManager.GetSession(c, h.Config.Auth.SessionName)
	if err != nil {
		h.Logger.Error("Failed to get session", err)
		return err
	}

	sess.AddFlash(message, "messages")
	return h.SessionManager.SaveSession(c, sess)
}

// GetFlashMessages retrieves and clears flash messages from the session
func (h *BaseHandler) GetFlashMessages(c echo.Context) ([]string, error) {
	if h.SessionManager == nil {
		return nil, session.ErrSessionNotFound
	}

	sess, err := h.SessionManager.GetSession(c, h.Config.Auth.SessionName)
	if err != nil {
		h.Logger.Error("Failed to get session", err)
		return nil, err
	}

	flashes := sess.Flashes("messages")
	messages := make([]string, len(flashes))
	for i, flash := range flashes {
		if str, ok := flash.(string); ok {
			messages[i] = str
		}
	}

	if err := h.SessionManager.SaveSession(c, sess); err != nil {
		h.Logger.Error("Failed to save session after getting flashes", err)
		return messages, err
	}

	return messages, nil
}

// ClearSession clears the current session
func (h *BaseHandler) ClearSession(c echo.Context) error {
	if h.SessionManager == nil {
		return session.ErrSessionNotFound
	}

	return h.SessionManager.ClearSession(c, h.Config.Auth.SessionName)
}

// GetSession retrieves the current session
func (h *BaseHandler) GetSession(c echo.Context) (*sessions.Session, error) {
	if h.SessionManager == nil {
		return nil, session.ErrSessionNotFound
	}
	return h.SessionManager.GetSession(c, h.Config.Auth.SessionName)
}

// SaveSession saves the session
func (h *BaseHandler) SaveSession(c echo.Context, sess *sessions.Session) error {
	if h.SessionManager == nil {
		return session.ErrSessionNotFound
	}
	return h.SessionManager.SaveSession(c, sess)
}

// SetSessionValues sets user values in the session
func (h *BaseHandler) SetSessionValues(sess *sessions.Session, user interface{}) {
	if h.SessionManager != nil {
		h.SessionManager.SetSessionValues(sess, user)
	}
}

// IsAuthenticated checks if the current session is authenticated
func (h *BaseHandler) IsAuthenticated(c echo.Context) bool {
	if h.SessionManager == nil {
		return false
	}
	return h.SessionManager.IsAuthenticated(c)
}
