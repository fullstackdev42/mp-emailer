package shared

import (
	"net/http"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
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
}

func NewBaseHandler(params BaseHandlerParams) BaseHandler {
	return BaseHandler{
		Store:            params.Store,
		Logger:           params.Logger,
		ErrorHandler:     params.ErrorHandler,
		Config:           params.Config,
		TemplateRenderer: params.TemplateRenderer,
	}
}

// Common rendering method that can be used by all handlers
func (h *BaseHandler) RenderTemplate(c echo.Context, name string, data interface{}) error {
	if err := h.TemplateRenderer.Render(c.Response().Writer, name, data, c); err != nil {
		return h.ErrorHandler.HandleHTTPError(c, err, "Failed to render "+name, http.StatusInternalServerError)
	}
	h.Logger.Debug(name + " template rendered successfully")
	return nil
}
