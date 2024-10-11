package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"embed"

	"github.com/fullstackdev42/mp-emailer/pkg/api"
	"github.com/fullstackdev42/mp-emailer/pkg/database"
	"github.com/fullstackdev42/mp-emailer/pkg/handlers"
	appmid "github.com/fullstackdev42/mp-emailer/pkg/middleware"
	"github.com/fullstackdev42/mp-emailer/pkg/services"
	"github.com/fullstackdev42/mp-emailer/pkg/templates"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	AppEnv        string
	AppPort       string
	MailgunDomain string
	MailgunAPIKey string
	MailpitHost   string
	MailpitPort   string
	DBUser        string
	DBPass        string
	DBName        string
	DBHost        string
	DBPort        string
	SessionSecret string
}

func loadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	config := &Config{
		AppEnv:        os.Getenv("APP_ENV"),
		AppPort:       os.Getenv("APP_PORT"),
		MailgunDomain: os.Getenv("MAILGUN_DOMAIN"),
		MailgunAPIKey: os.Getenv("MAILGUN_API_KEY"),
		MailpitHost:   os.Getenv("MAILPIT_HOST"),
		MailpitPort:   os.Getenv("MAILPIT_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPass:        os.Getenv("DB_PASS"),
		DBName:        os.Getenv("DB_NAME"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		SessionSecret: os.Getenv("SESSION_SECRET"),
	}

	if config.SessionSecret == "" {
		return nil, fmt.Errorf("SESSION_SECRET is not set in the environment")
	}

	if config.AppPort == "" {
		config.AppPort = "8080"
	}

	return config, nil
}

//go:embed web/public/* web/public/partials/*
var templateFS embed.FS

func dbMiddleware(db *database.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	}
}

func initializeEmailService(config *Config) services.EmailService {
	if config.AppEnv == "production" {
		return services.NewMailgunEmailService(config.MailgunDomain, config.MailgunAPIKey)
	}
	return services.NewMailpitEmailService(config.MailpitHost, config.MailpitPort)
}

func main() {
	logger, err := loggo.NewLogger("mp-emailer.log", loggo.LevelInfo)
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		return
	}

	config, err := loadConfig()
	if err != nil {
		logger.Error("Error loading configuration", err)
		return
	}

	emailService := initializeEmailService(config)

	client := api.NewClient(logger)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName)
	db, err := database.NewDB(dsn, logger, "./migrations")
	if err != nil {
		logger.Error("Error connecting to database", err)
		return
	}
	defer db.Close()

	e := echo.New()
	e.Static("/static", "web/public")

	tmplManager, err := templates.NewTemplateManager(templateFS)
	if err != nil {
		logger.Error("Error initializing templates", err)
		return
	}

	e.Renderer = echo.Renderer(tmplManager)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize session store
	store := sessions.NewCookieStore([]byte(config.SessionSecret))
	e.Use(session.Middleware(store))

	// Apply SetAuthStatusMiddleware after session middleware
	e.Use(appmid.SetAuthStatusMiddleware(store, logger))

	// Add database middleware
	e.Use(dbMiddleware(db))

	h := handlers.NewHandler(logger, client, store, emailService, tmplManager)

	// Public routes
	e.GET("/", h.HandleIndex)
	e.GET("/login", h.HandleLogin)
	e.POST("/login", h.HandleLogin)
	e.GET("/logout", h.HandleLogout)
	e.GET("/register", h.HandleRegister)
	e.POST("/register", h.HandleRegister)

	// Protected routes
	authGroup := e.Group("")
	authGroup.Use(appmid.RequireAuthMiddleware(store, logger))
	authGroup.GET("/submit", h.HandleSubmit)
	authGroup.POST("/submit", h.HandleSubmit)
	authGroup.POST("/echo", h.HandleEcho)

	// Campaign routes (protected)
	authGroup.GET("/campaigns", h.HandleGetCampaigns)
	authGroup.GET("/campaigns/new", h.HandleCreateCampaign)
	authGroup.POST("/campaigns/new", h.HandleCreateCampaign)
	authGroup.GET("/campaigns/:id", h.HandleGetCampaign)
	authGroup.POST("/campaigns/:id/update", h.HandleUpdateCampaign)
	authGroup.POST("/campaigns/:id/delete", h.HandleDeleteCampaign)

	port := config.AppPort
	if _, err := strconv.Atoi(port); err != nil {
		logger.Error("Invalid APP_PORT value", err)
		return
	}

	logger.Info(fmt.Sprintf("Attempting to start server on :%s", port))
	if err := e.Start(":" + port); err != http.ErrServerClosed {
		logger.Error("Error starting server", err)
	}
}
