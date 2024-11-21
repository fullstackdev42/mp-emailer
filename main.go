package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/mp-emailer/api"
	"github.com/jonesrussell/mp-emailer/campaign"
	"github.com/jonesrussell/mp-emailer/config"
	appMiddleware "github.com/jonesrussell/mp-emailer/middleware"
	"github.com/jonesrussell/mp-emailer/server"
	"github.com/jonesrussell/mp-emailer/session"
	"github.com/jonesrussell/mp-emailer/shared"
	"github.com/jonesrussell/mp-emailer/user"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func main() {
	// Initialize application using uber/fx dependency injection
	app := fx.New(
		fx.Options(
			shared.App,
			campaign.Module,
			user.Module,
			server.Module,
			api.Module,
			appMiddleware.Module,
			fx.Invoke(registerRoutes, startServer),
		),
	)

	app.Run()
}

// registerRoutes centralizes all route registration for the application
func registerRoutes(
	e *echo.Echo,
	serverHandler server.HandlerInterface,
	campaignHandler *campaign.Handler,
	userHandler *user.Handler,
	apiHandler *api.Handler,
	renderer shared.TemplateRendererInterface,
	middlewareManager *appMiddleware.Manager,
	cfg *config.Config,
	sessionManager session.Manager,
) {
	// Set custom template renderer for HTML responses
	e.Renderer = renderer

	// Register middleware first
	middlewareManager.Register(e)

	// Register route handlers after middleware
	registerHandlers(e, serverHandler, campaignHandler, userHandler, apiHandler, cfg, middlewareManager, sessionManager)

	// Serve static files
	e.Static("/static", "web/public")
}

// registerHandlers configures all route handlers for the application
func registerHandlers(
	e *echo.Echo,
	serverHandler server.HandlerInterface,
	campaignHandler *campaign.Handler,
	userHandler *user.Handler,
	apiHandler *api.Handler,
	cfg *config.Config,
	middlewareManager *appMiddleware.Manager,
	sessionManager session.Manager,
) {
	server.RegisterRoutes(serverHandler, e)
	campaign.RegisterRoutes(campaignHandler, e, cfg, sessionManager)
	user.RegisterRoutes(userHandler, e)
	api.RegisterRoutes(apiHandler, e, middlewareManager)
}

// startServer configures the server and starts it
func startServer(lc fx.Lifecycle, e *echo.Echo, cfg *config.Config, logger loggo.LoggerInterface, handler server.HandlerInterface) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			// Register health check endpoint
			e.GET("/health", handler.HealthCheck)

			go func() {
				addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
				logger.Info("Starting server", "host", cfg.App.Host, "port", cfg.App.Port)
				if err := e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Error("Server error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down server")

			// Set shutting down status for health checks
			if h, ok := handler.(*server.Handler); ok {
				h.IsShuttingDown = true
			}

			if err := e.Shutdown(ctx); err != nil {
				logger.Error("Error during shutdown", err)
				return err
			}

			logger.Info("Server shutdown completed")
			return nil
		},
	})
}
