package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func validateConfig(config *Config) error {
	if err := validateRequiredFields(config); err != nil {
		return err
	}

	if err := validateNumericFields(config); err != nil {
		return err
	}

	if err := validateJWTExpiry(config); err != nil {
		return err
	}

	if err := validateEnvironment(config); err != nil {
		return err
	}

	if err := validateLogLevel(config); err != nil {
		return err
	}

	return validatePaths(config)
}

func validateRequiredFields(config *Config) error {
	missing := []string{}
	commands := []string{}

	for envVar, command := range config.RequiredEnvVars() {
		if os.Getenv(envVar) == "" {
			missing = append(missing, envVar)
			commands = append(commands, command)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("Missing required environment variables: %v\nRun the following commands:\n%s",
			missing, strings.Join(commands, "\n"))
	}

	return nil
}

func validateNumericFields(config *Config) error {
	// Validate DB Port
	if _, err := strconv.Atoi(config.DBPort); err != nil {
		return fmt.Errorf("invalid DB_PORT '%s': must be a valid number", config.DBPort)
	}

	// Validate Mailpit Port
	if _, err := strconv.Atoi(config.MailpitPort); err != nil {
		return fmt.Errorf("invalid MAILPIT_PORT '%s': must be a valid number", config.MailpitPort)
	}

	// Validate App Port
	if _, err := strconv.Atoi(config.AppPort); err != nil {
		return fmt.Errorf("invalid APP_PORT '%s': must be a valid number", config.AppPort)
	}

	return nil
}

func validateJWTExpiry(config *Config) error {
	if _, err := time.ParseDuration(config.JWTExpiry); err != nil {
		return fmt.Errorf("invalid JWT_EXPIRY duration '%s': must be a valid duration (e.g., '24h', '30m', '15m'): %w",
			config.JWTExpiry, err)
	}
	return nil
}

func validateEnvironment(config *Config) error {
	if !config.AppEnv.IsValidEnvironment() {
		return fmt.Errorf("invalid APP_ENV '%s': must be one of: development, staging, production, testing",
			config.AppEnv)
	}
	return nil
}

func validateLogLevel(config *Config) error {
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[config.LogLevel] {
		return fmt.Errorf("invalid LOG_LEVEL '%s': must be one of: debug, info, warn, error",
			config.LogLevel)
	}
	return nil
}
