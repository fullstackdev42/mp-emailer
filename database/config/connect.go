package config

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/mp-emailer/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConnector interface {
	Connect(ctx context.Context, cfg *config.Config) (*gorm.DB, error)
}

type DefaultConnector struct{}

func (d *DefaultConnector) Connect(ctx context.Context, cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db.WithContext(ctx), nil
}

func ConnectWithRetry(ctx context.Context, cfg *config.Config, retryConfig *RetryConfig, logger loggo.LoggerInterface, connector DatabaseConnector) (*gorm.DB, error) {
	var db *gorm.DB

	operation := func() error {
		// Check if context is cancelled before attempting connection
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var err error
			db, err = connector.Connect(ctx, cfg)
			if err != nil {
				logger.Error("Failed to connect to database", err)
				return err
			}
			return nil
		}
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = retryConfig.InitialInterval
	expBackoff.MaxInterval = retryConfig.MaxInterval
	expBackoff.MaxElapsedTime = retryConfig.MaxElapsedTime
	expBackoff.Multiplier = retryConfig.MultiplicationFactor

	// Use backoff.WithContext to make the retry operation context-aware
	err := backoff.Retry(operation, backoff.WithContext(expBackoff, ctx))
	if err != nil {
		if err == context.Canceled {
			return nil, fmt.Errorf("database connection cancelled: %w", err)
		}
		if err == context.DeadlineExceeded {
			return nil, fmt.Errorf("database connection timeout: %w", err)
		}
		return nil, fmt.Errorf("failed to connect to database after retries: %w", err)
	}

	logger.Info("Successfully connected to database after retry")
	return db, nil
}

// RetryConfig holds configuration for connection retry behavior
type RetryConfig struct {
	InitialInterval      time.Duration
	MaxInterval          time.Duration
	MaxElapsedTime       time.Duration
	MultiplicationFactor float64
	MaxAttempts          int
}

// NewDefaultRetryConfig returns default retry configuration
func NewDefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		InitialInterval:      100 * time.Millisecond,
		MaxInterval:          10 * time.Second,
		MaxElapsedTime:       1 * time.Minute,
		MultiplicationFactor: 2.0,
		MaxAttempts:          5,
	}
}
