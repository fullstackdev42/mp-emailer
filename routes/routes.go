package routes

import (
	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/server"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterRoutes(e *echo.Echo, _ *server.Handler, ch *campaign.Handler, uh *user.Handler) {
	// Register user routes
	e.GET("/login", uh.LoginGET)
	e.POST("/login", uh.LoginPOST)
	e.GET("/logout", uh.LogoutGET)
	e.GET("/register", uh.RegisterGET)
	e.POST("/register", uh.RegisterPOST)

	// Register campaign routes
	e.GET("/campaigns", ch.GetAllCampaigns)
	e.GET("/campaigns/:id", ch.CampaignGET)
	e.POST("/campaigns/:id/send", ch.SendCampaign)
	e.GET("/campaigns/lookup-representatives", ch.HandleRepresentativeLookup)

	// Protected campaign routes
	e.GET("/campaigns/new", uh.RequireAuthMiddleware(ch.CreateCampaignForm))
	e.POST("/campaigns/new", uh.RequireAuthMiddleware(ch.CreateCampaign))
	e.GET("/campaigns/:id/edit", uh.RequireAuthMiddleware(ch.EditCampaignForm))
	e.POST("/campaigns/:id/edit", uh.RequireAuthMiddleware(ch.EditCampaign))
	e.DELETE("/campaigns/:id/delete", uh.RequireAuthMiddleware(ch.DeleteCampaign))

	// Apply SetAuthStatusMiddleware to all routes
	e.Use(user.SetAuthStatusMiddleware(uh.Store, uh.Logger, uh.SessionName))

	// Use Echo for middleware and other features
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
}
