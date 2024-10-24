package server

import (
	"fmt"
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
	templateManager TemplateRenderer
	userService     user.ServiceInterface
	campaignService campaign.ServiceInterface
	errorHandler    *shared.ErrorHandler
	EmailService    email.Service
}

// NewHandler creates a new Handler instance
func NewHandler(
	logger loggo.LoggerInterface,
	emailService email.Service,
	templateManager TemplateRenderer,
	userService user.ServiceInterface,
	campaignService campaign.ServiceInterface,
) *Handler {
	return &Handler{
		Logger:          logger,
		emailService:    emailService,
		templateManager: templateManager,
		userService:     userService,
		campaignService: campaignService,
		errorHandler:    shared.NewErrorHandler(logger),
	}
}

// Home page handler
func (h *Handler) HandleIndex(c echo.Context) error {
	h.Logger.Debug("Handling index request")
	isAuthenticated := false
	if auth, ok := c.Get("isAuthenticated").(bool); ok {
		isAuthenticated = auth
	}
	h.Logger.Debug("Authentication status", "isAuthenticated", isAuthenticated)

	// Fetch campaigns using the campaign service
	campaigns, err := h.campaignService.GetAllCampaigns()
	if err != nil {
		h.Logger.Error("Error fetching campaigns", err)
		return h.errorHandler.HandleHTTPError(c, err, "Error fetching campaigns", 500)
	}

	pageData := shared.PageData{
		Content:         campaigns,
		Title:           "Home",
		IsAuthenticated: isAuthenticated,
	}

	h.Logger.Debug("Attempting to render template", "template", "home")
	err = h.templateManager.Render(c.Response(), "home", pageData, c)
	if err != nil {
		h.Logger.Error("Error rendering template", err)
		return h.errorHandler.HandleHTTPError(c, err, "Error rendering page", 500)
	}

	h.Logger.Debug("Template rendered successfully")
	return nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to handle the request
	switch r.URL.Path {
	case "/":
		c := echo.New().NewContext(r, w)
		if err := h.HandleIndex(c); err != nil {
			h.Logger.Error("Error handling index", err)
			http.Error(w, fmt.Sprintf("Error handling index: %v", err), http.StatusInternalServerError)
		}
	}
}

// Pattern reports the path at which this handler is registered.
func (h *Handler) Pattern() string {
	return "/"
}
