package campaign_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jonesrussell/mp-emailer/campaign"
	mockdb "github.com/jonesrussell/mp-emailer/mocks/database"
	"github.com/jonesrussell/mp-emailer/shared"
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
					{BaseModel: shared.BaseModel{ID: uuid.New()}, Name: "Campaign 1"},
					{BaseModel: shared.BaseModel{ID: uuid.New()}, Name: "Campaign 2"},
				}

				s.mockDB.EXPECT().FindAll(
					mock.Anything,
					mock.AnythingOfType("*[]campaign.Campaign"),
					"1=1",
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

func (s *RepositoryTestSuite) TestUpdate() {
	s.Run("successful_update", func() {
		// Test data
		id := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		ownerID := uuid.New()
		dto := &campaign.UpdateCampaignDTO{
			ID:          id,
			Name:        "Updated Campaign",
			Description: "Updated Description",
			Template:    "Updated Template",
		}

		// Create existing campaign
		existingCampaign := &campaign.Campaign{
			BaseModel: shared.BaseModel{ID: id},
			OwnerID:   ownerID,
		}

		// Mock FindOne - Note the exact parameter matching
		s.mockDB.On("FindOne",
			mock.Anything,
			mock.MatchedBy(func(_ *campaign.Campaign) bool {
				// The campaign passed to FindOne will be empty initially
				return true
			}),
			"id = ?",
			id,
		).Run(func(args mock.Arguments) {
			// Set the values in the campaign object that's passed in
			campaign := args.Get(1).(*campaign.Campaign)
			campaign.ID = existingCampaign.ID
			campaign.OwnerID = existingCampaign.OwnerID
		}).Return(nil)

		// Mock Update - Note the exact parameter matching
		s.mockDB.On("Update",
			mock.Anything,
			mock.MatchedBy(func(c *campaign.Campaign) bool {
				return c.ID == id &&
					c.Name == dto.Name &&
					c.Description == dto.Description &&
					c.Template == dto.Template &&
					c.OwnerID == ownerID
			}),
		).Return(nil)

		// Execute update
		err := s.repo.Update(context.Background(), dto)

		// Assertions
		s.NoError(err)
		s.mockDB.AssertExpectations(s.T())
	})
}

func (s *RepositoryTestSuite) TestDelete() {
	tests := []struct {
		name    string
		dto     campaign.DeleteCampaignDTO
		wantErr bool
	}{
		{
			name: "successful delete",
			dto: campaign.DeleteCampaignDTO{
				ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest() // Reset mock for each test case

			// Mock Delete call
			s.mockDB.EXPECT().Delete(
				mock.Anything,
				mock.AnythingOfType("*campaign.Campaign"),
			).Return(nil)

			err := s.repo.Delete(context.Background(), tt.dto)

			if tt.wantErr {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
			}
		})
	}
}

func (s *RepositoryTestSuite) TestGetByID() {
	tests := []struct {
		name    string
		dto     campaign.GetCampaignDTO
		setup   func()
		wantErr bool
	}{
		{
			name: "successful retrieval by ID",
			dto: campaign.GetCampaignDTO{
				ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			},
			setup: func() {
				expectedCampaign := &campaign.Campaign{
					BaseModel: shared.BaseModel{
						ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					},
					Name:        "Test Campaign",
					Description: "Test Description",
					Template:    "Test Template",
				}

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
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest() // Reset mock for each test case
			if tt.setup != nil {
				tt.setup()
			}

			result, err := s.repo.GetByID(context.Background(), tt.dto)

			if tt.wantErr {
				assert.Error(s.T(), err)
				assert.Nil(s.T(), result)
			} else {
				assert.NoError(s.T(), err)
				assert.NotNil(s.T(), result)
				assert.Equal(s.T(), tt.dto.ID, result.ID)
			}
		})
	}
}
