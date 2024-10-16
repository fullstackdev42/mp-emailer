package server

import (
	"net/http"

	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Logger          loggo.LoggerInterface
	Store           sessions.Store
	emailService    email.Service
	templateManager *TemplateManager
	userService     user.ServiceInterface
}

func NewHandler(
	logger loggo.LoggerInterface,
	emailService email.Service,
	tmplManager *TemplateManager,
	userService user.ServiceInterface,
) *Handler {
	return &Handler{
		Logger:          logger,
		emailService:    emailService,
		templateManager: tmplManager,
		userService:     userService,
	}
}

func (h *Handler) HandleIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}
