package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/mp-emailer/config"
	"github.com/jonesrussell/mp-emailer/database/core"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Module defines the database module
//
//nolint:gochecknoglobals
var Module = fx.Options(
	fx.Provide(
		// Provide the core.DB implementation
		fx.Annotate(
			func(cfg *config.Config, logger loggo.LoggerInterface) (*core.DB, error) {
				gormConfig := &gorm.Config{
					PrepareStmt: true,
					Logger:      NewGormLogger(logger),
				}

				db, err := gorm.Open(mysql.Open(cfg.DSN()), gormConfig)
				if err != nil {
					logger.Error("Failed to connect to database", err)
					return nil, fmt.Errorf("failed to connect to database: %w", err)
				}

				sqlDB, err := db.DB()
				if err != nil {
					return nil, fmt.Errorf("failed to get underlying *sql.DB: %w", err)
				}

				// Configure connection pool
				sqlDB.SetMaxIdleConns(10)
				sqlDB.SetMaxOpenConns(100)
				sqlDB.SetConnMaxLifetime(time.Hour)

				return &core.DB{GormDB: db}, nil
			},
			fx.As(new(core.Interface)),
		),
		// Then provide the database interface
		fx.Annotate(
			NewDatabaseService,
			fx.As(new(Interface)),
		),
	),
	fx.Invoke(
		registerDatabaseHooks,
	),
)

// Params for dependency injection
type Params struct {
	fx.In
	Config *config.Config
	Logger loggo.LoggerInterface
	CoreDB core.Interface
}

// NewDatabaseService creates a new database service
func NewDatabaseService(params Params) (Interface, error) {
	database := &Database{
		logger: params.Logger,
		config: params.Config,
		db:     params.CoreDB,
	}

	return database, nil
}

func registerDatabaseHooks(lifecycle fx.Lifecycle, db Interface, logger loggo.LoggerInterface) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Database connection established")
			return nil
		},
		OnStop: func(context.Context) error {
			logger.Info("Closing database connection")
			return db.Close()
		},
	})
}
