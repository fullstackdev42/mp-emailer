package session_test

import (
	"context"
	"testing"
	"time"

	"github.com/jonesrussell/mp-emailer/mocks"
	mocksSession "github.com/jonesrussell/mp-emailer/mocks/session"
	"github.com/jonesrussell/mp-emailer/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewCleaner(t *testing.T) {
	// Arrange
	store := mocksSession.NewMockStore(t)
	logger := mocks.NewMockLoggerInterface(t)
	interval := 15 * time.Minute
	maxAge := 3600

	// Act
	cleaner := session.NewCleaner(store, interval, maxAge, logger)

	// Assert
	assert.NotNil(t, cleaner)
}

func TestCleanup(t *testing.T) {
	// Arrange
	store := mocksSession.NewMockStore(t)
	logger := mocks.NewMockLoggerInterface(t)
	interval := 15 * time.Millisecond
	maxAge := 3600

	// Set up expectations
	logger.EXPECT().Debug("Starting session cleanup routine",
		"interval", interval,
		"maxAge", maxAge,
	).Return()

	logger.EXPECT().Debug("Running session cleanup").Return()
	store.EXPECT().Cleanup(mock.AnythingOfType("time.Time")).Return(nil)

	cleaner := session.NewCleaner(store, interval, maxAge, logger)

	// Act
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	cleaner.StartCleanup(ctx)

	// Wait for at least one cleanup cycle
	time.Sleep(20 * time.Millisecond)
}

func TestStopCleanup(t *testing.T) {
	// Arrange
	store := mocksSession.NewMockStore(t)
	logger := mocks.NewMockLoggerInterface(t)
	interval := 100 * time.Millisecond
	maxAge := 3600

	cleaner := session.NewCleaner(store, interval, maxAge, logger)

	logger.EXPECT().Debug("Starting session cleanup routine",
		"interval", interval,
		"maxAge", maxAge,
	).Return()

	logger.EXPECT().Debug("Context cancelled, stopping cleanup routine").Return()

	// Act
	ctx, cancel := context.WithCancel(context.Background())
	cleaner.StartCleanup(ctx)

	// Wait briefly then cancel
	time.Sleep(50 * time.Millisecond)
	cancel()

	// Wait for cleanup to stop
	time.Sleep(50 * time.Millisecond)
}
