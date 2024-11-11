package config

import (
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/jonesrussell/loggo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectWithRetry(cfg *config.Config, retryConfig *RetryConfig, logger loggo.LoggerInterface) (*gorm.DB, error) {
	var db *gorm.DB

	operation := func() error {
		var err error
		db, err = Connect(cfg)
		if err != nil {
			logger.Error("Failed to connect to database", err)
			return err
		}
		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = retryConfig.InitialInterval
	expBackoff.MaxInterval = retryConfig.MaxInterval
	expBackoff.MaxElapsedTime = retryConfig.MaxElapsedTime
	expBackoff.Multiplier = retryConfig.MultiplicationFactor

	err := backoff.Retry(operation, expBackoff)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after retries: %w", err)
	}

	logger.Info("Successfully connected to database after retry")
	return db, nil
}

func Connect(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
