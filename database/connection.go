package database

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ConnectionConfig struct {
	DSN                  string
	MaxRetries           int
	InitialInterval      time.Duration
	MaxInterval          time.Duration
	MaxElapsedTime       time.Duration
	MultiplicationFactor float64
}

func NewConnection(ctx context.Context, cfg ConnectionConfig) (Database, error) {
	var db *gorm.DB

	operation := func() error {
		var err error
		db, err = gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
		return err
	}

	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.InitialInterval = cfg.InitialInterval
	backoffConfig.MaxInterval = cfg.MaxInterval
	backoffConfig.MaxElapsedTime = cfg.MaxElapsedTime
	backoffConfig.Multiplier = cfg.MultiplicationFactor

	if err := backoff.Retry(operation, backoff.WithContext(backoffConfig, ctx)); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &GormDB{db: db}, nil
}
