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

// Load loads the configuration
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		AppDebug:       getEnv("APP_DEBUG", "false"),
		AppEnv:         getEnv("APP_ENV", "development"),
		AppPort:        getEnv("APP_PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBName:         os.Getenv("DB_NAME"),
		DBPass:         os.Getenv("DB_PASS"),
		DBPort:         getEnv("DB_PORT", "3306"), // MariaDB default port
		DBUser:         os.Getenv("DB_USER"),
		MailgunAPIKey:  getEnv("MAILGUN_API_KEY", ""),
		MailgunDomain:  getEnv("MAILGUN_DOMAIN", ""),
		MailpitHost:    getEnv("MAILPIT_HOST", "localhost"),
		MailpitPort:    getEnv("MAILPIT_PORT", "1025"),
		MigrationsPath: getEnv("MIGRATIONS_PATH", "migrations"),
		SessionName:    getEnv("SESSION_NAME", "mpe"),
		SessionSecret:  os.Getenv("SESSION_SECRET"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}

	// Validate required variables
	if config.DBUser == "" || config.DBName == "" || config.DBPass == "" {
		return nil, fmt.Errorf("DB_USER, DB_NAME, and DB_PASS must be set in the environment")
	}
	if config.SessionSecret == "" {
		return nil, fmt.Errorf("SESSION_SECRET is not set in the environment")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
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

func (c *Config) DatabaseDSN() string {
	// DSN format specific to MariaDB
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}
