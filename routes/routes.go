package routes

import (
	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/server"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, h *server.Handler, ch *campaign.Handler, uh *user.Handler) {
	// Public routes
	e.GET("/", h.HandleIndex)

	// User routes
	registerUserRoutes(e, uh)

	// Protected routes
	authGroup := e.Group("")
	authGroup.Use(user.RequireAuthMiddleware(h.Store))
	authGroup.POST("/echo", h.HandleEcho)

	// Campaign routes
	registerCampaignRoutes(e, authGroup, ch)
}

func registerUserRoutes(e *echo.Echo, uh *user.Handler) {
	// Public user routes
	e.GET("/login", uh.LoginGET)
	e.POST("/login", uh.LoginPOST)

	e.GET("/logout", uh.HandleLogout)
	e.GET("/register", uh.RegisterGET)
	e.POST("/register", uh.RegisterPOST)

	// Add any additional user-related routes here
	// For example:
	// e.GET("/profile", uh.HandleProfile)
	// e.POST("/profile/update", uh.HandleProfileUpdate)
}

func registerCampaignRoutes(e *echo.Echo, authGroup *echo.Group, ch *campaign.Handler) {
	// Public campaign route
	e.GET("/campaigns/:id", ch.GetCampaign)
	e.POST("/campaigns/:id/send", ch.SendCampaign)

	// Protected campaign routes
	authGroup.GET("/campaigns", ch.GetAllCampaigns)
	authGroup.GET("/campaigns/new", ch.CreateCampaignForm)
	authGroup.POST("/campaigns/new", ch.CreateCampaign)
	authGroup.POST("/campaigns/:id/delete", ch.DeleteCampaign)
	authGroup.GET("/campaigns/:id/edit", ch.EditCampaignForm)
	authGroup.POST("/campaigns/:id/edit", ch.EditCampaign)

	// Add this new route
	e.POST("/lookup-representatives", ch.HandleRepresentativeLookup)
}
