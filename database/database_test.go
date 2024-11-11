package database_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/fullstackdev42/mp-emailer/config"
	dbconfig "github.com/fullstackdev42/mp-emailer/database/config"
	"github.com/jonesrussell/loggo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of the database with a connectAttempts counter.
type MockDB struct {
	mock.Mock
	connectAttempts int
}

// Connect mocks the database connection attempt.
func (m *MockDB) Connect() error {
	m.connectAttempts++
	args := m.Called()
	return args.Error(0)
}

// MockConnectWithRetry simulates the retry logic using the mock database.
func MockConnectWithRetry(_ *config.Config, retryConfig *dbconfig.RetryConfig, _ loggo.LoggerInterface, db *MockDB) error {
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = retryConfig.InitialInterval
	expBackoff.MaxInterval = retryConfig.MaxInterval
	expBackoff.MaxElapsedTime = retryConfig.MaxElapsedTime
	expBackoff.Multiplier = retryConfig.MultiplicationFactor
	expBackoff.Reset()

	operation := func() error {
		return db.Connect()
	}

	err := backoff.Retry(operation, expBackoff)
	if err != nil {
		return fmt.Errorf("failed to connect to database after retries: %w", err)
	}
	return nil
}

func TestDatabaseRetryMechanism(t *testing.T) {
	logger, _ := loggo.NewLogger("../storage/logs/database-test.log", loggo.LevelDebug)
	mockDB := new(MockDB)

	t.Run("successful connection after retries", func(t *testing.T) {
		mockDB.connectAttempts = 0
		mockDB.On("Connect").
			Return(fmt.Errorf("connection error")).
			Times(2)
		mockDB.On("Connect").
			Return(nil).
			Once()

		retryConfig := &dbconfig.RetryConfig{
			InitialInterval:      10 * time.Millisecond,
			MaxInterval:          100 * time.Millisecond,
			MaxElapsedTime:       300 * time.Millisecond,
			MultiplicationFactor: 2.0,
		}
		cfg := &config.Config{}

		err := MockConnectWithRetry(cfg, retryConfig, logger, mockDB)
		assert.NoError(t, err)
		assert.Equal(t, 3, mockDB.connectAttempts)
		mockDB.AssertExpectations(t)
	})

	t.Run("max retries exceeded", func(t *testing.T) {
		mockDB.connectAttempts = 0
		mockDB.On("Connect").
			Return(fmt.Errorf("connection error")).
			Times(5)
		mockDB.On("Connect").
			Return(fmt.Errorf("connection error")).
			Once()

		cfg := &config.Config{}
		retryConfig := &dbconfig.RetryConfig{
			InitialInterval:      5 * time.Millisecond,
			MaxInterval:          20 * time.Millisecond,
			MaxElapsedTime:       100 * time.Millisecond,
			MultiplicationFactor: 2.0,
		}

		err := MockConnectWithRetry(cfg, retryConfig, logger, mockDB)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to connect to database after retries")
		assert.Equal(t, 6, mockDB.connectAttempts)
		mockDB.AssertExpectations(t)
	})
}
