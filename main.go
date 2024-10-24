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
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

//go:embed web/templates/shared/* web/templates/partials/* web/templates/pages/*
var templateFS embed.FS

func main() {
	app := fx.New(
		fx.Provide(
			config.Load,
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

	app.Run()
}

func newLogger(cfg *config.Config) (loggo.LoggerInterface, error) {
	return loggo.NewLogger("mp-emailer.log", cfg.GetLogLevel())
}

func newDB(logger loggo.LoggerInterface, cfg *config.Config) (*database.DB, error) {
	logger.Info("Initializing database connection")
	dsn := cfg.DatabaseDSN()
	return database.NewDB(dsn, logger)
}

func newTemplateManager() (*server.TemplateManager, error) {
	return server.NewTemplateManager(templateFS)
}

func newSessionStore(cfg *config.Config) sessions.Store {
	return sessions.NewCookieStore([]byte(cfg.SessionSecret))
}

func newEcho(cfg *config.Config, logger loggo.LoggerInterface, tmplManager *server.TemplateManager) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Debug = cfg.IsDevelopment()

	// Set the renderer
	e.Renderer = tmplManager

	// Custom error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		HTTPErrorHandler(err, c, logger, cfg)
	}

	// Add middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	// Add static file serving
	e.Static("/static", "web/public")

	return e
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

// HandlerResult is the output struct for NewHandler
type HandlerResult struct {
	fx.Out
	Handler *server.Handler
}

// NewHandler creates a new server.Handler
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

// CampaignServiceResult is the output struct for NewCampaignService
type CampaignServiceResult struct {
	fx.Out
	Service campaign.ServiceInterface
}

// NewCampaignService creates a new campaign.Service
func NewCampaignService(repo campaign.RepositoryInterface) (CampaignServiceResult, error) {
	service, err := campaign.NewService(repo)
	if err != nil {
		return CampaignServiceResult{}, err
	}
	return CampaignServiceResult{Service: service}, nil
}

// UserServiceResult is the output struct for NewUserService
type UserServiceResult struct {
	fx.Out
	Service user.ServiceInterface
}

// NewUserService creates a new user.Service
func NewUserService(repo user.RepositoryInterface, logger loggo.LoggerInterface) (UserServiceResult, error) {
	service, err := user.NewService(repo.(*user.Repository), logger)
	if err != nil {
		return UserServiceResult{}, err
	}
	return UserServiceResult{Service: service}, nil
}

// EmailServiceResult is the output struct for NewEmailService
type EmailServiceResult struct {
	fx.Out
	Service email.Service
}

// NewEmailService creates a new email.Service
func NewEmailService(cfg *config.Config, logger loggo.LoggerInterface) (EmailServiceResult, error) {
	service, err := email.New(cfg, logger)
	if err != nil {
		return EmailServiceResult{}, err
	}
	return EmailServiceResult{Service: service}, nil
}

func HTTPErrorHandler(err error, c echo.Context, logger loggo.LoggerInterface, cfg *config.Config) {
	message := "Internal Server Error"
	statusCode := http.StatusInternalServerError

	if cfg.IsDevelopment() {
		message = err.Error()
	}

	if httpErr, ok := err.(*echo.HTTPError); ok {
		statusCode = httpErr.Code
		if msg, ok := httpErr.Message.(string); ok {
			message = msg
		}
	}
	errorHandler := shared.NewErrorHandler(logger)
	if err := errorHandler.HandleHTTPError(c, err, message, statusCode); err != nil {
		logger.Error("Failed to handle HTTP error", err)
	}
}
