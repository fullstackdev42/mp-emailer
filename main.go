package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/fullstackdev42/mp-emailer/api"
	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/internal/database"
	"github.com/fullstackdev42/mp-emailer/server"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/go-playground/validator/v10"
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
			provideValidator,
			func() email.Service { return email.NewMailpitEmailService("test@test.com", "test", nil) },
			provideTemplateFS,
			newEcho,
			provideTemplates,
			fx.Annotated{
				Name: "representativeLookupBaseURL",
				Target: func(cfg *config.Config) string {
					return cfg.RepresentativeLookupBaseURL
				},
			},
			fx.Annotated{
				Name: "representativeLogger",
				Target: func(logger loggo.LoggerInterface) loggo.LoggerInterface {
					return logger
				},
			},
		),
		shared.Module,
		user.Module,
		campaign.Module,
		server.Module,
		api.Module,
		fx.Invoke(registerRoutes, startServer),
	)
	app.Run()
}

// Central function to register routes
func registerRoutes(
	e *echo.Echo,
	serverHandler *server.Handler,
	campaignHandler *campaign.Handler,
	userHandler *user.Handler,
	apiHandler *api.Handler,
	renderer shared.TemplateRenderer,
	sessionStore sessions.Store,
	cfg *config.Config,
) {
	// Set the custom renderer
	e.Renderer = renderer

	// Middleware for logging
	e.Use(middleware.Logger())

	// Middleware for rate limiting
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	// Add the auth middleware
	e.Use(user.AuthMiddleware(sessionStore, cfg))

	// Register server routes
	server.RegisterRoutes(serverHandler, e)

	// Register campaign routes
	campaign.RegisterRoutes(campaignHandler, e)

	// Register user routes
	user.RegisterRoutes(userHandler, e)

	// Register API routes with JWT secret
	api.RegisterRoutes(apiHandler, e, cfg.JWTSecret)

	// Serve static files from the "static" directory
	e.Static("/static", "web/public")
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

//go:embed web/templates/**/*
var templateFS embed.FS

// Provide templateFS to the fx container
func provideTemplateFS() embed.FS {
	return templateFS
}

// Provide a *template.Template to the fx container
func provideTemplates() *template.Template {
	// Load your templates here
	return template.Must(template.ParseGlob("web/templates/**/*.gohtml"))
}

// Provide a new Echo instance
func newEcho() *echo.Echo {
	return echo.New()
}

// Provide a new logger
func newLogger(cfg *config.Config) (loggo.LoggerInterface, error) {
	logLevel := cfg.GetLogLevel()
	logger, err := loggo.NewLogger(cfg.LogFile, logLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	logger.Debug("Logger initialized with level: %s", logLevel)
	return logger, nil
}

func newDB(logger loggo.LoggerInterface, cfg *config.Config) (*database.DB, error) {
	logger.Info("Initializing database connection")
	dsn := cfg.DatabaseDSN()
	var db *database.DB
	var err error
	for retries := 5; retries > 0; retries-- {
		db, err = database.NewDB(dsn, logger)
		if err == nil {
			return db, nil
		}
		logger.Warn("Failed to connect to database, retrying...", "error", err)
		time.Sleep(5 * time.Second)
	}
	return nil, fmt.Errorf("failed to connect to database after multiple attempts: %w", err)
}

func newSessionStore(cfg *config.Config) sessions.Store {
	return sessions.NewCookieStore([]byte(cfg.SessionSecret))
}

func provideValidator() *validator.Validate {
	return validator.New()
}
