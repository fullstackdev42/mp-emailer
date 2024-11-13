package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/jonesrussell/loggo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigPriorityOrder(t *testing.T) {
	// Setup
	envFile := `
APP_PORT=3000
DB_USER=envfile_user
`
	err := os.WriteFile(".env", []byte(envFile), 0644)
	require.NoError(t, err)
	defer os.Remove(".env")

	// Set environment variable (highest priority)
	os.Setenv("APP_PORT", "8080")
	defer os.Unsetenv("APP_PORT")

	cfg, err := config.Load()
	require.NoError(t, err)

	// Environment variable should take precedence over .env file
	assert.Equal(t, 8080, cfg.AppPort)
}

func TestRequiredFieldValidation(t *testing.T) {
	// Clear any existing env vars that might interfere
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("SESSION_SECRET")

	err := config.CheckRequired()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DB_USER")
	assert.Contains(t, err.Error(), "DB_PASSWORD")
}

func TestEnvironmentValidation(t *testing.T) {
	tests := []struct {
		name        string
		env         config.Environment
		shouldBeVal bool
	}{
		{"Valid Development", config.EnvDevelopment, true},
		{"Valid Production", config.EnvProduction, true},
		{"Valid Staging", config.EnvStaging, true},
		{"Valid Testing", config.EnvTesting, true},
		{"Invalid Environment", config.Environment("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.shouldBeVal, tt.env.IsValidEnvironment())
		})
	}
}

func TestLogLevelConversion(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
		want     loggo.Level
	}{
		{"Debug Level", "debug", loggo.LevelDebug},
		{"Info Level", "info", loggo.LevelInfo},
		{"Warn Level", "warn", loggo.LevelWarn},
		{"Error Level", "error", loggo.LevelError},
		{"Default Level", "invalid", loggo.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{LogLevel: tt.logLevel}
			assert.Equal(t, tt.want, cfg.GetLogLevel())
		})
	}
}

func TestDefaultValues(t *testing.T) {
	// Clear environment
	os.Clearenv()

	// Set required fields to pass validation
	os.Setenv("DB_USER", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("JWT_SECRET", "test")
	os.Setenv("SESSION_SECRET", "test")

	cfg, err := config.Load()
	require.NoError(t, err)

	// Get the current working directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Construct the expected log file path dynamically
	expectedLogFilePath := fmt.Sprintf("%s/storage/logs/app.log", cwd)

	// Test default values
	assert.Equal(t, false, cfg.AppDebug)
	assert.Equal(t, "localhost", cfg.AppHost)
	assert.Equal(t, 8080, cfg.AppPort)
	assert.Equal(t, 3306, cfg.DBPort)
	assert.Equal(t, "smtp", string(cfg.EmailProvider))
	assert.Equal(t, "24h", cfg.JWTExpiry)
	assert.Equal(t, expectedLogFilePath, cfg.LogFile) // Use the dynamic path here
	assert.Equal(t, "info", cfg.LogLevel)
}
