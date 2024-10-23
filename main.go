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

//go:embed web/templates/shared/* web/templates/partials/* web/templates/pages/*
var templateFS embed.FS

func main() {
	app := fx.New(
		fx.Provide(
			campaign.NewDefaultClient,
			campaign.NewHandler,
			campaign.NewRepository,
			campaign.NewRepresentativeLookupService,
			campaign.NewService,
			config.Load,
			email.New,
			newDB,
			newLogger,
			newSessionStore,
			newTemplateManager,
			provideHandler,
			server.New,
			user.NewHandler,
			user.NewRepository,
			fx.Annotate(
				user.NewRepository,
				fx.As(new(user.RepositoryInterface)),
			),
			user.NewService,
		),
		fx.Invoke(registerRoutes, startServer),
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
	dsn := config.DatabaseDSN()
	db, err := database.NewDB(dsn, logger)
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
	emailService email.Service,
	tmplManager *server.TemplateManager,
	userService user.ServiceInterface,
	campaignService campaign.ServiceInterface,
) *server.Handler {
	return server.NewHandler(
		logger,
		emailService,
		tmplManager,
		userService,
		campaignService,
	)
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
