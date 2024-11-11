package database_test

import (
	"fmt"
	"testing"

	"github.com/fullstackdev42/mp-emailer/config"
	dbconfig "github.com/fullstackdev42/mp-emailer/database/config"
	"github.com/jonesrussell/loggo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
	connectAttempts int
}

func (m *MockDB) Connect() error {
	m.connectAttempts++
	args := m.Called()
	return args.Error(0)
}

func TestDatabaseRetryMechanism(t *testing.T) {
	logger, _ := loggo.NewLogger("../storage/logs/database_test.log", loggo.LevelDebug)
	mockDB := new(MockDB)

	// Test successful connection after 2 retries
	t.Run("successful connection after retries", func(t *testing.T) {
		mockDB.connectAttempts = 0 // Reset counter
		mockDB.On("Connect").
			Return(fmt.Errorf("connection error")).
			Times(2)
		mockDB.On("Connect").
			Return(nil).
			Once()

		retryConfig := dbconfig.NewDefaultRetryConfig()
		cfg := &config.Config{}

		_, err := dbconfig.ConnectWithRetry(cfg, retryConfig, logger) // Use mockDB here
		assert.NoError(t, err)
		assert.Equal(t, 3, mockDB.connectAttempts)
		mockDB.AssertExpectations(t)
	})

	// Test max retries exceeded
	t.Run("max retries exceeded", func(t *testing.T) {
		mockDB.connectAttempts = 0 // Reset counter
		mockDB.On("Connect").
			Return(fmt.Errorf("connection error")).
			Times(5)

		cfg := &config.Config{}
		retryConfig := dbconfig.NewDefaultRetryConfig()

		_, err := dbconfig.ConnectWithRetry(cfg, retryConfig, logger) // Use mockDB here
		assert.Error(t, err)
		assert.Equal(t, 5, mockDB.connectAttempts)
		mockDB.AssertExpectations(t)
	})
}
