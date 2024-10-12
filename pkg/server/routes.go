package server

import (
	"github.com/fullstackdev42/mp-emailer/pkg/handlers"
	appmid "github.com/fullstackdev42/mp-emailer/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, h *handlers.Handler) {
	// Public routes
	e.GET("/", h.HandleIndex)
	e.GET("/login", h.HandleLogin)
	e.POST("/login", h.HandleLogin)
	e.GET("/logout", h.HandleLogout)
	e.GET("/register", h.HandleRegister)
	e.POST("/register", h.HandleRegister)

	// Protected routes
	authGroup := e.Group("")
	authGroup.Use(appmid.RequireAuthMiddleware(h.Store, h.Logger))
	authGroup.GET("/submit", h.HandleSubmit)
	authGroup.POST("/submit", h.HandleSubmit)
	authGroup.POST("/echo", h.HandleEcho)

	// Campaign routes (protected)
	authGroup.GET("/campaigns", h.HandleGetCampaigns)
	authGroup.GET("/campaigns/new", h.HandleCreateCampaign)
	authGroup.POST("/campaigns/new", h.HandleCreateCampaign)
	authGroup.GET("/campaigns/:id", h.HandleGetCampaign)
	authGroup.POST("/campaigns/:id/delete", h.HandleDeleteCampaign)
	e.GET("/campaigns/:id/edit", h.HandleEditCampaign)
	e.POST("/campaigns/:id/edit", h.HandleEditCampaign)
}
