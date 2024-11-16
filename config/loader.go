package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Load loads the configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{}

	// 1. Load base config file (lowest priority)
	if err := loadConfigFile(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	// 2. Override with .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env file not found, using environment variables")
	}

	// 3. Override with environment variables (highest priority)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env vars: %w", err)
	}

	// Handle paths
	if err := cfg.setupPaths(); err != nil {
		return nil, fmt.Errorf("failed to setup paths: %w", err)
	}

	// Validate configuration
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
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

		if len(missing) > 0 {
			return fmt.Errorf("Configuration incomplete. Missing: %v\n\nRun:\n%s",
				missing, strings.Join(commands, "\n"))
		}
	}
	return nil
}

func loadConfigFile(cfg *Config) error {
	// Try loading from config.yaml first
	if err := loadYAMLConfig("config.yaml", cfg); err != nil {
		// Fall back to default config
		if err := loadYAMLConfig("config/config.default.yaml", cfg); err != nil {
			return fmt.Errorf("failed to load default config: %w", err)
		}
	}
	return nil
}

func loadYAMLConfig(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, cfg)
}

func (c *Config) validate() error {
	// Check required fields
	if err := c.validateRequired(); err != nil {
		return err
	}

	// Check environment
	if !c.App.Env.IsValidEnvironment() {
		return fmt.Errorf("invalid environment: %s", c.App.Env)
	}

	// Validate paths are absolute
	if !filepath.IsAbs(c.Log.File) {
		return fmt.Errorf("log file path must be absolute: %s", c.Log.File)
	}
	if !filepath.IsAbs(c.Server.MigrationsPath) {
		return fmt.Errorf("migrations path must be absolute: %s", c.Server.MigrationsPath)
	}

	return nil
}

func (c *Config) validateRequired() error {
	missing := []string{}

	// Database checks
	if c.Database.User == "" {
		missing = append(missing, "DB_USER")
	}
	if c.Database.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if c.Database.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if c.Database.Name == "" {
		missing = append(missing, "DB_NAME")
	}

	// Auth checks
	if c.Auth.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if c.Auth.SessionSecret == "" {
		missing = append(missing, "SESSION_SECRET")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}
	return nil
}
