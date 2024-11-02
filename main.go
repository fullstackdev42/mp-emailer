package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/fullstackdev42/mp-emailer/api"
	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/server"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		shared.App,
		campaign.Module,
		user.Module,
		server.Module,
		api.Module,
		fx.Invoke(registerRoutes, startServer),
	)

	app.Run()
}

// Central function to register routes
func registerRoutes(
	e *echo.Echo,
	serverHandler server.HandlerInterface,
	campaignHandler *campaign.Handler,
	userHandler *user.Handler,
	apiHandler *api.Handler,
	renderer *shared.CustomTemplateRenderer,
	sessionStore sessions.Store,
	cfg *config.Config,
) {
	// Set the custom renderer
	e.Renderer = renderer

	// Add session store to context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("store", sessionStore)
			return next(c)
		}
	})

	// Middleware for logging
	e.Use(middleware.Logger())

	// Middleware for rate limiting
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	// Add the auth middleware
	e.Use(user.AuthMiddleware(sessionStore, cfg))

	// Register routes
	server.RegisterRoutes(serverHandler, e)
	campaign.RegisterRoutes(campaignHandler, e)
	user.RegisterRoutes(userHandler, e)
	api.RegisterRoutes(apiHandler, e, cfg.JWTSecret)

	// Serve static files from the "static" directory
	e.Static("/static", "web/public")
}

// Start the server
func startServer(lc fx.Lifecycle, e *echo.Echo, config *config.Config, logger loggo.LoggerInterface) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				logger.Info(fmt.Sprintf("Starting server on :%s", config.AppPort))
				if err := e.Start(":" + config.AppPort); !errors.Is(err, http.ErrServerClosed) {
					logger.Error("Error starting server", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Server shutting down")
			if err := e.Shutdown(ctx); err != nil {
				logger.Error("Error shutting down server", err)
				return err
			}
			return nil
		},
	})
}
