package campaign_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jonesrussell/mp-emailer/campaign"
	mocksDatabase "github.com/jonesrussell/mp-emailer/mocks/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	mockDB *mocksDatabase.MockInterface
	repo   campaign.RepositoryInterface
}

// SetupTest sets up the test environment
func (s *RepositoryTestSuite) SetupTest() {
	s.mockDB = mocksDatabase.NewMockInterface(s.T())
	s.repo = campaign.NewRepository(s.mockDB).(*campaign.Repository)
}

// TearDownTest tears down the test environment
func (s *RepositoryTestSuite) TearDownTest() {
	// Clean up database connection
	if s.mockDB != nil {
		s.mockDB.AssertExpectations(s.T())
	}
}

// TestRepositoryTestSuite runs the Repository test suite
func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

// TestCreate tests the Create method of Repository
func (s *RepositoryTestSuite) TestCreate() {
	tests := []struct {
		name    string
		setup   func()
		dto     *campaign.CreateCampaignDTO
		wantErr bool
	}{
		{
			name: "successful creation",
			setup: func() {
				s.mockDB.On("Create", mock.AnythingOfType("*campaign.Campaign")).
					Return(nil)
			},
			dto: &campaign.CreateCampaignDTO{
				Name:        "Test Campaign",
				Description: "Test Description",
				Template:    "Test Template",
				OwnerID:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			},
			wantErr: false,
		},
		{
			name: "database error",
			setup: func() {
				s.mockDB.On("Create", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
					return campaign.Name == "Test Campaign" &&
						campaign.Description == "Test Description" &&
						campaign.Template == "Test Template" &&
						campaign.OwnerID == uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				})).Return(errors.New("db error"))
			},
			dto: &campaign.CreateCampaignDTO{
				Name:        "Test Campaign",
				Description: "Test Description",
				Template:    "Test Template",
				OwnerID:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest() // Reset mock for each test case
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
				mockResult := &mocksDatabase.MockResult{}
				mockResult.On("Error").Return(nil)
				mockResult.On("Scan", mock.AnythingOfType("*[]campaign.Campaign")).
					Run(func(args mock.Arguments) {
						dest := args.Get(0).(*[]campaign.Campaign)
						*dest = []campaign.Campaign{
							{Name: "Campaign 1"},
							{Name: "Campaign 2"},
						}
					}).
					Return(mockResult)

				s.mockDB.On("Query", "SELECT * FROM campaigns").
					Return(mockResult)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "database error",
			setup: func() {
				mockResult := &mocksDatabase.MockResult{}
				mockResult.On("Error").Return(errors.New("db error"))
				mockResult.On("Scan", mock.AnythingOfType("*[]campaign.Campaign")).
					Return(mockResult)

				s.mockDB.On("Query", "SELECT * FROM campaigns").
					Return(mockResult)
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest() // Reset mock for each test case
			tt.setup()

			campaigns, err := s.repo.GetAll()
			if tt.wantErr {
				assert.Error(s.T(), err)
				assert.Nil(s.T(), campaigns)
			} else {
				assert.NoError(s.T(), err)
				assert.Len(s.T(), campaigns, tt.wantCount)
			}
		})
	}
}

// TestUpdate tests the Update method
func (s *RepositoryTestSuite) TestUpdate() {
	testID := uuid.New()
	tests := []struct {
		name    string
		setup   func()
		dto     *campaign.UpdateCampaignDTO
		wantErr bool
	}{
		{
			name: "successful update",
			setup: func() {
				s.mockDB.On("Exists", mock.AnythingOfType("*campaign.Campaign"), "id = ?", testID).
					Return(true, nil)
				s.mockDB.On("Exec",
					"UPDATE campaigns SET name = ?, description = ?, template = ? WHERE id = ?",
					"Updated Campaign", "Updated Description", "Updated Template", testID).
					Return(nil)
			},
			dto: &campaign.UpdateCampaignDTO{
				ID:          testID,
				Name:        "Updated Campaign",
				Description: "Updated Description",
				Template:    "Updated Template",
			},
			wantErr: false,
		},
		{
			name: "campaign not found",
			setup: func() {
				s.mockDB.On("Exists", mock.AnythingOfType("*campaign.Campaign"), "id = ?", testID).
					Return(false, nil)
			},
			dto: &campaign.UpdateCampaignDTO{
				ID: testID,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.setup()

			err := s.repo.Update(tt.dto)
			if tt.wantErr {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
			}
		})
	}
}

// TestDelete tests the Delete method
func (s *RepositoryTestSuite) TestDelete() {
	testID := uuid.New()
	tests := []struct {
		name    string
		setup   func()
		dto     campaign.DeleteCampaignDTO
		wantErr bool
	}{
		{
			name: "successful deletion",
			setup: func() {
				s.mockDB.On("Exists", mock.AnythingOfType("*campaign.Campaign"), "id = ?", testID).
					Return(true, nil)
				s.mockDB.On("Exec", "DELETE FROM campaigns WHERE id = ?", testID).
					Return(nil)
			},
			dto: campaign.DeleteCampaignDTO{
				ID: testID,
			},
			wantErr: false,
		},
		{
			name: "campaign not found",
			setup: func() {
				s.mockDB.On("Exists", mock.AnythingOfType("*campaign.Campaign"), "id = ?", testID).
					Return(false, nil)
			},
			dto: campaign.DeleteCampaignDTO{
				ID: testID,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.setup()

			err := s.repo.Delete(tt.dto)
			if tt.wantErr {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
			}
		})
	}
}

// TestGetByID tests the GetByID method
func (s *RepositoryTestSuite) TestGetByID() {
	testID := uuid.New()
	tests := []struct {
		name    string
		setup   func()
		dto     campaign.GetCampaignDTO
		wantErr bool
	}{
		{
			name: "successful retrieval",
			setup: func() {
				s.mockDB.On("FindOne", mock.AnythingOfType("*campaign.Campaign"), "id = ?", testID).
					Run(func(args mock.Arguments) {
						campaign := args.Get(0).(*campaign.Campaign)
						campaign.ID = testID
						campaign.Name = "Test Campaign"
					}).
					Return(nil)
			},
			dto: campaign.GetCampaignDTO{
				ID: testID,
			},
			wantErr: false,
		},
		{
			name: "campaign not found",
			setup: func() {
				s.mockDB.On("FindOne", mock.AnythingOfType("*campaign.Campaign"), "id = ?", testID).
					Return(errors.New("not found"))
			},
			dto: campaign.GetCampaignDTO{
				ID: testID,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.setup()

			campaign, err := s.repo.GetByID(tt.dto)
			if tt.wantErr {
				assert.Error(s.T(), err)
				assert.Nil(s.T(), campaign)
			} else {
				assert.NoError(s.T(), err)
				assert.NotNil(s.T(), campaign)
				assert.Equal(s.T(), testID, campaign.ID)
			}
		})
	}
}
