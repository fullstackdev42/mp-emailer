package server

import (
	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/pkg/handlers"
	appmid "github.com/fullstackdev42/mp-emailer/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, h *handlers.Handler, ch *campaign.Handler) {
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

	// Campaign routes
	registerCampaignRoutes(e, authGroup, ch)
}

func registerCampaignRoutes(e *echo.Echo, authGroup *echo.Group, ch *campaign.Handler) {
	// Public campaign route
	e.GET("/campaigns/:id", ch.GetCampaign)

	// Protected campaign routes
	authGroup.GET("/campaigns", ch.GetAllCampaigns)
	authGroup.GET("/campaigns/new", ch.CreateCampaignForm)
	authGroup.POST("/campaigns/new", ch.CreateCampaign)
	authGroup.POST("/campaigns/:id/delete", ch.DeleteCampaign)
	authGroup.GET("/campaigns/:id/edit", ch.EditCampaignForm)
	authGroup.POST("/campaigns/:id/edit", ch.EditCampaign)
}
