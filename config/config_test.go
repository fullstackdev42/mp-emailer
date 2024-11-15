package config_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jonesrussell/loggo"
)

func TestConfigPriorityOrder(t *testing.T) {
	// First, clear any existing environment variables
	os.Clearenv()

	// Setup
	envFile := `
APP_PORT=3000
DB_USER=envfile_user
DB_PASSWORD=test_password
DB_HOST=localhost
DB_NAME=testdb
JWT_SECRET=test_jwt_secret
SESSION_SECRET=test_session_secret
`
	err := os.WriteFile(".env", []byte(envFile), 0644)
	require.NoError(t, err)
	defer os.Remove(".env")

	// Set environment variable (highest priority)
	os.Setenv("APP_PORT", "8080")
	// Don't set DB_USER to test .env file fallback

	// Set other required environment variables
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("JWT_SECRET", "test_jwt_secret")
	os.Setenv("SESSION_SECRET", "test_session_secret")

	defer func() {
		os.Unsetenv("APP_PORT")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("SESSION_SECRET")
	}()

	cfg, err := config.Load()
	require.NoError(t, err)

	// Environment variable should take precedence over .env file
	assert.Equal(t, 8080, cfg.AppPort)
	// Verify env file value is used when no environment variable is set
	assert.Equal(t, "envfile_user", cfg.DBUser)
}

func TestRequiredFieldValidation(t *testing.T) {
	// Clear any existing env vars that might interfere
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("SESSION_SECRET")

	_, err := config.Load()
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
		{"Debug Level", "debug", loggo.LevelDebug},
		{"Info Level", "info", loggo.LevelInfo},
		{"Warn Level", "warn", loggo.LevelWarn},
		{"Error Level", "error", loggo.LevelError},
		{"Default Level", "invalid", loggo.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{LogLevel: tt.logLevel}
			level := cfg.GetLogLevel()
			assert.Equal(t, tt.want, level)
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
	assert.Equal(t, "0.0.0.0", cfg.AppHost)
	assert.Equal(t, 8080, cfg.AppPort)
	assert.Equal(t, 3306, cfg.DBPort)
	assert.Equal(t, "smtp", string(cfg.EmailProvider))
	assert.Equal(t, "24h", cfg.JWTExpiry)
	assert.True(t, strings.HasSuffix(cfg.LogFile, expectedLogFilePath))
	assert.Equal(t, "info", cfg.LogLevel)
}
