package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/fx"

	// Import MySQL driver for database migrations
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrationModule is the module for database migrations
// nolint:gochecknoglobals
var MigrationModule = fx.Options(
	fx.Invoke(RunMigrations),
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
}

// NewMigrationService creates a new instance of MigrationService
func NewMigrationService(config *config.Config, migrationsPath string) *MigrationService {
	return &MigrationService{
		dsn:            config.DatabaseDSN(),
		migrationsPath: migrationsPath,
	}
}

// Run executes the migrations
func (ms *MigrationService) Run(migrator Migrator) error {
	// Validate DSN
	if ms.dsn == "" {
		return fmt.Errorf("DSN is required")
	}

	// Ensure the migrations directory exists
	if _, err := os.Stat(ms.migrationsPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %w", err)
	}

	// Run migrations
	switch err := migrator.Up(); {
	case err != nil && !errors.Is(err, migrate.ErrNoChange):
		return fmt.Errorf("error running migrations: %w", err)
	}

	// Always close the migrator
	if err := migrator.Close(); err != nil {
		return fmt.Errorf("error closing migrations: %w", err)
	}

	return nil
}

// RunMigrations executes the database migrations
func (ms *MigrationService) RunMigrations() error {
	// Create a new migrator instance
	m, err := migrate.New(
		ms.migrationsPath,
		ms.dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Create a default migrator
	migrator := &DefaultMigrator{Migrate: m}

	// Run the migrations
	return ms.Run(migrator)
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

type MigrationParams struct {
	fx.In
	Config *config.Config
}

func RunMigrations(p MigrationParams) error {
	migrationService := NewMigrationService(p.Config, "file://migrations")
	return migrationService.RunMigrations()
}
