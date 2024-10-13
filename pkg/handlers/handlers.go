package handlers

import (
	"net/http"

	"github.com/fullstackdev42/mp-emailer/pkg/api"
	"github.com/fullstackdev42/mp-emailer/pkg/services"
	"github.com/fullstackdev42/mp-emailer/pkg/templates"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Logger          loggo.LoggerInterface
	client          api.ClientInterface
	Store           sessions.Store
	emailService    services.EmailService
	templateManager *templates.TemplateManager
}

func NewHandler(logger loggo.LoggerInterface, client api.ClientInterface, store sessions.Store, emailService services.EmailService, tmplManager *templates.TemplateManager) *Handler {

	return &Handler{
		Logger:          logger,
		client:          client,
		Store:           store,
		emailService:    emailService,
		templateManager: tmplManager,
	}
}

func (h *Handler) HandleIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func (h *Handler) handleError(err error, statusCode int, message string) error {
	h.Logger.Error(message, err)
	return echo.NewHTTPError(statusCode, message)
}

func (h *Handler) HandleEcho(c echo.Context) error {
	type EchoRequest struct {
		Message string `json:"message"`
	}

	req := new(EchoRequest)
	if err := c.Bind(req); err != nil {
		return h.handleError(err, http.StatusBadRequest, "Error binding request")
	}

	return c.JSON(http.StatusOK, req)
}
