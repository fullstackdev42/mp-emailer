package database

import (
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/mp-emailer/config"
	dbconfig "github.com/jonesrussell/mp-emailer/database/config"
	"github.com/jonesrussell/mp-emailer/database/core"
	"github.com/jonesrussell/mp-emailer/database/decorators"
)

func ProvideDatabase(cfg *config.Config, logger loggo.LoggerInterface, retryConfig *dbconfig.RetryConfig, connector dbconfig.DatabaseConnector) (core.Interface, error) {
	gormDB, err := dbconfig.ConnectWithRetry(cfg, retryConfig, logger, connector)
	if err != nil {
		return nil, err
	}

	return decorators.NewLoggingDecorator(&core.DB{GormDB: gormDB}, logger), nil
}
