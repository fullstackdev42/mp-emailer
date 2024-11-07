package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Load loads the configuration from environment variables.
func Load() (*Config, error) {
	// Load .env file if it exists, ignore error if it doesn't
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env file not found, using environment variables")
	}

	// Check for required SESSION_SECRET early
	if os.Getenv("SESSION_SECRET") == "" {
		return nil, fmt.Errorf("SESSION_SECRET not set.\nRun: \nRun: echo 'SESSION_SECRET='$(openssl rand -base64 32) >> .env")
	}

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
		MigrationsPath:              normalizePath(getEnv("MIGRATIONS_PATH", "migrations"), appRoot),
		RepresentativeLookupBaseURL: getEnv("REPRESENTATIVE_LOOKUP_BASE_URL", "https://represent.opennorth.ca"),
		SessionName:                 getEnv("SESSION_NAME", "mpe"),
		SessionSecret:               os.Getenv("SESSION_SECRET"),
		EmailProvider:               EmailProvider(getEnv("EMAIL_PROVIDER", string(EmailProviderSMTP))),
		SMTPHost:                    getEnv("SMTP_HOST", "mailpit"),
		SMTPPort:                    getEnv("SMTP_PORT", "1025"),
		SMTPUsername:                getEnv("SMTP_USERNAME", ""),
		SMTPPassword:                getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:                    getEnv("SMTP_FROM", "test@example.com"),
	}
}

func printConfig(config *Config) {
	// Create a Log struct that masks sensitive fields
	configLog := Log{Config: config}

	fmt.Printf("Loaded configuration: %+v\n", configLog)
}

// CheckRequired verifies all required configuration is present before loading full config
func CheckRequired() error {
	if err := godotenv.Load(); err != nil {
		// Create empty config to get required vars
		cfg := &Config{}
		missing := []string{}
		commands := []string{}

		for envVar, command := range cfg.RequiredEnvVars() {
			if os.Getenv(envVar) == "" {
				missing = append(missing, envVar)
				commands = append(commands, command)
			}
		}

		return fmt.Errorf("Configuration incomplete. Missing: %v\n\nRun:\n%s",
			missing, strings.Join(commands, "\n"))
	}
	return nil
}
