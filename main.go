package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"

	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/internal/database"
	"github.com/fullstackdev42/mp-emailer/routes"
	"github.com/fullstackdev42/mp-emailer/server"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

//go:embed web/templates/* web/templates/partials/*
var templateFS embed.FS

func main() {
	app := fx.New(
		fx.Provide(
			config.Load,
			newLogger,
			newDB,
			email.NewEmailService,
			newTemplateManager,
			newSessionStore,
			campaign.NewRepository,
			campaign.NewService,
			campaign.NewRepresentativeLookupService,
			campaign.NewDefaultClient,
			campaign.NewHandler,
			user.NewRepository,
			user.NewService,
			user.NewHandler,
			server.New,
			provideHandler,
		),
		fx.Invoke(
			registerRoutes,
			startServer,
		),
	)
	app.Run()
}

func newLogger(config *config.Config) (loggo.LoggerInterface, error) {
	logger, err := loggo.NewLogger("mp-emailer.log", config.GetLogLevel())
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func newDB(logger loggo.LoggerInterface, config *config.Config) (*database.DB, error) {
	logger.Info("Initializing database connection")
	dsn := config.DatabaseDSN()            // Get the DSN from the config
	db, err := database.NewDB(dsn, logger) // Call NewDB with DSN and logger
	if err != nil {
		logger.Error("Failed to initialize database", err)
		return nil, err
	}
	return db, nil
}

func newTemplateManager() (*server.TemplateManager, error) {
	return server.NewTemplateManager(templateFS)
}

func newSessionStore(config *config.Config) sessions.Store {
	return sessions.NewCookieStore([]byte(config.SessionSecret))
}

func provideHandler(
	logger loggo.LoggerInterface,
	store sessions.Store,
	emailService email.Service,
	tmplManager *server.TemplateManager,
	userService user.ServiceInterface,
) *server.Handler {
	return server.NewHandler(logger.(*loggo.Logger), store, emailService, tmplManager, userService)
}

func registerRoutes(
	e *echo.Echo,
	handler *server.Handler,
	campaignHandler *campaign.Handler,
	userHandler *user.Handler,
) {
	routes.RegisterRoutes(e, handler, campaignHandler, userHandler)
}

func startServer(
	lc fx.Lifecycle,
	e *echo.Echo,
	config *config.Config,
	logger loggo.LoggerInterface,
) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				logger.Info(fmt.Sprintf("Starting server on :%s", config.AppPort))
				if err := e.Start(":" + config.AppPort); err != http.ErrServerClosed {
					logger.Error("Error starting server", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})
}
