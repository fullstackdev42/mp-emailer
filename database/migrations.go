package database

import (
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationConfig struct {
	DSN            string
	MigrationsPath string
	AllowDirty     bool
}

func RunMigrations(cfg MigrationConfig) error {
	absPath, err := filepath.Abs(cfg.MigrationsPath)
	if err != nil {
		return fmt.Errorf("invalid migrations path: %w", err)
	}

	m, err := migrate.New(
		"file://"+absPath,
		"mysql://"+cfg.DSN,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error running migrations: %w", err)
	}

	return nil
}
