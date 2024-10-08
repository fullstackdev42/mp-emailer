package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fullstackdev42/mp-emailer/pkg/api"
	"github.com/fullstackdev42/mp-emailer/pkg/database"
	"github.com/fullstackdev42/mp-emailer/pkg/handlers"
	"github.com/fullstackdev42/mp-emailer/pkg/templates"
	"github.com/joho/godotenv"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}

	logger, err := loggo.NewLogger("mp-emailer.log", loggo.LevelInfo)
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		return
	}

	// Create API client
	client := api.NewClient(logger)

	// Initialize database connection
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := database.NewDB(dsn, logger, "./migrations")
	if err != nil {
		logger.Error("Error connecting to database", err)
		return
	}
	defer db.Close()

	// Create a new Echo instance
	e := echo.New()

	// Set renderer
	e.Renderer = templates.NewRenderer()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Get the session secret from environment variables
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		logger.Error("SESSION_SECRET is not set in the environment", nil)
		return
	}

	// Create a new handler with the logger, client, session secret, and database
	h := handlers.NewHandler(logger, client, sessionSecret, db)

	e.Use(h.AuthMiddleware)

	// Routes
	e.GET("/", h.HandleIndex)
	e.GET("/login", h.HandleLogin)
	e.POST("/login", h.HandleLogin)
	e.GET("/logout", h.HandleLogout)

	e.GET("/register", h.HandleRegister)
	e.POST("/register", h.HandleRegister)

	// Protected routes
	e.GET("/submit", h.HandleSubmit, h.AuthMiddleware)
	e.POST("/submit", h.HandleSubmit, h.AuthMiddleware)
	e.POST("/echo", h.HandleEcho, h.AuthMiddleware)

	// Start server
	logger.Info("Starting server on :8080")
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		logger.Error("Error starting server", err)
	}
}
