package api

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(h *Handler, e *echo.Echo, jwtSecret string) {
	api := e.Group("/api")

	// Public routes
	api.POST("/user/register", h.RegisterUser)
	api.POST("/user/login", h.LoginUser)

	// Protected routes
	protected := api.Group("")
	protected.Use(JWTMiddleware(jwtSecret))

	// Campaign routes
	campaigns := protected.Group("/campaign")
	campaigns.GET("", h.GetCampaigns)
	campaigns.GET("/:id", h.GetCampaign)
	campaigns.POST("", h.CreateCampaign)
	campaigns.PUT("/:id", h.UpdateCampaign)
	campaigns.DELETE("/:id", h.DeleteCampaign)

	// User routes
	users := protected.Group("/user")
	users.GET("/:username", h.GetUser)
}
