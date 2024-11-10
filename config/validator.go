package config

import (
	"fmt"
	"os"
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

	if err := validateEmailConfig(config); err != nil {
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
	if config.DBPort <= 0 {
		return fmt.Errorf("invalid DB_PORT '%d': must be a positive number", config.DBPort)
	}

	// Validate Mailpit Port
	if config.SMTPPort <= 0 {
		return fmt.Errorf("invalid SMTP_PORT '%d': must be a valid number", config.SMTPPort)
	}

	// Validate App Port
	if config.AppPort <= 0 {
		return fmt.Errorf("invalid APP_PORT '%d': must be a positive number", config.AppPort)
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

func validateEmailConfig(config *Config) error {
	switch config.EmailProvider {
	case EmailProviderSMTP:
		if config.SMTPHost == "" || config.SMTPPort <= 0 || config.SMTPFrom == "" {
			return fmt.Errorf("SMTP_HOST, SMTP_PORT, and SMTP_FROM are required when EMAIL_PROVIDER=smtp")
		}
	case EmailProviderMailgun:
		if config.MailgunAPIKey == "" || config.MailgunDomain == "" {
			return fmt.Errorf("MAILGUN_API_KEY and MAILGUN_DOMAIN are required when EMAIL_PROVIDER=mailgun")
		}
	default:
		return fmt.Errorf("invalid EMAIL_PROVIDER '%s': must be one of: smtp, mailgun", config.EmailProvider)
	}
	return nil
}

func validateSMTPConfig(config *Config) error {
	if config.EmailProvider == EmailProviderSMTP {
		if config.SMTPHost == "" {
			return fmt.Errorf("SMTP_HOST is required when EMAIL_PROVIDER is smtp")
		}
		if config.SMTPPort == 0 {
			return fmt.Errorf("SMTP_PORT is required when EMAIL_PROVIDER is smtp")
		}
		if config.SMTPUsername == "" {
			return fmt.Errorf("SMTP_USERNAME is required when EMAIL_PROVIDER is smtp")
		}
		if config.SMTPPassword == "" {
			return fmt.Errorf("SMTP_PASSWORD is required when EMAIL_PROVIDER is smtp")
		}
		if config.SMTPFrom == "" {
			return fmt.Errorf("SMTP_FROM is required when EMAIL_PROVIDER is smtp")
		}

		// Validate port range
		if config.SMTPPort < 1 || config.SMTPPort > 65535 {
			return fmt.Errorf("invalid SMTP_PORT %d: must be between 1 and 65535", config.SMTPPort)
		}
	}
	return nil
}
