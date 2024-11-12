package config

import (
	"os"
	"testing"

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

	cfg := &Config{}
	err = LoadConfig(cfg)
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

	err := CheckRequired()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DB_USER")
	assert.Contains(t, err.Error(), "DB_PASSWORD")
}

func TestEnvironmentValidation(t *testing.T) {
	tests := []struct {
		name        string
		env         Environment
		shouldBeVal bool
	}{
		{"Valid Development", EnvDevelopment, true},
		{"Valid Production", EnvProduction, true},
		{"Valid Staging", EnvStaging, true},
		{"Valid Testing", EnvTesting, true},
		{"Invalid Environment", Environment("invalid"), false},
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
		want     string
	}{
		{"Debug Level", "debug", "debug"},
		{"Info Level", "info", "info"},
		{"Warn Level", "warn", "warn"},
		{"Error Level", "error", "error"},
		{"Default Level", "invalid", "info"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{LogLevel: tt.logLevel}
			level := cfg.GetLogLevel()
			assert.Equal(t, tt.want, level.String())
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

	cfg := &Config{}
	err := LoadConfig(cfg)
	require.NoError(t, err)

	// Test default values
	assert.Equal(t, false, cfg.AppDebug)
	assert.Equal(t, "localhost", cfg.AppHost)
	assert.Equal(t, 8080, cfg.AppPort)
	assert.Equal(t, 3306, cfg.DBPort)
	assert.Equal(t, "smtp", string(cfg.EmailProvider))
	assert.Equal(t, "24h", cfg.JWTExpiry)
	assert.Equal(t, "storage/logs/app.log", cfg.LogFile)
	assert.Equal(t, "info", cfg.LogLevel)
}
