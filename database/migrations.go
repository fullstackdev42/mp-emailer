package database

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/fx"

	// Import MySQL driver for database migrations
	_ "github.com/go-sql-driver/mysql"
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
	// Add debug logging
	pwd, _ := os.Getwd()
	fmt.Printf("Current working directory: %s\n", pwd)
	fmt.Printf("Looking for migrations in: %s\n", ms.migrationsPath)

	// Validate DSN
	if ms.dsn == "" {
		return fmt.Errorf("DSN is required")
	}

	// Ensure the migrations directory exists
	if _, err := os.Stat(strings.TrimPrefix(ms.migrationsPath, "file://")); os.IsNotExist(err) {
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
	fmt.Println("Starting RunMigrations...")

	sourceURL := fmt.Sprintf("file://%s", strings.TrimPrefix(ms.migrationsPath, "file://"))
	dbURL := fmt.Sprintf("mysql://%s&multiStatements=true", ms.dsn)

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Use a more robust defer for cleanup
	defer func() {
		if m != nil {
			if _, dbErr := m.Close(); dbErr != nil {
				fmt.Printf("Error closing migrations: %v\n", dbErr)
			}
		}
	}()

	// Run migrations with proper error handling
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error running migrations: %w", err)
	}

	fmt.Println("Migrations completed successfully.")
	return nil
}

// DefaultMigrator wraps the migrate.Migrate struct
type DefaultMigrator struct {
	*migrate.Migrate
}

// Up runs the migrations
func (dm *DefaultMigrator) Up() error {
	fmt.Println("Starting migration process...")
	err := dm.Migrate.Up()
	if err != nil {
		fmt.Printf("Error during migration: %v\n", err)
	} else {
		fmt.Println("Migration completed successfully.")
	}
	return err
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
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Construct and verify migrations path
	migrationsDir := filepath.Clean(filepath.Join(pwd, "database", "migrations"))
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory not found at %s: %w", migrationsDir, err)
	}

	// List migration files for debugging
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}
	fmt.Printf("Found migration files in %s:\n", migrationsDir)
	for _, file := range files {
		fmt.Printf("- %s\n", file.Name())
	}

	// Ensure the path is in the correct format for the migrate library
	sourceURL := "file://" + filepath.ToSlash(migrationsDir)
	fmt.Printf("Using migrations source URL: %s\n", sourceURL)

	migrationService := NewMigrationService(p.Config, sourceURL)
	if err := migrationService.RunMigrations(); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
