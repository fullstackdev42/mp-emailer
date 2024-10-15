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
			"APP_DEBUG":       "true",
			"APP_ENV":         "test",
			"APP_PORT":        "9090",
			"DB_HOST":         "testhost",
			"DB_NAME":         "testdb",
			"DB_PASS":         "testpass",
			"DB_PORT":         "3307",
			"DB_USER":         "testuser",
			"MAILGUN_API_KEY": "testkey",
			"MAILGUN_DOMAIN":  "testdomain",
			"MAILPIT_HOST":    "testmailpit",
			"MAILPIT_PORT":    "2025",
			"MIGRATIONS_PATH": "testmigrations",
			"SESSION_NAME":    "testsession",
			"SESSION_SECRET":  "testsecret",
			"LOG_LEVEL":       "debug",
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
		assert.Equal(t, "9090", config.AppPort)
		assert.Equal(t, "testhost", config.DBHost)
		assert.Equal(t, "testdb", config.DBName)
		assert.Equal(t, "testpass", config.DBPass)
		assert.Equal(t, "3307", config.DBPort)
		assert.Equal(t, "testuser", config.DBUser)
		assert.Equal(t, "testkey", config.MailgunAPIKey)
		assert.Equal(t, "testdomain", config.MailgunDomain)
		assert.Equal(t, "testmailpit", config.MailpitHost)
		assert.Equal(t, "2025", config.MailpitPort)
		assert.Equal(t, "testmigrations", config.MigrationsPath)
		assert.Equal(t, "testsession", config.SessionName)
		assert.Equal(t, "testsecret", config.SessionSecret)
		assert.Equal(t, "debug", config.LogLevel)
	})

	// Test case 2: Missing required environment variables
	t.Run("MissingRequiredVariables", func(t *testing.T) {
		// Clear all environment variables
		os.Clearenv()

		// Load configuration
		_, err := Load()

		// Assert error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "DB_USER, DB_NAME, and DB_PASS must be set in the environment")
	})

	// Test case 3: Missing SESSION_SECRET
	t.Run("MissingSessionSecret", func(t *testing.T) {
		// Clear all environment variables
		os.Clearenv()

		// Set required DB variables
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("DB_PASS", "testpass")

		// Load configuration
		_, err := Load()

		// Assert error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "SESSION_SECRET is not set in the environment")
	})

	// Test case 4: Default values
	t.Run("DefaultValues", func(t *testing.T) {
		// Clear all environment variables
		os.Clearenv()

		// Set only required variables
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("DB_PASS", "testpass")
		os.Setenv("SESSION_SECRET", "testsecret")

		// Load configuration
		config, err := Load()

		// Assert no error
		assert.NoError(t, err)

		// Assert default values
		assert.Equal(t, "false", config.AppDebug)
		assert.Equal(t, "development", config.AppEnv)
		assert.Equal(t, "8080", config.AppPort)
		assert.Equal(t, "localhost", config.DBHost)
		assert.Equal(t, "3306", config.DBPort)
		assert.Equal(t, "", config.MailgunAPIKey)
		assert.Equal(t, "", config.MailgunDomain)
		assert.Equal(t, "localhost", config.MailpitHost)
		assert.Equal(t, "1025", config.MailpitPort)
		assert.Equal(t, "migrations", config.MigrationsPath)
		assert.Equal(t, "session", config.SessionName)
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
		DBPort: "3306",
		DBName: "testdb",
	}

	expectedDSN := "testuser:testpass@tcp(localhost:3306)/testdb?parseTime=true"
	assert.Equal(t, expectedDSN, config.DatabaseDSN())
}
