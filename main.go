package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"net/http"

	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/internal/database"
	"github.com/fullstackdev42/mp-emailer/server"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

//go:embed web/templates/layout/* web/templates/shared/* web/templates/partials/* web/templates/pages/*
var templateFS embed.FS

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
		),
		campaign.Module,
		server.ProvideModule(),
		user.ProvideModule(),
		fx.Invoke(registerRoutes, startServer),
	)
	app.Run()
}

// Provide a new Echo instance
func newEcho() *echo.Echo {
	return echo.New()
}

// Central function to register routes
func registerRoutes(e *echo.Echo, serverHandler *server.Handler, campaignHandler *campaign.Handler, userHandler *user.Handler) {
	// Register server routes
	server.RegisterRoutes(serverHandler, e)
	// Register campaign routes
	campaign.RegisterRoutes(campaignHandler, e)
	// Register user routes
	user.RegisterRoutes(userHandler, e)
	// Add more route registrations as needed
}

// Provide templateFS to the fx container
func provideTemplateFS() embed.FS {
	return templateFS
}

func newLogger(cfg *config.Config) (loggo.LoggerInterface, error) {
	return loggo.NewLogger("mp-emailer.log", cfg.GetLogLevel())
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
			return e.Shutdown(ctx)
		},
	})
}

// HandlerResult is the output struct for NewHandler
type HandlerResult struct {
	fx.Out
	ServerHandler   *server.Handler
	CampaignHandler *campaign.Handler
}

func HTTPErrorHandler(err error, c echo.Context, logger loggo.LoggerInterface, cfg *config.Config) {
	message := "Internal Server Error"
	statusCode := http.StatusInternalServerError
	if cfg.IsDevelopment() {
		message = err.Error()
	}
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		statusCode = httpErr.Code
		if msg, ok := httpErr.Message.(string); ok {
			message = msg
		}
	}
	c.Logger().Error(err)
	if err := c.Render(statusCode, "error", map[string]interface{}{
		"Title":   fmt.Sprintf("%d - %s", statusCode, http.StatusText(statusCode)),
		"Message": message,
	}); err != nil {
		logger.Error("Failed to render error page", err)
		if err := c.String(http.StatusInternalServerError, "Internal Server Error"); err != nil {
			logger.Error("Failed to send error response", err)
		}
	}
}
