package campaign

import (
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers the campaign routes
func RegisterRoutes(h *Handler, e *echo.Echo, cfg *config.Config, logger loggo.LoggerInterface) {
	// Public routes (no authentication required)
	e.GET("/campaigns", h.GetCampaigns)
	e.GET("/campaign/:id", h.CampaignGET)

	// Protected routes (require authentication)
	protected := e.Group("/campaign")
	protected.Use(ValidateSessionWithLogging(cfg.SessionName, logger))

	// Protected campaign routes
	protected.GET("/new", h.CreateCampaignForm)
	protected.POST("", h.CreateCampaign)
	protected.GET("/:id/edit", h.EditCampaignForm)
	protected.PUT("/:id", h.EditCampaign)
	protected.DELETE("/:id", h.DeleteCampaign)
	protected.POST("/:id/compose", h.ComposeEmail)
	protected.POST("/:id/send", h.SendCampaign)
}
