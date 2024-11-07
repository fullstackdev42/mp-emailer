package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DatabaseDSN returns the formatted DSN string for MySQL
func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// IsDevelopment checks if the app is running in a development environment.
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == EnvDevelopment
}

// IsStaging checks if the app is running in a staging environment.
func (c *Config) IsStaging() bool {
	return c.AppEnv == EnvStaging
}

// IsProduction checks if the app is running in a production environment.
func (c *Config) IsProduction() bool {
	return c.AppEnv == EnvProduction
}

// IsTesting checks if the app is running in a testing environment.
func (c *Config) IsTesting() bool {
	return c.AppEnv == EnvTesting
}

// IsDebugEnabled checks if debug logging is enabled.
func (c *Config) IsDebugEnabled() bool {
	return c.AppDebug
}

// ShouldShowDetailedErrors checks if detailed errors should be shown.
func (c *Config) ShouldShowDetailedErrors() bool {
	return c.IsDevelopment() || c.IsStaging()
}

// RequireHTTPS checks if HTTPS is required.
func (c *Config) RequireHTTPS() bool {
	return c.IsProduction()
}

// AllowCORS checks if CORS is allowed.
func (c *Config) AllowCORS() bool {
	return !c.IsProduction()
}

// GetEnvironment returns the current environment.
func (c *Config) GetEnvironment() Environment {
	return c.AppEnv
}

// validatePaths checks if paths are valid and accessible
func validatePaths(c *Config) error {
	// Ensure migrations directory exists and is readable
	migrationsInfo, err := os.Stat(c.MigrationsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("migrations directory does not exist: %s", c.MigrationsPath)
		}
		return fmt.Errorf("error accessing migrations directory '%s': %w", c.MigrationsPath, err)
	}
	if !migrationsInfo.IsDir() {
		return fmt.Errorf("migrations path '%s' is not a directory", c.MigrationsPath)
	}

	// Ensure log file directory exists and is writable
	logDir := filepath.Dir(c.LogFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory '%s': %w", logDir, err)
	}

	// Test if log file is writable
	f, err := os.OpenFile(c.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("log file '%s' is not writable: %w", c.LogFile, err)
	}
	f.Close()

	return nil
}

// GetAbsolutePath returns the absolute path of a given path
func (c *Config) GetAbsolutePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return path // Return original path if conversion fails
	}
	return abs
}

// GetMigrationsPath returns the path to the migrations directory
func (c *Config) GetMigrationsPath() string {
	return c.MigrationsPath
}

// GetLogFilePath returns the path to the log file
func (c *Config) GetLogFilePath() string {
	return c.LogFile
}

// GetJWTExpiryDuration returns the parsed JWT expiry duration
func (c *Config) GetJWTExpiryDuration() (time.Duration, error) {
	return time.ParseDuration(c.JWTExpiry)
}
