package shared

import (
	"errors"
	"net/http"

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
		h.Logger.Error("Session manager not initialized", errors.New("nil session manager"))
		return errors.New("session manager not initialized")
	}

	session, err := h.SessionManager.GetSession(c, h.Config.Auth.SessionName)
	if err != nil {
		h.Logger.Error("Failed to get session", err)
		return err
	}

	session.AddFlash(message, "messages")
	if err := h.SessionManager.SaveSession(c, session); err != nil {
		h.Logger.Error("Failed to save session after adding flash", err)
		return err
	}

	h.Logger.Debug("Flash message added successfully")
	return nil
}

// GetFlashMessages retrieves and clears flash messages from the session
func (h *BaseHandler) GetFlashMessages(c echo.Context) ([]string, error) {
	if h.SessionManager == nil {
		h.Logger.Error("Session manager not initialized", errors.New("nil session manager"))
		return nil, errors.New("session manager not initialized")
	}

	session, err := h.SessionManager.GetSession(c, h.Config.Auth.SessionName)
	if err != nil {
		h.Logger.Error("Failed to get session", err)
		return nil, err
	}

	flashes := session.Flashes("messages")
	messages := make([]string, len(flashes))
	for i, flash := range flashes {
		if str, ok := flash.(string); ok {
			messages[i] = str
		}
	}

	if err := h.SessionManager.SaveSession(c, session); err != nil {
		h.Logger.Error("Failed to save session after getting flashes", err)
		return messages, err
	}

	return messages, nil
}
