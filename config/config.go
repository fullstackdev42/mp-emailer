package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/jonesrussell/loggo"
)

// Config holds the application's configuration values.
type Config struct {
	AppDebug                    string
	AppEnv                      string
	AppPort                     string
	DBHost                      string
	DBName                      string
	DBPassword                  string
	DBPort                      string
	DBUser                      string
	JWTExpiry                   string
	JWTSecret                   string
	LogFile                     string
	LogLevel                    string
	MailgunAPIKey               string
	MailgunDomain               string
	MailpitHost                 string
	MailpitPort                 string
	MigrationsPath              string
	RepresentativeLookupBaseURL string
	SessionName                 string
	SessionSecret               string
}

// Load loads the configuration from environment variables.
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		AppDebug:                    getEnv("APP_DEBUG", "false"),
		AppEnv:                      getEnv("APP_ENV", "production"),
		AppPort:                     getEnv("APP_PORT", "8080"),
		DBHost:                      getEnv("DB_HOST", "localhost"),
		DBName:                      getEnv("DB_NAME", "db"),
		DBPassword:                  getEnv("DB_PASSWORD", "db"),
		DBPort:                      getEnv("DB_PORT", "3306"), // MariaDB default port
		DBUser:                      getEnv("DB_USER", "db"),
		JWTExpiry:                   getEnv("JWT_EXPIRY", "1h"),
		JWTSecret:                   getEnv("JWT_SECRET", ""),
		LogFile:                     getEnv("LOG_FILE", "mp-emailer.log"),
		LogLevel:                    getEnv("LOG_LEVEL", "info"),
		MailgunAPIKey:               getEnv("MAILGUN_API_KEY", ""),
		MailgunDomain:               getEnv("MAILGUN_DOMAIN", ""),
		MailpitHost:                 getEnv("MAILPIT_HOST", "localhost"),
		MailpitPort:                 getEnv("MAILPIT_PORT", "1025"),
		MigrationsPath:              getEnv("MIGRATIONS_PATH", "migrations"),
		RepresentativeLookupBaseURL: getEnv("REPRESENTATIVE_LOOKUP_BASE_URL", "https://represent.opennorth.ca"),
		SessionName:                 getEnv("SESSION_NAME", "mpe"),
		SessionSecret:               os.Getenv("SESSION_SECRET"),
	}

	// Validate required variables
	if config.DBUser == "" || config.DBName == "" || config.DBPassword == "" {
		return nil, fmt.Errorf("DB_USER, DB_NAME, and DB_PASSWORD must be set in the environment")
	}
	if config.SessionSecret == "" {
		return nil, fmt.Errorf("SESSION_SECRET is not set in the environment")
	}

	// Print the Config struct
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling config to JSON: %v\n", err)
	} else {
		fmt.Printf("Loaded configuration:\n%s\n", string(configJSON))
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value if not set.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetLogLevel maps the log level string to loggo.Level.
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

// DatabaseDSN returns the Data Source Name for connecting to the database.
func (c *Config) DatabaseDSN() string {
	// DSN format specific to MariaDB
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// IsDevelopment checks if the app is running in a development environment.
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}
