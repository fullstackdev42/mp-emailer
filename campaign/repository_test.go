package campaign

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fullstackdev42/mp-emailer/database"
	"github.com/fullstackdev42/mp-emailer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type RepositoryTestSuite struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	db         *database.DB
	repo       *Repository
	mockLogger *mocks.MockLoggerInterface
}

// SetupTest sets up the test environment
func (s *RepositoryTestSuite) SetupTest() {
	// Use an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(s.T(), err)

	s.mockLogger = new(mocks.MockLoggerInterface)

	// Setup expected logger calls
	s.mockLogger.On("Debug", "Connecting to database", []interface{}{}).Return(nil)

	// Migrate the schema for testing
	err = db.AutoMigrate(&Campaign{})
	assert.NoError(s.T(), err)

	s.db = &database.DB{GormDB: db}
	s.repo = &Repository{db: s.db}
}

// TearDownTest tears down the test environment
func (s *RepositoryTestSuite) TearDownTest() {
	s.mock.ExpectClose()
	sqlDB, err := s.db.GormDB.DB()
	if err != nil {
		s.T().Fatal(err)
	}
	sqlDB.Close()
}

// TestRepositoryTestSuite runs the Repository test suite
func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

// TestCreate tests the Create method of Repository
func (s *RepositoryTestSuite) TestCreate() {
	tests := []struct {
		name    string
		dto     *CreateCampaignDTO
		setup   func()
		wantErr bool
	}{
		{
			name: "successful creation",
			dto: &CreateCampaignDTO{
				Name:        "Test Campaign",
				Description: "Test Description",
				Template:    "Test Template",
				OwnerID:     "123e4567-e89b-12d3-a456-426614174000",
			},
			setup: func() {
				s.mock.ExpectExec("INSERT INTO campaigns").
					WithArgs("Test Campaign", "Test Description", "Test Template", "123e4567-e89b-12d3-a456-426614174000").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "database error",
			dto: &CreateCampaignDTO{
				Name:        "Test Campaign",
				Description: "Test Description",
				Template:    "Test Template",
				OwnerID:     "123e4567-e89b-12d3-a456-426614174000",
			},
			setup: func() {
				s.mock.ExpectExec("INSERT INTO campaigns").
					WithArgs("Test Campaign", "Test Description", "Test Template", "123e4567-e89b-12d3-a456-426614174000").
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setup()

			campaign, err := s.repo.Create(tt.dto)
			if tt.wantErr {
				assert.Error(s.T(), err)
				assert.Nil(s.T(), campaign)
			} else {
				assert.NoError(s.T(), err)
				assert.NotNil(s.T(), campaign)
				assert.Equal(s.T(), tt.dto.Name, campaign.Name)
			}

			assert.NoError(s.T(), s.mock.ExpectationsWereMet())
		})
	}
}

// TestGetAll tests the GetAll method of Repository
func (s *RepositoryTestSuite) TestGetAll() {
	tests := []struct {
		name      string
		setup     func()
		wantCount int
		wantErr   bool
	}{
		{
			name: "successful retrieval",
			setup: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "template", "owner_id", "created_at", "updated_at"}).
					AddRow(1, "Campaign 1", "Desc 1", "Template 1", 1, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05")).
					AddRow(2, "Campaign 2", "Desc 2", "Template 2", 1, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
				s.mock.ExpectQuery("SELECT (.+) FROM campaigns").WillReturnRows(rows)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "database error",
			setup: func() {
				s.mock.ExpectQuery("SELECT (.+) FROM campaigns").
					WillReturnError(errors.New("database error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setup()

			campaigns, err := s.repo.GetAll()
			if tt.wantErr {
				assert.Error(s.T(), err)
				assert.Nil(s.T(), campaigns)
			} else {
				assert.NoError(s.T(), err)
				assert.Len(s.T(), campaigns, tt.wantCount)
			}

			assert.NoError(s.T(), s.mock.ExpectationsWereMet())
		})
	}
}
