package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/jonesrussell/loggo"
)

// Config holds the application's configuration values.
type Config struct {
	AppDebug                    bool
	AppEnv                      Environment
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

// Log is used for logging configuration without sensitive fields
type Log struct {
	*Config
	JWTSecret     string `json:"-"`
	MailgunAPIKey string `json:"-"`
	SessionSecret string `json:"-"`
}

// Load loads the configuration from environment variables.
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Get application root directory
	appRoot, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	config := &Config{
		AppDebug:                    getEnv("APP_DEBUG", "false") == "true",
		AppEnv:                      Environment(getEnv("APP_ENV", string(EnvProduction))),
		AppPort:                     getEnv("APP_PORT", "8080"),
		DBHost:                      getEnv("DB_HOST", "localhost"),
		DBName:                      getEnv("DB_NAME", "db"),
		DBPassword:                  getEnv("DB_PASSWORD", "db"),
		DBPort:                      getEnv("DB_PORT", "3306"), // MariaDB default port
		DBUser:                      getEnv("DB_USER", "db"),
		JWTExpiry:                   getEnv("JWT_EXPIRY", "24h"),
		JWTSecret:                   getEnv("JWT_SECRET", ""),
		LogFile:                     normalizePath(getEnv("LOG_FILE", "mp-emailer.log"), appRoot),
		LogLevel:                    getEnv("LOG_LEVEL", "info"),
		MailgunAPIKey:               getEnv("MAILGUN_API_KEY", ""),
		MailgunDomain:               getEnv("MAILGUN_DOMAIN", ""),
		MailpitHost:                 getEnv("MAILPIT_HOST", "localhost"),
		MailpitPort:                 getEnv("MAILPIT_PORT", "1025"),
		MigrationsPath:              normalizePath(getEnv("MIGRATIONS_PATH", "migrations"), appRoot),
		RepresentativeLookupBaseURL: getEnv("REPRESENTATIVE_LOOKUP_BASE_URL", "https://represent.opennorth.ca"),
		SessionName:                 getEnv("SESSION_NAME", "mpe"),
		SessionSecret:               os.Getenv("SESSION_SECRET"),
	}

	// Validate required database configuration
	if config.DBUser == "" {
		return nil, fmt.Errorf("DB_USER is not set in the environment")
	}
	if config.DBName == "" {
		return nil, fmt.Errorf("DB_NAME is not set in the environment")
	}
	if config.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD is not set in the environment")
	}

	// Validate required security configuration
	if config.SessionSecret == "" {
		return nil, fmt.Errorf("SESSION_SECRET is not set in the environment")
	}
	if config.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set in the environment")
	}

	// Validate JWT expiry duration
	if _, err := time.ParseDuration(config.JWTExpiry); err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY duration '%s': must be a valid duration (e.g., '24h', '30m', '15m'): %w",
			config.JWTExpiry, err)
	}

	// Validate numeric ports
	if _, err := strconv.Atoi(config.DBPort); err != nil {
		return nil, fmt.Errorf("invalid DB_PORT '%s': must be a valid number", config.DBPort)
	}
	if _, err := strconv.Atoi(config.MailpitPort); err != nil {
		return nil, fmt.Errorf("invalid MAILPIT_PORT '%s': must be a valid number", config.MailpitPort)
	}
	if _, err := strconv.Atoi(config.AppPort); err != nil {
		return nil, fmt.Errorf("invalid APP_PORT '%s': must be a valid number", config.AppPort)
	}

	// Validate environment value
	if !config.AppEnv.IsValidEnvironment() {
		return nil, fmt.Errorf("invalid APP_ENV '%s': must be one of: development, staging, production, testing",
			config.AppEnv)
	}

	// Validate log level
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[config.LogLevel] {
		return nil, fmt.Errorf("invalid LOG_LEVEL '%s': must be one of: debug, info, warn, error",
			config.LogLevel)
	}

	// Validate file paths
	if err := validatePaths(config); err != nil {
		return nil, err
	}

	// Print the Config struct with sensitive fields masked
	configLog := Log{Config: config}
	configJSON, err := json.MarshalIndent(configLog, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling config to JSON: %v\n", err)
	} else {
		fmt.Printf("Loaded configuration:\n%s\n", string(configJSON))
	}

	return config, nil
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
	return c.AppEnv == EnvDevelopment
}

// IsStaging checks if the app is running in a staging environment.
func (c *Config) IsStaging() bool {
	return c.AppEnv == EnvStaging
}

// IsProduction checks if the app is running in a production environment.
func (c *Config) IsProduction() bool {
	return c.AppEnv == EnvProduction
}

// IsTesting checks if the app is running in a testing environment.
func (c *Config) IsTesting() bool {
	return c.AppEnv == EnvTesting
}

// IsDebugEnabled checks if debug logging is enabled.
func (c *Config) IsDebugEnabled() bool {
	return c.AppDebug
}

// ShouldShowDetailedErrors checks if detailed errors should be shown.
func (c *Config) ShouldShowDetailedErrors() bool {
	return c.IsDevelopment() || c.IsStaging()
}

// RequireHTTPS checks if HTTPS is required.
func (c *Config) RequireHTTPS() bool {
	return c.IsProduction()
}

// AllowCORS checks if CORS is allowed.
func (c *Config) AllowCORS() bool {
	return !c.IsProduction()
}

// GetEnvironment returns the current environment.
func (c *Config) GetEnvironment() Environment {
	return c.AppEnv
}

// normalizePath ensures paths are absolute and cleaned
func normalizePath(path, baseDir string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Join(baseDir, path)
	}
	return filepath.Clean(path)
}

// validatePaths checks if paths are valid and accessible
func validatePaths(c *Config) error {
	// Ensure migrations directory exists and is readable
	migrationsInfo, err := os.Stat(c.MigrationsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("migrations directory does not exist: %s", c.MigrationsPath)
		}
		return fmt.Errorf("error accessing migrations directory '%s': %w", c.MigrationsPath, err)
	}
	if !migrationsInfo.IsDir() {
		return fmt.Errorf("migrations path '%s' is not a directory", c.MigrationsPath)
	}

	// Ensure log file directory exists and is writable
	logDir := filepath.Dir(c.LogFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory '%s': %w", logDir, err)
	}

	// Test if log file is writable
	f, err := os.OpenFile(c.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("log file '%s' is not writable: %w", c.LogFile, err)
	}
	f.Close()

	return nil
}

// Helper method to get the absolute path
func (c *Config) GetAbsolutePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return path // Return original path if conversion fails
	}
	return abs
}

// Add these methods to get normalized paths
func (c *Config) GetMigrationsPath() string {
	return c.MigrationsPath
}

func (c *Config) GetLogFilePath() string {
	return c.LogFile
}

// getEnv gets an environment variable or returns a default value if not set.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper method to get parsed JWT expiry duration
func (c *Config) GetJWTExpiryDuration() (time.Duration, error) {
	return time.ParseDuration(c.JWTExpiry)
}
