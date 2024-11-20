package campaign_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jonesrussell/mp-emailer/campaign"
	mockdb "github.com/jonesrussell/mp-emailer/mocks/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	mockDB *mockdb.MockDatabase
	repo   campaign.RepositoryInterface
	params campaign.RepositoryParams
}

func (s *RepositoryTestSuite) SetupTest() {
	s.mockDB = mockdb.NewMockDatabase(s.T())
	s.params = campaign.RepositoryParams{
		DB: s.mockDB,
	}
	s.repo = campaign.NewRepository(s.params)
}

func TestCampaignRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) TestCreate() {
	tests := []struct {
		name    string
		dto     *campaign.CreateCampaignDTO
		wantErr bool
	}{
		{
			name: "successful creation",
			dto: &campaign.CreateCampaignDTO{
				Name:        "Test Campaign",
				Description: "Test Description",
				Template:    "Test Template",
				OwnerID:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest() // Reset mock for each test case

			expectedCampaign := &campaign.Campaign{
				Name:        tt.dto.Name,
				Description: tt.dto.Description,
				Template:    tt.dto.Template,
				OwnerID:     tt.dto.OwnerID,
			}

			// Mock Create call
			s.mockDB.EXPECT().Create(
				mock.Anything,
				mock.AnythingOfType("*campaign.Campaign"),
			).Run(func(_ context.Context, value interface{}) {
				if campaign, ok := value.(*campaign.Campaign); ok {
					campaign.ID = uuid.New() // Set an ID as the database would
				}
			}).Return(nil)

			// Mock FindOne call that happens after creation
			s.mockDB.EXPECT().FindOne(
				mock.Anything,
				mock.AnythingOfType("*campaign.Campaign"),
				"id = ?",
				mock.AnythingOfType("uuid.UUID"),
			).Run(func(_ context.Context, dest interface{}, _ string, _ ...interface{}) {
				if campaign, ok := dest.(*campaign.Campaign); ok {
					*campaign = *expectedCampaign
				}
			}).Return(nil)

			result, err := s.repo.Create(context.Background(), tt.dto)

			if tt.wantErr {
				assert.Error(s.T(), err)
				assert.Nil(s.T(), result)
			} else {
				assert.NoError(s.T(), err)
				assert.NotNil(s.T(), result)
				assert.Equal(s.T(), tt.dto.Name, result.Name)
			}
		})
	}
}

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
				campaigns := []campaign.Campaign{
					{Name: "Campaign 1"},
					{Name: "Campaign 2"},
				}

				s.mockDB.EXPECT().FindOne(
					mock.Anything,
					mock.AnythingOfType("*[]campaign.Campaign"),
					"",
				).Run(func(_ context.Context, dest interface{}, _ string, _ ...interface{}) {
					if destSlice, ok := dest.(*[]campaign.Campaign); ok {
						*destSlice = campaigns
					}
				}).Return(nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest() // Reset mock for each test case
			if tt.setup != nil {
				tt.setup()
			}

			campaigns, err := s.repo.GetAll(context.Background())
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
