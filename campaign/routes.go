package campaign

import (
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers the campaign routes
func RegisterRoutes(h *Handler, e *echo.Echo, cfg *config.Config) {
	// Public routes (no authentication required)
	e.GET("/campaigns", h.GetCampaigns)
	e.GET("/campaign/:id", h.CampaignGET)
	e.POST("/campaign/:id/compose", h.ComposeEmail)
	e.POST("/campaign/:id/send", h.SendCampaign)

	// Protected routes (require authentication)
	protected := e.Group("/campaign")
	protected.Use(ValidateSession(cfg.SessionName))
	protected.GET("/new", h.CreateCampaignForm)
	protected.POST("", h.CreateCampaign)
	protected.GET("/:id/edit", h.EditCampaignForm)
	protected.PUT("/:id", h.EditCampaign)
	protected.DELETE("/:id", h.DeleteCampaign)
}
