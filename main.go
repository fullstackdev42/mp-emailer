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

	// Add this line to check if debug logging is working
	fmt.Printf("Log level: %v\n", cfg.GetLogLevel())

	logger, err := newLogger(cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}

	// Add these lines to test debug logging
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Error("This is an error message", fmt.Errorf("sample error"))

	app := fx.New(
		fx.Provide(
			func() (*config.Config, error) { return cfg, nil },
			newLogger,
			newDB,
			newTemplateManager,
			newSessionStore,
			newEcho,
			userRepositoryProvider,
			NewUserService,
			user.NewHandler,
			campaignRepositoryProvider,
			NewCampaignService,
			campaign.NewHandler,
			campaign.NewRepresentativeLookupService,
			campaign.NewDefaultClient,
			NewEmailService,
			NewHandler,
		),
		fx.Invoke(registerRoutes, startServer),
	)

	// Use the existing logger to log the start of the application
	logger.Debug("Starting application")

	app.Run()
}

func newLogger(cfg *config.Config) (loggo.LoggerInterface, error) {
	logLevel := cfg.GetLogLevel()
	logger, err := loggo.NewLogger("mp-emailer.log", logLevel)
	if err != nil {
		return nil, err
	}

	// Add this line to verify the log level
	fmt.Printf("Logger created with level: %v\n", logLevel)

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

func campaignRepositoryProvider(db *database.DB) (campaign.RepositoryInterface, error) {
	return campaign.NewRepository(db), nil
}

func registerRoutes(e *echo.Echo, handler *server.Handler, campaignHandler *campaign.Handler, userHandler *user.Handler, logger loggo.LoggerInterface) {
	logger.Debug("Registering routes")
	routes.RegisterRoutes(e, handler, campaignHandler, userHandler)
	logger.Debug("Routes registered successfully")
}

func startServer(lc fx.Lifecycle, e *echo.Echo, config *config.Config, logger loggo.LoggerInterface) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				logger.Debug("Server starting")
				logger.Info(fmt.Sprintf("Starting server on :%s", config.AppPort))
				if err := e.Start(":" + config.AppPort); err != http.ErrServerClosed {
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

type HandlerResult struct {
	fx.Out
	Handler *server.Handler
}

func NewHandler(
	logger loggo.LoggerInterface,
	emailService email.Service,
	tmplManager *server.TemplateManager,
	userService user.ServiceInterface,
	campaignService campaign.ServiceInterface,
) (HandlerResult, error) {
	logger.Debug("Creating new handler")
	handler := server.NewHandler(
		logger,
		emailService,
		tmplManager,
		userService,
		campaignService,
	)
	logger.Debug("Handler created successfully")
	return HandlerResult{Handler: handler}, nil
}

type CampaignServiceResult struct {
	fx.Out
	Service campaign.ServiceInterface
}

func NewCampaignService(repo campaign.RepositoryInterface) (CampaignServiceResult, error) {
	service, err := campaign.NewService(repo)
	if err != nil {
		return CampaignServiceResult{}, err
	}
	return CampaignServiceResult{Service: service}, nil
}

type UserServiceResult struct {
	fx.Out
	Service user.ServiceInterface
}

func NewUserService(repo user.RepositoryInterface, logger loggo.LoggerInterface) (UserServiceResult, error) {
	service, err := user.NewService(repo.(*user.Repository), logger)
	if err != nil {
		return UserServiceResult{}, err
	}
	return UserServiceResult{Service: service}, nil
}

type EmailServiceResult struct {
	fx.Out
	Service email.Service
}

func NewEmailService(cfg *config.Config, logger loggo.LoggerInterface) (EmailServiceResult, error) {
	service, err := email.New(cfg, logger)
	if err != nil {
		return EmailServiceResult{}, err
	}
	return EmailServiceResult{Service: service}, nil
}
