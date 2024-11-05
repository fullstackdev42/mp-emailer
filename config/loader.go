package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Load loads the configuration from environment variables.
func Load() (*Config, error) {
	_ = godotenv.Load()

	appRoot, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	config := loadFromEnv(appRoot)

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	printConfig(config)

	return config, nil
}

func loadFromEnv(appRoot string) *Config {
	return &Config{
		AppDebug:                    getEnv("APP_DEBUG", "false") == "true",
		AppEnv:                      Environment(getEnv("APP_ENV", string(EnvProduction))),
		AppPort:                     getEnv("APP_PORT", "8080"),
		DBHost:                      getEnv("DB_HOST", "localhost"),
		DBName:                      getEnv("DB_NAME", "db"),
		DBPassword:                  getEnv("DB_PASSWORD", "db"),
		DBPort:                      getEnv("DB_PORT", "3306"),
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
}

func printConfig(config *Config) {
	// Create a Log struct that masks sensitive fields
	configLog := Log{Config: config}

	fmt.Printf("Loaded configuration: %+v\n", configLog)
}
