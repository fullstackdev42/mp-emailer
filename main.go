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
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	app := fx.New(
		fx.Provide(
			func() (*config.Config, error) { return cfg, nil },
			newLogger,
			newDB,
			newTemplateManager,
			newSessionStore,
			newEcho,
			userRepositoryProvider,
			userServiceProvider,
			user.NewHandler,
			campaignRepositoryProvider,
			campaignServiceProvider,
			campaign.NewHandler,
			campaign.NewRepresentativeLookupService,
			campaign.NewDefaultClient,
			email.New,
			provideHandler,
		),
		fx.Invoke(registerRoutes, startServer),
	)
	app.Run()
}

func newLogger(cfg *config.Config) (loggo.LoggerInterface, error) {
	logger, err := loggo.NewLogger("mp-emailer.log", cfg.GetLogLevel())
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func newDB(logger loggo.LoggerInterface, cfg *config.Config) (*database.DB, error) {
	logger.Info("Initializing database connection")
	dsn := cfg.DatabaseDSN()
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

func newSessionStore(cfg *config.Config) sessions.Store {
	return sessions.NewCookieStore([]byte(cfg.SessionSecret))
}

func newEcho() *echo.Echo {
	return echo.New()
}

func userRepositoryProvider(db *database.DB, logger loggo.LoggerInterface) (user.RepositoryInterface, error) {
	return user.NewRepository(db, logger), nil
}

func userServiceProvider(repo user.RepositoryInterface, logger loggo.LoggerInterface) (user.ServiceInterface, error) {
	return user.NewService(repo.(*user.Repository), logger), nil
}

func campaignRepositoryProvider(db *database.DB) (campaign.RepositoryInterface, error) {
	return campaign.NewRepository(db), nil
}

func campaignServiceProvider(repo campaign.RepositoryInterface) (campaign.ServiceInterface, error) {
	return campaign.NewService(repo), nil
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

func registerRoutes(e *echo.Echo, handler *server.Handler, campaignHandler *campaign.Handler, userHandler *user.Handler) {
	routes.RegisterRoutes(e, handler, campaignHandler, userHandler)
}

func startServer(lc fx.Lifecycle, e *echo.Echo, config *config.Config, logger loggo.LoggerInterface) {
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
