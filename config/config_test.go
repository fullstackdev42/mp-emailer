package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jonesrussell/mp-emailer/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// setupTestDirectory sets up a temporary directory for testing and clears environment variables.
func setupTestDirectory(t *testing.T) (string, string, func()) {
	t.Helper()

	// Clear any existing environment variables
	os.Clearenv()

	// Create temporary test directory
	tmpDir := t.TempDir()

	// Save current working directory
	originalWd, err := os.Getwd()
	require.NoError(t, err, "could not get current directory")

	// Change to temp directory for test
	err = os.Chdir(tmpDir)
	require.NoError(t, err, "could not change to temp directory")

	// Return cleanup function to restore original state
	cleanup := func() {
		err := os.Chdir(originalWd)
		require.NoError(t, err, "could not restore original directory")
	}

	return tmpDir, originalWd, cleanup
}

func TestConfigPriorityOrder(t *testing.T) {
	_, _, cleanup := setupTestDirectory(t)
	defer cleanup()

	// Create config file
	defaultConfig := `
app:
  port: 3000
database:
  user: default_user
`
	err := os.WriteFile("config.yaml", []byte(defaultConfig), 0644)
	require.NoError(t, err, "could not write config.yaml")

	// Setup .env file
	envFile := `
APP_PORT=3000
DB_USER=envfile_user
DB_PASSWORD=test_password
DB_HOST=localhost
DB_NAME=testdb
JWT_SECRET=test_jwt_secret
SESSION_SECRET=test_session_secret
`
	err = os.WriteFile(".env", []byte(envFile), 0644)
	require.NoError(t, err, "could not write .env file")
	defer os.Remove(".env")

	// Set environment variables (highest priority)
	os.Setenv("APP_PORT", "8080")
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("JWT_SECRET", "test_jwt_secret")
	os.Setenv("SESSION_SECRET", "test_session_secret")

	cfg, err := config.Load()
	require.NoError(t, err, "could not load config")

	// Environment variable should take precedence over .env file
	assert.Equal(t, 8080, cfg.App.Port)
	// Verify .env file value is used when no environment variable is set
	assert.Equal(t, "envfile_user", cfg.Database.User)
}

func TestRequiredFieldValidation(t *testing.T) {
	_, _, cleanup := setupTestDirectory(t)
	defer cleanup()

	// Create minimal config file
	defaultConfig := `
app:
  port: 3000
`
	err := os.WriteFile("config.yaml", []byte(defaultConfig), 0644)
	require.NoError(t, err, "could not write config.yaml")

	// Clear any existing environment variables that might interfere
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("SESSION_SECRET")

	_, err = config.Load()
	assert.Error(t, err, "config.Load() should have failed due to missing required fields")
	assert.Contains(t, err.Error(), "DB_USER", "error should mention missing DB_USER")
	assert.Contains(t, err.Error(), "DB_PASSWORD", "error should mention missing DB_PASSWORD")
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
			assert.Equal(t, tt.shouldBeVal, tt.env.IsValidEnvironment(), "environment validity check failed")
		})
	}
}

func TestLogLevelConversion(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
		want     zapcore.Level
	}{
		{"Debug Level", "debug", zapcore.DebugLevel},
		{"Info Level", "info", zapcore.InfoLevel},
		{"Warn Level", "warn", zapcore.WarnLevel},
		{"Error Level", "error", zapcore.ErrorLevel},
		{"Default Level", "invalid", zapcore.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{Log: config.LogConfig{Level: tt.logLevel}}
			level := cfg.GetLogLevel()
			assert.Equal(t, tt.want, level, "log level conversion failed")
		})
	}
}

func TestDefaultValues(t *testing.T) {
	_, _, cleanup := setupTestDirectory(t)
	defer cleanup()

	// Create minimal config file
	defaultConfig := `
app:
  host: 0.0.0.0
  port: 8080
log:
  file: "storage/logs/app.log"
  format: "json"
`
	err := os.WriteFile("config.yaml", []byte(defaultConfig), 0644)
	require.NoError(t, err, "could not write config.yaml")

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
	require.NoError(t, err, "could not load config")

	// Get the current working directory
	cwd, err := os.Getwd()
	require.NoError(t, err, "could not get current working directory")

	// Construct the expected log file path dynamically
	expectedLogFilePath := fmt.Sprintf("%s/storage/logs/app.log", cwd)

	// Test default values
	assert.Equal(t, false, cfg.App.Debug, "unexpected value for App.Debug")
	assert.Equal(t, "0.0.0.0", cfg.App.Host, "unexpected value for App.Host")
	assert.Equal(t, 8080, cfg.App.Port, "unexpected value for App.Port")
	assert.Equal(t, 3306, cfg.Database.Port, "unexpected value for Database.Port")
	assert.Equal(t, config.EmailProviderSMTP, cfg.Email.Provider, "unexpected value for Email.Provider")
	assert.Equal(t, "24h", cfg.Auth.JWTExpiry, "unexpected value for Auth.JWTExpiry")
	assert.Equal(t, expectedLogFilePath, cfg.Log.File, "unexpected value for Log.File")
	assert.Equal(t, "json", cfg.Log.Format, "unexpected value for Log.Format")
	assert.Equal(t, "info", cfg.Log.Level, "unexpected value for Log.Level")
}

func TestFeatureFlagConfiguration(t *testing.T) {
	_, _, cleanup := setupTestDirectory(t)
	defer cleanup()

	// Create minimal config file
	defaultConfig := `
featureFlags:
  enableSMTP: true
  enableMailgun: false
`
	err := os.WriteFile("config.yaml", []byte(defaultConfig), 0644)
	require.NoError(t, err, "could not write config.yaml")

	// Set required environment variables to pass validation
	os.Setenv("DB_USER", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("JWT_SECRET", "test")
	os.Setenv("SESSION_SECRET", "test")

	// Test default values
	cfg, err := config.Load()
	require.NoError(t, err, "could not load config")

	assert.False(t, cfg.FeatureFlags.EnableMailgun, "expected EnableMailgun to be false")
	assert.True(t, cfg.FeatureFlags.EnableSMTP, "expected EnableSMTP to be true")

	// Test environment override
	os.Setenv("FEATURE_MAILGUN", "true")
	cfg, err = config.Load()
	require.NoError(t, err, "could not load config with environment override")
	assert.True(t, cfg.FeatureFlags.EnableMailgun, "expected EnableMailgun to be true after environment override")
}

func TestVersionConfiguration(t *testing.T) {
	_, _, cleanup := setupTestDirectory(t)
	defer cleanup()

	// Create minimal config file
	defaultConfig := `
app:
  port: 3000
`
	err := os.WriteFile("config.yaml", []byte(defaultConfig), 0644)
	require.NoError(t, err, "could not write config.yaml")

	os.Setenv("APP_VERSION", "1.0.0")
	os.Setenv("BUILD_DATE", "2024-03-21T12:00:00Z")
	os.Setenv("GIT_COMMIT", "abc123")

	// Set required fields
	os.Setenv("DB_USER", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("JWT_SECRET", "test")
	os.Setenv("SESSION_SECRET", "test")

	cfg, err := config.Load()
	require.NoError(t, err)

	status := cfg.GetStatus()
	version := status["version"].(map[string]string)

	assert.Equal(t, "1.0.0", version["version"])
	assert.Equal(t, "2024-03-21T12:00:00Z", version["buildDate"])
	assert.Equal(t, "abc123", version["commit"])
}
