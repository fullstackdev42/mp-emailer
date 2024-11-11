package config

import (
	"fmt"
	"path/filepath"
	"time"
)

// DSN returns the database connection string
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
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
