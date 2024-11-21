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
	ctx := context.Background()

	// Set up expectations with correct types:
	// Debug(string, string, time.Duration, string, int)
	logger.On("Debug",
		"Starting session cleanup routine",
		"interval", interval, // time.Duration
		"maxAge", maxAge, // int
	).Return()

	// Debug(string)
	logger.On("Debug", "Running session cleanup").Return()
	store.On("Cleanup", mock.AnythingOfType("time.Time")).Return(nil)

	cleaner := session.NewCleaner(store, interval, maxAge, logger)

	// Act
	cleaner.StartCleanup(ctx)

	// Wait for at least one cleanup cycle
	time.Sleep(20 * time.Millisecond)

	// Cleanup
	cleaner.StopCleanup()

	// Assert
	mock.AssertExpectationsForObjects(t, store, logger)
}

func TestStopCleanup(t *testing.T) {
	// Arrange
	store := mocksSession.NewMockStore(t)
	logger := mocks.NewMockLoggerInterface(t)
	interval := 100 * time.Millisecond
	maxAge := 3600
	ctx := context.Background()

	// Set up expectations
	logger.EXPECT().Debug("Starting session cleanup routine",
		"interval", interval,
		"maxAge", maxAge).Return()
	logger.EXPECT().Debug("Running session cleanup").Return()
	store.EXPECT().Cleanup(mock.AnythingOfType("time.Time")).Return(nil)

	cleaner := session.NewCleaner(store, interval, maxAge, logger)

	// Act
	cleaner.StartCleanup(ctx)

	// Wait a moment to ensure goroutine is running
	time.Sleep(150 * time.Millisecond)

	// Stop the cleanup
	cleaner.StopCleanup()

	// Wait a moment to ensure goroutine has stopped
	time.Sleep(150 * time.Millisecond)

	// Assert
	mock.AssertExpectationsForObjects(t, store, logger)
}
