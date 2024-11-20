package campaign_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jonesrussell/mp-emailer/campaign"
	mocks "github.com/jonesrussell/mp-emailer/mocks"
	mocksCampaign "github.com/jonesrussell/mp-emailer/mocks/campaign"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CampaignServiceTestSuite struct {
	suite.Suite
	service    *campaign.Service
	mockRepo   *mocksCampaign.MockRepositoryInterface
	validate   *validator.Validate
	mockLogger *mocks.MockLoggerInterface
}

func (s *CampaignServiceTestSuite) SetupTest() {
	s.mockRepo = new(mocksCampaign.MockRepositoryInterface)
	s.validate = validator.New()
	s.mockLogger = new(mocks.MockLoggerInterface)

	// Register the UUID validator
	err := s.validate.RegisterValidation("uuid4", func(fl validator.FieldLevel) bool {
		uuidStr := fl.Field().String()
		_, err := uuid.Parse(uuidStr)
		return err == nil
	})

	if err != nil {
		s.T().Fatalf("failed to register uuid4 validator: %v", err)
	}

	s.service = campaign.NewService(s.mockRepo, s.validate, s.mockLogger).(*campaign.Service)

	s.mockRepo.On("GetByID",
		mock.Anything,
		mock.MatchedBy(func(_ campaign.GetCampaignParams) bool {
			return true
		}),
	).Return(nil, fmt.Errorf("campaign not found"))
}

func TestCampaignServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CampaignServiceTestSuite))
}

func (s *CampaignServiceTestSuite) TestCreateCampaign() {
	tests := []struct {
		name    string
		dto     *campaign.CreateCampaignDTO
		setup   func()
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
			setup: func() {
				s.mockRepo.EXPECT().Create(
					mock.Anything,
					mock.AnythingOfType("*campaign.CreateCampaignDTO"),
				).Return(&campaign.Campaign{
					Name:        "Test Campaign",
					Description: "Test Description",
					Template:    "Test Template",
					OwnerID:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				}, nil)

				s.mockLogger.EXPECT().Info(
					"Campaign created successfully",
					"id",
					mock.AnythingOfType("uuid.UUID"),
				).Return()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.mockRepo.ExpectedCalls = nil
			s.mockRepo.Calls = nil
			s.mockLogger.ExpectedCalls = nil
			s.mockLogger.Calls = nil

			tt.setup()
			got, err := s.service.CreateCampaign(context.Background(), tt.dto)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
			s.Equal(tt.dto.Name, got.Name)
			s.Equal(tt.dto.Description, got.Description)
			s.Equal(tt.dto.Template, got.Template)
			s.Equal(tt.dto.OwnerID, got.OwnerID)
			s.mockRepo.AssertExpectations(s.T())
			s.mockLogger.AssertExpectations(s.T())
		})
	}
}

func (s *CampaignServiceTestSuite) TestComposeEmail() {
	tests := []struct {
		name    string
		params  campaign.ComposeEmailParams
		want    string
		wantErr bool
	}{
		{
			name: "successful email composition",
			params: campaign.ComposeEmailParams{
				MP: campaign.Representative{
					Name:  "John Doe",
					Email: "john@example.com",
				},
				Campaign: &campaign.Campaign{
					Template: "Dear {{MP's Name}}, This is a test email. Your email is {{MPEmail}}. Date: {{Date}}",
				},
				UserData: map[string]string{
					"CustomField": "Custom Value",
				},
			},
			want:    "Dear John Doe, This is a test email. Your email is john@example.com. Date: " + time.Now().Format("2006-01-02"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.service.ComposeEmail(context.Background(), tt.params)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
			s.Equal(tt.want, got)
		})
	}
}

func (s *CampaignServiceTestSuite) TestGetCampaignByID() {
	tests := []struct {
		name    string
		id      uuid.UUID
		setup   func()
		want    *campaign.Campaign
		wantErr bool
	}{
		{
			name: "successful retrieval",
			id:   uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			setup: func() {
				s.mockRepo.On("GetByID",
					mock.Anything,
					campaign.GetCampaignParams{
						ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					},
				).Return(&campaign.Campaign{
					Name:        "Test Campaign",
					Description: "Test Description",
					Template:    "Test Template",
					OwnerID:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				}, nil)
			},
			want: &campaign.Campaign{
				Name:        "Test Campaign",
				Description: "Test Description",
				Template:    "Test Template",
				OwnerID:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			},
			wantErr: false,
		},
		{
			name: "campaign not found",
			id:   uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
			setup: func() {
				s.mockRepo.EXPECT().GetByID(
					mock.Anything,
					mock.MatchedBy(func(params campaign.GetCampaignParams) bool {
						return params.ID == uuid.MustParse("123e4567-e89b-12d3-a456-426614174002")
					}),
				).Return(nil, fmt.Errorf("campaign not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.mockRepo.ExpectedCalls = nil
			s.mockRepo.Calls = nil
			tt.setup()
			got, err := s.service.GetCampaignByID(context.Background(), campaign.GetCampaignParams{ID: tt.id})
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
			s.Equal(tt.want, got)
			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *CampaignServiceTestSuite) TestDeleteCampaign() {
	tests := []struct {
		name    string
		id      uuid.UUID
		setup   func()
		wantErr bool
	}{
		{
			name: "successful deletion",
			id:   uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			setup: func() {
				s.mockRepo.EXPECT().GetByID(
					mock.Anything,
					mock.MatchedBy(func(params campaign.GetCampaignParams) bool {
						return params.ID == uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
					}),
				).Return(&campaign.Campaign{}, nil)

				s.mockRepo.EXPECT().Delete(
					mock.Anything,
					mock.MatchedBy(func(dto campaign.DeleteCampaignDTO) bool {
						return dto.ID == uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
					}),
				).Return(nil)

				s.mockLogger.EXPECT().Info(
					"Campaign deleted successfully",
					"id",
					uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				).Return()
			},
			wantErr: false,
		},
		{
			name: "campaign not found",
			id:   uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
			setup: func() {
				s.mockRepo.EXPECT().GetByID(
					mock.Anything,
					mock.MatchedBy(func(params campaign.GetCampaignParams) bool {
						return params.ID == uuid.MustParse("123e4567-e89b-12d3-a456-426614174002")
					}),
				).Return(nil, fmt.Errorf("campaign not found"))

				s.mockLogger.EXPECT().Error(
					"Failed to fetch campaign for deletion",
					mock.MatchedBy(func(err error) bool {
						return err.Error() == "campaign not found"
					}),
					"id",
					uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
				).Return()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.mockRepo.ExpectedCalls = nil
			s.mockRepo.Calls = nil

			tt.setup()
			err := s.service.DeleteCampaign(context.Background(), campaign.DeleteCampaignDTO{ID: tt.id})
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
			s.mockRepo.AssertExpectations(s.T())
		})
	}
}
