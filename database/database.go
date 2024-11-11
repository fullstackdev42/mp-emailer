package database

import (
	"github.com/fullstackdev42/mp-emailer/config"
	dbconfig "github.com/fullstackdev42/mp-emailer/database/config"
	"github.com/fullstackdev42/mp-emailer/database/core"
	"github.com/fullstackdev42/mp-emailer/database/decorators"
	"github.com/jonesrussell/loggo"
)

func ProvideDatabase(cfg *config.Config, logger loggo.LoggerInterface) (core.Interface, error) {
	retryConfig := dbconfig.NewDefaultRetryConfig()
	gormDB, err := dbconfig.ConnectWithRetry(cfg, retryConfig, logger)
	if err != nil {
		return nil, err
	}

	return decorators.NewLoggingDecorator(&core.DB{GormDB: gormDB}, logger), nil
}
