package config

import (
	"os"
	"testing"

	"github.com/jonesrussell/loggo"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Test case 1: All required environment variables are set
	t.Run("AllVariablesSet", func(t *testing.T) {
		// Clear any existing environment variables
		os.Clearenv()

		// Set up environment variables
		envVars := map[string]string{
			"APP_DEBUG":      "true",
			"APP_ENV":        "test",
			"APP_PORT":       "8080",
			"DB_HOST":        "localhost",
			"DB_NAME":        "testdb",
			"DB_PASS":        "password",
			"DB_PORT":        "5432",
			"DB_USER":        "testuser",
			"SESSION_NAME":   "testsession",
			"SESSION_SECRET": "testsecret",
			"LOG_LEVEL":      "debug",
		}

		for key, value := range envVars {
			os.Setenv(key, value)
		}

		// Defer cleanup of environment variables
		defer os.Clearenv()

		// Load configuration
		config, err := Load()

		// Assert no error
		assert.NoError(t, err)

		// Assert config values
		assert.Equal(t, "true", config.AppDebug)
		assert.Equal(t, "test", config.AppEnv)
		assert.Equal(t, "8080", config.AppPort)
		assert.Equal(t, "localhost", config.DBHost)
		assert.Equal(t, "testdb", config.DBName)
		assert.Equal(t, "password", config.DBPass)
		assert.Equal(t, "5432", config.DBPort)
		assert.Equal(t, "testuser", config.DBUser)
		assert.Equal(t, "testsession", config.SessionName)
		assert.Equal(t, "testsecret", config.SessionSecret)
		assert.Equal(t, "debug", config.LogLevel)
	})

	// Test case 2: Missing required environment variable
	t.Run("MissingRequiredVariable", func(t *testing.T) {
		// Clear all environment variables
		os.Clearenv()

		// Load configuration
		_, err := Load()

		// Assert error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "SESSION_SECRET is not set")
	})

	// Test case 3: Default values
	t.Run("DefaultValues", func(t *testing.T) {
		// Clear all environment variables
		os.Clearenv()

		// Set only required variables
		os.Setenv("SESSION_SECRET", "testsecret")

		// Load configuration
		config, err := Load()

		// Assert no error
		assert.NoError(t, err)

		// Assert default values
		assert.Equal(t, "8080", config.AppPort)
		assert.Equal(t, "info", config.LogLevel)
	})
}

func TestConfig_GetLogLevel(t *testing.T) {
	testCases := []struct {
		name          string
		logLevel      string
		expectedLevel loggo.Level
	}{
		{"Debug", "debug", loggo.LevelDebug},
		{"Info", "info", loggo.LevelInfo},
		{"Warn", "warn", loggo.LevelWarn},
		{"Error", "error", loggo.LevelError},
		{"Default", "invalid", loggo.LevelInfo},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{LogLevel: tc.logLevel}
			assert.Equal(t, tc.expectedLevel, config.GetLogLevel())
		})
	}
}

func TestConfig_DatabaseDSN(t *testing.T) {
	config := &Config{
		DBUser: "testuser",
		DBPass: "testpass",
		DBHost: "localhost",
		DBPort: "5432",
		DBName: "testdb",
	}

	expectedDSN := "testuser:testpass@tcp(localhost:5432)/testdb"
	assert.Equal(t, expectedDSN, config.DatabaseDSN())
}
