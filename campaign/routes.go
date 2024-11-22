package campaign

import (
	"net/http"

	"github.com/jonesrussell/mp-emailer/session"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers the campaign routes
func RegisterRoutes(h *Handler, e *echo.Echo, sessionManager session.Manager) {
	// Public routes (no authentication required)
	e.GET("/campaigns", h.GetCampaigns)
	e.GET("/campaign/:id", h.CampaignGET)

	// Protected routes (require authentication)
	protected := e.Group("/campaign")
	protected.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := sessionManager.ValidateSession(c); err != nil {
				return c.Redirect(http.StatusSeeOther, "/user/login")
			}
			return next(c)
		}
	})

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
