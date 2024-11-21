package database

import (
	"database/sql"
	"fmt"
	"path/filepath"

	// mysql driver is imported here to register its SQL driver init function
	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
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

	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("mysql"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, absPath); err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	return nil
}
