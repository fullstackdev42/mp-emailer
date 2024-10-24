package server

import (
	"net/http"

	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

// Handler struct
type Handler struct {
	Logger          loggo.LoggerInterface
	Store           sessions.Store
	emailService    email.Service
	templateManager *TemplateManager
	userService     user.ServiceInterface
	campaignService campaign.ServiceInterface
	errorHandler    *shared.ErrorHandler
}

// NewHandler creates a new Handler instance
func NewHandler(
	logger loggo.LoggerInterface,
	emailService email.Service,
	tmplManager *TemplateManager,
	userService user.ServiceInterface,
	campaignService campaign.ServiceInterface,
) *Handler {
	return &Handler{
		Logger:          logger,
		emailService:    emailService,
		templateManager: tmplManager,
		userService:     userService,
		campaignService: campaignService,
		errorHandler:    shared.NewErrorHandler(logger),
	}
}

// Home page handler
func (h *Handler) HandleIndex(c echo.Context) error {
	h.Logger.Debug("server.HandleIndex", "message", "Handling index request")
	isAuthenticated := c.Get("isAuthenticated").(bool)
	h.Logger.Debug("server.HandleIndex", "isAuthenticated", isAuthenticated)

	var internalError bool
	defer func() {
		if r := recover(); r != nil {
			h.Logger.Debug("server.HandleIndex: Panic recovered", "error", r)
			internalError = true
		}
	}()

	// Fetch campaigns using the campaign service
	campaigns, err := h.campaignService.GetAllCampaigns()
	if err != nil {
		h.Logger.Debug("server.HandleIndex: Error fetching campaigns", "error", err)
		return c.HTML(http.StatusInternalServerError, "<h1>Error fetching campaigns</h1>")
	}

	pageData := shared.PageData{
		Content:         campaigns,
		Title:           "Home",
		IsAuthenticated: isAuthenticated,
	}

	err = h.templateManager.Render(c.Response(), "home.html", pageData, c)
	if err != nil {
		h.Logger.Debug("server.HandleIndex: Error rendering template", "error", err)
		return c.HTML(http.StatusInternalServerError, "<h1>Error rendering page</h1>")
	}

	if internalError {
		return c.HTML(http.StatusInternalServerError, "<h1>Internal Server Error</h1>")
	}

	return nil
}
