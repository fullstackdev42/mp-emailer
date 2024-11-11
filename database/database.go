package database

import (
	"context"

	"github.com/fullstackdev42/mp-emailer/config"
	dbconfig "github.com/fullstackdev42/mp-emailer/database/config"
	"github.com/fullstackdev42/mp-emailer/database/core"
	"github.com/fullstackdev42/mp-emailer/database/decorators"
	"github.com/jonesrussell/loggo"
	"go.uber.org/fx"
)

func ProvideDatabase(lc fx.Lifecycle, cfg *config.Config, logger loggo.LoggerInterface) (core.Interface, error) {
	retryConfig := dbconfig.NewDefaultRetryConfig()
	db, err := dbconfig.ConnectWithRetry(cfg, retryConfig, logger)
	if err != nil {
		return nil, err
	}

	decoratedDB := decorators.NewLoggingDecorator(db, logger)

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			sqlDB, err := decoratedDB.GetSQLDB()
			if err != nil {
				return err
			}
			return sqlDB.Close()
		},
	})

	return decoratedDB, nil
}
