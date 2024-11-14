package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fullstackdev42/mp-emailer/config"
	dbconfig "github.com/fullstackdev42/mp-emailer/database/config"
	"github.com/fullstackdev42/mp-emailer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MockConnector struct {
	mock.Mock
}

func (m *MockConnector) Connect(cfg *config.Config) (*gorm.DB, error) {
	args := m.Called(cfg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gorm.DB), args.Error(1)
}

func TestProvideDatabase(t *testing.T) {
	t.Run("successful database connection", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      mockDB,
			SkipInitializeWithVersion: true,
		}))
		require.NoError(t, err)

		connector := &MockConnector{}
		connector.On("Connect", mock.Anything).Return(gormDB, nil)

		logger := mocks.NewMockLoggerInterface(t)
		logger.EXPECT().Info("Successfully connected to database after retry").Return()

		cfg := &config.Config{
			DBHost:     "localhost",
			DBPort:     3306,
			DBUser:     "test",
			DBPassword: "test",
			DBName:     "testdb",
		}

		retryConfig := dbconfig.NewDefaultRetryConfig()

		db, err := ProvideDatabase(cfg, logger, retryConfig, connector)
		assert.NoError(t, err)
		assert.NotNil(t, db)
	})
}
