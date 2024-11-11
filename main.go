package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fullstackdev42/mp-emailer/api"
	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/database"
	appMiddleware "github.com/fullstackdev42/mp-emailer/middleware"
	"github.com/fullstackdev42/mp-emailer/server"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	// Check for required configuration before starting the application
	if err := config.CheckRequired(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// Initialize application using uber/fx dependency injection
	app := fx.New(
		fx.Options(
			shared.App,
			campaign.Module,
			user.Module,
			server.Module,
			api.Module,
			database.MigrationModule,
			appMiddleware.Module,
			fx.Invoke(registerRoutes, startServer),
		),
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ConsoleLogger{W: os.Stdout}
		}),
	)

	app.Run()
}

// registerRoutes centralizes all route registration for the application
// It takes in all necessary handlers and services via dependency injection
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
) {
	// Set custom template renderer for HTML responses
	e.Renderer = renderer

	// Register middleware first
	middlewareManager.Register(e)

	// Register route handlers after middleware
	registerHandlers(e, serverHandler, campaignHandler, userHandler, apiHandler, cfg, middlewareManager)

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
) {
	server.RegisterRoutes(serverHandler, e)
	campaign.RegisterRoutes(campaignHandler, e, cfg, middlewareManager)
	user.RegisterRoutes(userHandler, e)
	api.RegisterRoutes(apiHandler, e, middlewareManager)
}

// startServer configures the server and starts it
func startServer(lc fx.Lifecycle, e *echo.Echo, cfg *config.Config, logger loggo.LoggerInterface, handler server.HandlerInterface) {
	// Create a channel to signal shutdown
	shutdown := make(chan struct{})

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			// Register health check endpoint
			e.GET("/health", handler.HealthCheck)

			go func() {
				addr := fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort)
				logger.Info("Starting server", "host", cfg.AppHost, "port", cfg.AppPort)
				if err := e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Error("Server error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Initiating graceful shutdown")

			// Signal shutdown status
			close(shutdown)

			// Set shutting down status for health checks
			if h, ok := handler.(*server.Handler); ok {
				h.IsShuttingDown = true
			}

			// Create shutdown context with timeout
			shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			// Perform cleanup procedures
			cleanup := make(chan bool)
			go func() {
				// Add your cleanup procedures here
				logger.Info("Running cleanup procedures")

				// Example: Wait for active requests to complete
				time.Sleep(2 * time.Second)

				cleanup <- true
			}()

			// Wait for cleanup or timeout
			select {
			case <-cleanup:
				logger.Info("Cleanup completed successfully")
			case <-shutdownCtx.Done():
				logger.Warn("Cleanup timed out")
			}

			// Shutdown the server
			if err := e.Shutdown(shutdownCtx); err != nil {
				logger.Error("Error during shutdown", err)
				return err
			}

			logger.Info("Server shutdown completed")
			return nil
		},
	})
}
