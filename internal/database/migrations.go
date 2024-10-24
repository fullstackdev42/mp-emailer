package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	// Import MySQL driver for database migrations
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jonesrussell/loggo"
)

// Migrator interface for running and closing migrations
type Migrator interface {
	Up() error
	Close() error
}

// MigrationService struct to encapsulate migration logic
type MigrationService struct {
	dsn            string
	migrationsPath string
	logger         loggo.LoggerInterface
}

// Run executes the migrations
func (ms *MigrationService) Run(migrator Migrator) error {
	// Validate DSN
	if ms.dsn == "" {
		ms.logger.Error("DSN is empty", nil)
		return fmt.Errorf("DSN is required")
	}

	// Log the migrations path for debugging
	ms.logger.Debug("Migrations path: " + ms.migrationsPath)

	// Ensure the migrations directory exists
	if _, err := os.Stat(ms.migrationsPath); os.IsNotExist(err) {
		ms.logger.Error("Migrations directory does not exist", err, "path", ms.migrationsPath)
		return fmt.Errorf("migrations directory does not exist: %w", err)
	}

	// Run migrations
	switch err := migrator.Up(); {
	case err != nil && !errors.Is(err, migrate.ErrNoChange):
		ms.logger.Error("Error running migrations", err)
		return fmt.Errorf("error running migrations: %w", err)
	}

	// Always close the migrator
	if err := migrator.Close(); err != nil {
		ms.logger.Error("Error closing migration instance", err)
		return fmt.Errorf("error closing migrations: %w", err)
	}

	ms.logger.Info("Migrations completed successfully")
	return nil
}

// DefaultMigrator wraps the migrate.Migrate struct
type DefaultMigrator struct {
	*migrate.Migrate
}

// Up runs the migrations
func (dm *DefaultMigrator) Up() error {
	return dm.Migrate.Up()
}

// Close closes the migrations
func (dm *DefaultMigrator) Close() error {
	_, err := dm.Migrate.Close()
	return err
}
