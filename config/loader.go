package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// Load loads the configuration from environment variables.
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env file not found, using environment variables")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Normalize paths
	appRoot, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}
	cfg.LogFile = normalizePath(cfg.LogFile, appRoot)
	cfg.MigrationsPath = normalizePath(cfg.MigrationsPath, appRoot)

	printConfig(cfg)

	return cfg, nil
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

		if len(missing) > 0 {
			return fmt.Errorf("Configuration incomplete. Missing: %v\n\nRun:\n%s",
				missing, strings.Join(commands, "\n"))
		}
	}
	return nil
}
