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
	"github.com/labstack/echo/v4/middleware"
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
			newTemplateManager,
			newSessionStore,
			NewHandler,
			newEcho,
		),
		campaign.ProvideModule(),
		user.ProvideModule(),
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

func newEcho(cfg *config.Config, logger loggo.LoggerInterface, tmplManager *server.TemplateManager, routes []server.Route) *echo.Echo {
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

	// Register the routes
	for _, route := range routes {
		e.Add("GET", route.Pattern(), echo.WrapHandler(route))
	}

	return e
}

func registerRoutes(
	e *echo.Echo,
	serverHandler *server.Handler,
	userHandler *user.Handler,
	campaignHandler *campaign.Handler,
	logger loggo.LoggerInterface,
) {
	logger.Debug("Registering routes")

	// Server routes
	e.GET("/", serverHandler.HandleIndex)

	// User routes
	e.GET("/user/register", userHandler.RegisterGET)
	e.POST("/user/register", userHandler.RegisterPOST)
	e.GET("/user/login", userHandler.LoginGET)
	e.POST("/user/login", userHandler.LoginPOST)
	e.GET("/user/logout", userHandler.LogoutGET)

	// Campaign routes
	e.GET("/campaign", campaignHandler.GetAllCampaigns)
	e.POST("/campaign", campaignHandler.CreateCampaign)
	e.GET("/campaign/:id", campaignHandler.CampaignGET)
	e.PUT("/campaign/:id", campaignHandler.EditCampaign)
	e.DELETE("/campaign/:id", campaignHandler.DeleteCampaign)
	e.POST("/campaign/:id/send", campaignHandler.SendCampaign)

	logger.Debug("Routes registered successfully")
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
	UserHandler     *user.Handler `name:"mainUserHandler"`
	CampaignHandler *campaign.Handler
}

// NewHandler creates a new server.Handler
func NewHandler(
	logger loggo.LoggerInterface,
	emailService email.Service,
	tmplManager *server.TemplateManager,
	userService user.ServiceInterface,
	campaignService campaign.ServiceInterface,
	config *config.Config,
	userRepo user.RepositoryInterface,
	representativeLookupService campaign.RepresentativeLookupServiceInterface,
	campaignClient campaign.ClientInterface,
) (HandlerResult, error) {
	// Create the server handler
	serverHandler := server.NewHandler(
		logger,
		emailService,
		tmplManager,
		userService,
		campaignService,
	)
	// Create the user handler
	userHandler, err := user.NewHandler(userRepo, userService, logger, sessions.NewCookieStore([]byte(config.SessionSecret)), config)
	if err != nil {
		return HandlerResult{}, fmt.Errorf("failed to create user handler: %w", err)
	}
	// Create the campaign handler
	campaignHandler := campaign.NewHandler(
		campaignService,
		logger,
		representativeLookupService,
		emailService,
		campaignClient,
	)

	return HandlerResult{
		ServerHandler:   serverHandler,
		UserHandler:     userHandler.Handler,
		CampaignHandler: campaignHandler,
	}, nil
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
