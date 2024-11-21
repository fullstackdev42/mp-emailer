package campaign

import (
	"github.com/jonesrussell/mp-emailer/config"
	"github.com/jonesrussell/mp-emailer/session"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers the campaign routes
func RegisterRoutes(h *Handler, e *echo.Echo, cfg *config.Config, sessionManager session.Manager) {
	// Public routes (no authentication required)
	e.GET("/campaigns", h.GetCampaigns)
	e.GET("/campaign/:id", h.CampaignGET)

	// Protected routes (require authentication)
	protected := e.Group("/campaign")
	protected.Use(sessionManager.ValidateSession(cfg.Auth.SessionName))

	// Protected campaign routes
	protected.GET("/new", h.CreateCampaignForm)
	protected.POST("", h.CreateCampaign)
	protected.GET("/:id/edit", h.EditCampaignForm)
	protected.PUT("/:id", h.EditCampaign)
	protected.DELETE("/:id", h.DeleteCampaign)
	protected.POST("/:id/compose", h.ComposeEmail)
	protected.POST("/:id/send", h.SendCampaign)

	// Debug logging
	for _, route := range e.Routes() {
		h.Logger.Debug("Registered route",
			"method", route.Method,
			"path", route.Path)
	}
}
