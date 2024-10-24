package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/internal/database"
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
		fx.Provide(
			config.Load,
			newLogger,
			newDB,
			newSessionStore,
			func() email.Service { return email.NewMailpitEmailService("test@test.com", "test", nil) },
			provideTemplateFS,
			newEcho,
			provideTemplates,
			shared.NewTemplateRenderer,
		),
		shared.Module,
		user.Module,
		campaign.Module,
		server.Module,
		fx.Invoke(registerRoutes, startServer),
	)
	app.Run()
}

// Provide a new Echo instance
func newEcho() *echo.Echo {
	return echo.New()
}

// Provide a *template.Template to the fx container
func provideTemplates() (*template.Template, error) {
	templates, err := template.ParseGlob("web/templates/**/*.gohtml")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}
	return templates, nil
}

// Central function to register routes
func registerRoutes(
	e *echo.Echo,
	serverHandler *server.Handler,
	campaignHandler *campaign.Handler,
	userHandler *user.Handler,
	renderer *shared.CustomTemplateRenderer,
) {
	// Set the custom renderer
	e.Renderer = renderer

	// Middleware for logging
	e.Use(middleware.Logger())
	// Middleware for rate limiting
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	// Register server routes
	server.RegisterRoutes(serverHandler, e)
	// Register campaign routes
	campaign.RegisterRoutes(campaignHandler, e)
	// Register user routes
	user.RegisterRoutes(userHandler, e)

	// Serve static files from the "static" directory
	e.Static("/static", "web/public")
}

//go:embed web/templates/**/*
var templateFS embed.FS

// Provide templateFS to the fx container
func provideTemplateFS() embed.FS {
	return templateFS
}

func newLogger(cfg *config.Config) (loggo.LoggerInterface, error) {
	logLevel := cfg.GetLogLevel()
	logger, err := loggo.NewLogger("mp-emailer.log", logLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	logger.Debug("Logger initialized with DEBUG level") // Example debug log
	return logger, nil
}

func newDB(logger loggo.LoggerInterface, cfg *config.Config) (*database.DB, error) {
	logger.Info("Initializing database connection")
	dsn := cfg.DatabaseDSN()
	return database.NewDB(dsn, logger)
}

func newSessionStore(cfg *config.Config) sessions.Store {
	return sessions.NewCookieStore([]byte(cfg.SessionSecret))
}

func startServer(lc fx.Lifecycle, e *echo.Echo, config *config.Config, logger loggo.LoggerInterface) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				logger.Debug("Server starting")
				logger.Info(fmt.Sprintf("Starting server on :%s", config.AppPort))
				if err := e.Start(":" + config.AppPort); !errors.Is(err, http.ErrServerClosed) {
					logger.Error("Error starting server", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Debug("Server shutting down")
			if err := e.Shutdown(ctx); err != nil {
				logger.Error("Error shutting down server", err)
				return err
			}
			return nil
		},
	})
}

// HandlerResult is the output struct for NewHandler
type HandlerResult struct {
	fx.Out
	ServerHandler   *server.Handler
	CampaignHandler *campaign.Handler
}
