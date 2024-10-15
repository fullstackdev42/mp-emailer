package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/jonesrussell/loggo"
)

type Config struct {
	AppDebug       string
	AppEnv         string
	AppPort        string
	DBHost         string
	DBName         string
	DBPass         string
	DBPort         string
	DBUser         string
	MailgunAPIKey  string
	MailgunDomain  string
	MailpitHost    string
	MailpitPort    string
	MigrationsPath string
	SessionName    string
	SessionSecret  string
	LogLevel       string
}

func Load() (*Config, error) {
	// Load .env file if it exists, but don't return an error if it doesn't
	_ = godotenv.Load()

	config := &Config{
		AppDebug:       os.Getenv("APP_DEBUG"),
		AppEnv:         os.Getenv("APP_ENV"),
		AppPort:        os.Getenv("APP_PORT"),
		DBHost:         os.Getenv("DB_HOST"),
		DBName:         os.Getenv("DB_NAME"),
		DBPass:         os.Getenv("DB_PASS"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		MailgunAPIKey:  os.Getenv("MAILGUN_API_KEY"),
		MailgunDomain:  os.Getenv("MAILGUN_DOMAIN"),
		MailpitHost:    os.Getenv("MAILPIT_HOST"),
		MailpitPort:    os.Getenv("MAILPIT_PORT"),
		MigrationsPath: os.Getenv("MIGRATIONS_PATH"),
		SessionName:    os.Getenv("SESSION_NAME"),
		SessionSecret:  os.Getenv("SESSION_SECRET"),
		LogLevel:       os.Getenv("LOG_LEVEL"),
	}

	// Set default values or perform validations
	if config.AppDebug == "true" {
		config.LogLevel = "debug"
	}
	if config.SessionSecret == "" {
		return nil, fmt.Errorf("SESSION_SECRET is not set in the environment")
	}
	if config.AppPort == "" {
		config.AppPort = "8080" // Set default port
	}
	if config.LogLevel == "" {
		config.LogLevel = "info" // Set default log level
	}

	return config, nil
}

func (c *Config) GetLogLevel() loggo.Level {
	switch c.LogLevel {
	case "debug":
		return loggo.LevelDebug
	case "info":
		return loggo.LevelInfo
	case "warn":
		return loggo.LevelWarn
	case "error":
		return loggo.LevelError
	default:
		return loggo.LevelInfo
	}
}

// Add this method to the Config struct
func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}
