package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/fullstackdev42/mp-emailer/pkg/database"
	"github.com/fullstackdev42/mp-emailer/pkg/handlers"
	"github.com/fullstackdev42/mp-emailer/pkg/services"
	"github.com/joho/godotenv"
	"github.com/jonesrussell/loggo"
	"golang.org/x/crypto/bcrypt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading .env file: %v", err)
	}

	logger, err := loggo.NewLogger("register_user.log", loggo.LevelInfo)
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		return
	}

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

	// Log email service configuration
	mailpitHost := os.Getenv("MAILPIT_HOST")
	mailpitPort := os.Getenv("MAILPIT_PORT")
	logger.Info("Email service configuration:",
		"host", mailpitHost,
		"port", mailpitPort)

	// Check if Mailpit environment variables are set
	if mailpitHost == "" || mailpitPort == "" {
		logger.Error("Mailpit host or port environment variables are not set", err)
		return
	}

	emailService := services.NewMailpitEmailService(mailpitHost, mailpitPort)

	username := randString(8)
	email := fmt.Sprintf("%s@example.com", randString(8))
	password := randString(12)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Error hashing password", err)
		return
	}

	err = db.CreateUser(username, email, string(hashedPassword))
	if err != nil {
		logger.Error("Error creating user", err)
		return
	}
	handler := handlers.NewHandler(logger, nil, "", db, emailService, nil)
	if err := handler.SendAdminNotification(username, email); err != nil {
		logger.Error("Failed to send admin notification email", err)
		logger.Info("Continuing execution despite email sending failure")
	} else {
		logger.Info("Admin notification email sent successfully")
	}

	fmt.Printf("User registered successfully:\nUsername: %s\nEmail: %s\nPassword: %s\n", username, email, password)
}
