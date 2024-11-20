package database

import (
	"fmt"

	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/mp-emailer/config"
	"github.com/jonesrussell/mp-emailer/database/core"
	"github.com/jonesrussell/mp-emailer/database/decorators"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Module defines the database module
//
//nolint:gochecknoglobals
var Module = fx.Module("database",
	fx.Provide(
		fx.Annotated{
			Name: "database",
			Target: func(cfg *config.Config, logger loggo.LoggerInterface) (core.Interface, error) {
				db, err := ProvideDatabase(cfg, logger)
				if err != nil {
					logger.Error("Failed to provide database", err)
					return nil, err // Let fx handle the error
				}
				return db, nil
			},
		},
	),
)

// ProvideDatabase creates and configures the database connection with retry logic and logging
func ProvideDatabase(cfg *config.Config, logger loggo.LoggerInterface) (core.Interface, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return decorators.NewLoggingDecorator(&core.DB{GormDB: db}, logger), nil
}
