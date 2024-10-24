package routes

import (
	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/server"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, sh *server.Handler, ch *campaign.Handler, uh *user.Handler) {
	// Apply SetAuthStatusMiddleware to all routes
	e.Use(user.SetAuthStatusMiddleware(uh.Store, uh.Logger, uh.SessionName))

	// Public routes
	e.GET("/", sh.HandleIndex)

	// User routes
	registerUserRoutes(e, uh)

	// Campaign routes (both public and protected)
	registerCampaignRoutes(e, ch)

	// Protected routes
	authGroup := e.Group("/campaigns")
	authGroup.Use(user.RequireAuthMiddleware(uh.Store, uh.SessionName))

	// Protected campaign routes
	registerProtectedCampaignRoutes(authGroup, ch)
}

func registerUserRoutes(e *echo.Echo, uh *user.Handler) {
	// Public user routes
	e.GET("/login", uh.LoginGET)   // Handle GET request for login
	e.POST("/login", uh.LoginPOST) // Handle POST request for login

	e.GET("/logout", uh.LogoutGET)       // Handle GET request for logout
	e.GET("/register", uh.RegisterGET)   // Handle GET request for registration
	e.POST("/register", uh.RegisterPOST) // Handle POST request for registration
}

func registerCampaignRoutes(e *echo.Echo, ch *campaign.Handler) {
	// Public campaign routes
	e.GET("/campaigns", ch.GetAllCampaigns)
	e.GET("/campaigns/:id", ch.CampaignGET)
	e.POST("/campaigns/:id/send", ch.SendCampaign)
	e.POST("/campaigns/lookup-representatives", ch.HandleRepresentativeLookup)
}

func registerProtectedCampaignRoutes(authGroup *echo.Group, ch *campaign.Handler) {
	// Protected campaign routes
	authGroup.GET("/new", ch.CreateCampaignForm)
	authGroup.POST("/new", ch.CreateCampaign)
	authGroup.POST("/:id/delete", ch.DeleteCampaign)
	authGroup.GET("/:id/edit", ch.EditCampaignForm)
	authGroup.POST(":id/edit", ch.EditCampaign)
}
