package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/fullstackdev42/mp-emailer/mocks"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jonesrussell/loggo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMigrate is a mock implementation of migrate.Migrate
type MockMigrate struct {
	mock.Mock
}

// Mock the Up method
func (m *MockMigrate) Up() error {
	args := m.Called()
	return args.Error(0)
}

// Mock the Close method
func (m *MockMigrate) Close() (error, error) {
	args := m.Called()
	return args.Error(0), args.Error(1)
}

// TestRunMigrations tests the RunMigrations function
func TestRunMigrations(t *testing.T) {
	// Test case 1: Migrations directory does not exist
	mockLogger := new(mocks.MockLoggerInterface)
	mockLogger.On("Info", mock.Anything).Return()
	err := RunMigrations("test_dsn", "/non/existent/path", mockLogger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "migrations directory does not exist")

	// Test case 2: Successful migration
	tempDir, err := os.MkdirTemp("", "test_migrations")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	mockMigrate := new(MockMigrate)
	mockMigrate.On("Up").Return(nil)
	mockMigrate.On("Close").Return(nil, nil)
	err = runMigrationsWithInstance(mockMigrate, mockLogger)
	assert.NoError(t, err)
	mockMigrate.AssertExpectations(t)
	mockLogger.AssertCalled(t, "Info", "Migrations completed successfully")

	// Test case 3: Migration error
	mockMigrate = new(MockMigrate)
	mockMigrate.On("Up").Return(fmt.Errorf("migration error"))
	mockMigrate.On("Close").Return(nil, nil)
	err = runMigrationsWithInstance(mockMigrate, mockLogger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error running migrations: migration error")
	mockMigrate.AssertExpectations(t)
}

// Update the migrator interface
type migrator interface {
	Up() error
	Close() (error, error)
}

// Update the runMigrationsWithInstance function
func runMigrationsWithInstance(m migrator, logger loggo.LoggerInterface) error {
	var migrationErr error
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		migrationErr = fmt.Errorf("error running migrations: %w", err)
	} else {
		logger.Info("Migrations completed successfully")
	}

	// Always call Close(), regardless of whether there was an error during migration
	closeErr, _ := m.Close()
	if closeErr != nil {
		return fmt.Errorf("error closing migrations: %w", closeErr)
	}

	// If there was a migration error, return it now
	if migrationErr != nil {
		return migrationErr
	}

	return nil
}
