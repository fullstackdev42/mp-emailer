package campaign

import (
	"fmt"
	"testing"

	mocksDatabase "github.com/fullstackdev42/mp-emailer/mocks/database"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	mockDB *mocksDatabase.MockInterface
	repo   *Repository
}

// SetupTest sets up the test environment
func (s *RepositoryTestSuite) SetupTest() {
	s.mockDB = mocksDatabase.NewMockInterface(s.T())
	s.repo = NewRepository(s.mockDB).(*Repository)
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
		dto     *CreateCampaignDTO
		wantErr bool
	}{
		{
			name: "successful creation",
			setup: func() {
				s.mockDB.On("Create", mock.AnythingOfType("*campaign.Campaign")).
					Return(nil)
			},
			dto: &CreateCampaignDTO{
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
				s.mockDB.On("Create", mock.AnythingOfType("*campaign.Campaign")).
					Return(fmt.Errorf("db error"))
			},
			dto: &CreateCampaignDTO{
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
	campaigns := []*Campaign{
		{Name: "Campaign 1"},
		{Name: "Campaign 2"},
	}

	tests := []struct {
		name      string
		setup     func()
		wantCount int
		wantErr   bool
	}{
		{
			name: "successful retrieval",
			setup: func() {
				s.mockDB.On("Find", mock.AnythingOfType("*[]*campaign.Campaign")).
					Run(func(args mock.Arguments) {
						result := args.Get(0).(*[]*Campaign)
						*result = campaigns
					}).
					Return(nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "database error",
			setup: func() {
				s.mockDB.On("Find", mock.AnythingOfType("*[]*campaign.Campaign")).
					Return(fmt.Errorf("db error"))
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
		})
	}
}
