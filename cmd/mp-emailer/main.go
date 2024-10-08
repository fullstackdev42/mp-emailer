package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/fullstackdev42/mp-emailer/pkg/api"
	"github.com/fullstackdev42/mp-emailer/pkg/database"
	"github.com/fullstackdev42/mp-emailer/pkg/handlers"
	"github.com/fullstackdev42/mp-emailer/pkg/services"
	"github.com/fullstackdev42/mp-emailer/pkg/templates"
	"github.com/joho/godotenv"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	AppEnv        string
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
	AppPort       string
}

func loadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	config := &Config{
		AppEnv:        os.Getenv("APP_ENV"),
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
		AppPort:       os.Getenv("APP_PORT"),
	}

	if config.SessionSecret == "" {
		return nil, fmt.Errorf("SESSION_SECRET is not set in the environment")
	}

	if config.AppPort == "" {
		config.AppPort = "8080"
	}

	return config, nil
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

	var emailService services.EmailService
	if config.AppEnv == "production" {
		emailService = services.NewMailgunEmailService(config.MailgunDomain, config.MailgunAPIKey)
	} else {
		emailService = services.NewMailpitEmailService(config.MailpitHost, config.MailpitPort)
	}

	client := api.NewClient(logger)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName)
	db, err := database.NewDB(dsn, logger, "./migrations")
	if err != nil {
		logger.Error("Error connecting to database", err)
		return
	}
	defer db.Close()

	e := echo.New()
	e.Renderer = templates.NewRenderer()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := handlers.NewHandler(logger, client, config.SessionSecret, db, emailService)
	e.Use(h.AuthMiddleware)

	e.GET("/", h.HandleIndex)
	e.GET("/login", h.HandleLogin)
	e.POST("/login", h.HandleLogin)
	e.GET("/logout", h.HandleLogout)
	e.GET("/register", h.HandleRegister)
	e.POST("/register", h.HandleRegister)
	e.GET("/submit", h.HandleSubmit, h.AuthMiddleware)
	e.POST("/submit", h.HandleSubmit, h.AuthMiddleware)
	e.POST("/echo", h.HandleEcho, h.AuthMiddleware)

	// Campaign routes
	e.GET("/campaigns", h.HandleGetCampaigns, h.AuthMiddleware)
	e.GET("/campaigns/new", h.HandleCreateCampaign, h.AuthMiddleware)
	e.POST("/campaigns/new", h.HandleCreateCampaign, h.AuthMiddleware)
	e.POST("/campaigns/:id/update", h.HandleUpdateCampaign, h.AuthMiddleware)
	e.POST("/campaigns/:id/delete", h.HandleDeleteCampaign, h.AuthMiddleware)

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
