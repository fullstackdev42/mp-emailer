package database

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	// Import the MySQL driver for database/sql
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	// Import the file source for migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jonesrussell/loggo"
)

func RunMigrations(dsn string, migrationsPath string, logger *loggo.Logger) error {
	// Log the migrations path for debugging
	logger.Info("Migrations path: " + migrationsPath)

	// Ensure the migrations directory exists
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", migrationsPath)
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("mysql://%s", dsn),
	)
	if err != nil {
		return fmt.Errorf("error creating migration instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error running migrations: %w", err)
	}

	logger.Info("Migrations completed successfully")
	return nil
}
