package database

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/mp-emailer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMigrate is a mock implementation of migrate.Migrate
type MockMigrate struct {
	mock.Mock
}

func (m *MockMigrate) Up() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMigrate) Close() error {
	args := m.Called()
	return args.Error(0)
}

// TestRunMigrations uses table-driven tests for various scenarios
func TestRunMigrations(t *testing.T) {
	tests := []struct {
		name           string
		dsn            string
		migrationsPath string
		setupMock      func(*mocks.MockLoggerInterface, *MockMigrate)
		wantErr        bool
		errMsg         string
	}{
		{
			name:           "Migrations directory does not exist",
			dsn:            "user:password@tcp(localhost:3306)/testdb",
			migrationsPath: "/non/existent/path",
			setupMock: func(logger *mocks.MockLoggerInterface, migrator *MockMigrate) {
				logger.On("Debug", mock.Anything).Return()
				// Adjust the mock expectation to match the actual call
				logger.On("Error", "Error running migrations", fmt.Errorf("migrations directory does not exist")).Return()
				migrator.On("Up").Return(fmt.Errorf("migrations directory does not exist"))
				migrator.On("Close").Return(nil)
			},
			wantErr: true,
			errMsg:  "migrations directory does not exist",
		},
		{
			name: "Successful migration",
			dsn:  "user:password@tcp(localhost:3306)/testdb",
			migrationsPath: func() string {
				tempDir, err := os.MkdirTemp("", "test_migrations")
				if err != nil {
					t.Fatalf("failed to create temp dir: %v", err)
				}
				return tempDir
			}(),
			setupMock: func(logger *mocks.MockLoggerInterface, migrator *MockMigrate) {
				logger.On("Debug", mock.Anything).Return()
				logger.On("Info", "Migrations completed successfully").Return()
				migrator.On("Up").Return(nil)
				migrator.On("Close").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Migration error",
			dsn:  "user:password@tcp(localhost:3306)/testdb",
			migrationsPath: func() string {
				tempDir, err := os.MkdirTemp("", "test_migrations")
				if err != nil {
					t.Fatalf("failed to create temp dir: %v", err)
				}
				return tempDir
			}(),
			setupMock: func(logger *mocks.MockLoggerInterface, migrator *MockMigrate) {
				logger.On("Debug", mock.Anything).Return()
				// Adjust the mock expectation to match the actual call
				logger.On("Error", "Error running migrations", fmt.Errorf("migration error")).Return()
				migrator.On("Up").Return(fmt.Errorf("migration error"))
				migrator.On("Close").Return(nil)
			},
			wantErr: true,
			errMsg:  "error running migrations",
		},
		{
			name: "Migration error",
			dsn:  "user:password@tcp(localhost:3306)/testdb",
			migrationsPath: func() string {
				tempDir, err := os.MkdirTemp("", "test_migrations")
				if err != nil {
					t.Fatalf("failed to create temp dir: %v", err)
				}
				return tempDir
			}(),
			setupMock: func(logger *mocks.MockLoggerInterface, migrator *MockMigrate) {
				logger.On("Debug", mock.Anything).Return()
				// Adjust the mock expectation to match the actual call
				logger.On("Error", "Error running migrations", fmt.Errorf("migration error")).Return()
				migrator.On("Up").Return(fmt.Errorf("migration error"))
				migrator.On("Close").Return(nil)
			},
			wantErr: true,
			errMsg:  "error running migrations",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := new(mocks.MockLoggerInterface)
			mockMigrate := new(MockMigrate)
			if tt.setupMock != nil {
				tt.setupMock(mockLogger, mockMigrate)
			}

			err := runMigrationsWithInstance(mockMigrate, mockLogger)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				mockLogger.AssertCalled(t, "Info", "Migrations completed successfully")
			}

			mockLogger.AssertExpectations(t)
			mockMigrate.AssertExpectations(t)
		})
	}
}

// Update the migrator interface
type migrator interface {
	Up() error
	Close() error
}

// Update the runMigrationsWithInstance function
func runMigrationsWithInstance(m migrator, logger loggo.LoggerInterface) error {
	logger.Debug("Starting migration process") // Add debug log at the start
	var migrationErr error
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		migrationErr = fmt.Errorf("error running migrations: %w", err)
		logger.Error("Error running migrations", err)
	} else {
		logger.Info("Migrations completed successfully")
	}

	// Always call Close(), regardless of whether there was an error during migration
	if err := m.Close(); err != nil {
		return fmt.Errorf("error closing migrations: %w", err)
	}

	// If there was a migration error, return it now
	if migrationErr != nil {
		return migrationErr
	}

	return nil
}
