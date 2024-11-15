package shared

import (
	"net/http"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"go.uber.org/fx"
)

// BaseHandlerParams defines the input parameters for BaseHandler
type BaseHandlerParams struct {
	fx.In

	Store            sessions.Store
	Logger           loggo.LoggerInterface
	ErrorHandler     ErrorHandlerInterface
	Config           *config.Config
	TemplateRenderer TemplateRendererInterface
}

type BaseHandler struct {
	Store            sessions.Store
	Logger           loggo.LoggerInterface
	ErrorHandler     ErrorHandlerInterface
	Config           *config.Config
	TemplateRenderer TemplateRendererInterface
	MapError         func(error) (int, string)
}

func NewBaseHandler(params BaseHandlerParams) BaseHandler {
	return BaseHandler{
		Store:            params.Store,
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
