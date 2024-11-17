package config

import (
	"fmt"
	"os"

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

	// 4. Warn if using default secrets in production
	if cfg.App.Env == EnvProduction {
		if cfg.Auth.JWTSecret == "dev_jwt_secret_do_not_use_in_production" ||
			cfg.Auth.SessionSecret == "dev_session_secret_do_not_use_in_production" {
			return nil, fmt.Errorf("cannot use default secrets in production environment")
		}
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

func loadConfigFile(cfg *Config) error {
	if err := loadYAMLConfig("config.yaml", cfg); err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
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
	// Check environment
	if !c.App.Env.IsValidEnvironment() {
		return fmt.Errorf("invalid environment: %s", c.App.Env)
	}

	return nil
}
