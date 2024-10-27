package api

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(h *Handler, e *echo.Echo, jwtSecret string) {
	api := e.Group("/api")

	// Public routes
	api.POST("/users/register", h.RegisterUser)
	api.POST("/users/login", h.LoginUser)

	// Protected routes
	protected := api.Group("")
	protected.Use(JWTMiddleware(jwtSecret))

	// Campaign routes
	campaigns := protected.Group("/campaigns")
	campaigns.GET("", h.GetCampaigns)
	campaigns.GET("/:id", h.GetCampaign)
	campaigns.POST("", h.CreateCampaign)
	campaigns.PUT("/:id", h.UpdateCampaign)
	campaigns.DELETE("/:id", h.DeleteCampaign)

	// User routes
	users := protected.Group("/users")
	users.GET("/:username", h.GetUser)
}
