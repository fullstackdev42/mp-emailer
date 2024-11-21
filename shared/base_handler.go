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

// BaseHandlerParams defines the input parameters for BaseHandler
type BaseHandlerParams struct {
	fx.In

	Logger           loggo.LoggerInterface
	ErrorHandler     ErrorHandlerInterface
	Config           *config.Config
	TemplateRenderer TemplateRendererInterface
}

type BaseHandler struct {
	Logger           loggo.LoggerInterface
	ErrorHandler     ErrorHandlerInterface
	Config           *config.Config
	TemplateRenderer TemplateRendererInterface
	MapError         func(error) (int, string)
}

func NewBaseHandler(params BaseHandlerParams) BaseHandler {
	return BaseHandler{
		Logger:           params.Logger,
		ErrorHandler:     params.ErrorHandler,
		Config:           params.Config,
		TemplateRenderer: params.TemplateRenderer,
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

	sessionManager, ok := c.Get("session_manager").(session.Manager)
	if !ok {
		h.Logger.Error("Failed to get session manager from context", errors.New("session manager not found"))
		return errors.New("session manager not found")
	}

	session, err := sessionManager.GetSession(c, h.Config.Auth.SessionName)
	if err != nil {
		h.Logger.Error("Failed to get session", err)
		return err
	}

	session.AddFlash(message, "messages")
	if err := sessionManager.SaveSession(c, session); err != nil {
		h.Logger.Error("Failed to save session after adding flash", err)
		return err
	}

	h.Logger.Debug("Flash message added successfully")
	return nil
}
