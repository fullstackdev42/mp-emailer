package campaign_test

import (
	"testing"
	"time"

	"github.com/fullstackdev42/mp-emailer/campaign"
	mocks "github.com/fullstackdev42/mp-emailer/mocks"
	mocksCampaign "github.com/fullstackdev42/mp-emailer/mocks/campaign"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
}

func TestCampaignServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CampaignServiceTestSuite))
}

func (s *CampaignServiceTestSuite) TestCreateCampaign() {
	tests := []struct {
		name    string
		input   *campaign.CreateCampaignDTO
		mock    func(*campaign.CreateCampaignDTO)
		want    *campaign.Campaign
		wantErr bool
	}{
		{
			name: "successful creation",
			input: &campaign.CreateCampaignDTO{
				Name:        "Test Campaign",
				Description: "Test Description",
				Template:    "Test Template",
				OwnerID:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			},
			mock: func(input *campaign.CreateCampaignDTO) {
				expectedCampaign := &campaign.Campaign{
					Name:        "Test Campaign",
					Description: "Test Description",
					Template:    "Test Template",
					OwnerID:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Tokens:      []string{},
				}
				s.mockRepo.On("Create", input).Return(expectedCampaign, nil)
				s.mockLogger.On("Info", "Campaign created successfully", "id", mock.Anything).Return()
			},
			want: &campaign.Campaign{
				Name:        "Test Campaign",
				Description: "Test Description",
				Template:    "Test Template",
				OwnerID:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Tokens:      []string{},
			},
			wantErr: false,
		},
		// ... other test cases ...
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.mockRepo.ExpectedCalls = nil
			s.mockRepo.Calls = nil

			tt.mock(tt.input)
			got, err := s.service.CreateCampaign(tt.input)
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
			got, err := s.service.ComposeEmail(tt.params)
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
		params  campaign.GetCampaignParams
		mock    func(campaign.GetCampaignParams)
		want    *campaign.Campaign
		wantErr bool
	}{
		{
			name: "successful retrieval",
			params: campaign.GetCampaignParams{
				ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			},
			mock: func(params campaign.GetCampaignParams) {
				expectedCampaign := &campaign.Campaign{
					Name:        "Test Campaign",
					Description: "Test Description",
					Template:    "Test Template",
					Tokens:      []string{},
				}
				s.mockRepo.On("GetByID", campaign.GetCampaignDTO(params)).Return(expectedCampaign, nil)
			},
			want: &campaign.Campaign{
				Name:        "Test Campaign",
				Description: "Test Description",
				Template:    "Test Template",
				Tokens:      []string{},
			},
			wantErr: false,
		},
		{
			name: "campaign not found",
			params: campaign.GetCampaignParams{
				ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
			},
			mock: func(params campaign.GetCampaignParams) {
				s.mockRepo.On("GetByID", campaign.GetCampaignDTO(params)).Return(nil, campaign.ErrCampaignNotFound)
				s.mockLogger.On("Debug", "Campaign not found", "id", params.ID).Return()
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.mockRepo.ExpectedCalls = nil
			s.mockRepo.Calls = nil

			tt.mock(tt.params)
			got, err := s.service.GetCampaignByID(tt.params)
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
		dto     campaign.DeleteCampaignDTO
		mock    func(campaign.DeleteCampaignDTO)
		wantErr bool
	}{
		{
			name: "successful deletion",
			dto: campaign.DeleteCampaignDTO{
				ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			},
			mock: func(dto campaign.DeleteCampaignDTO) {
				c := &campaign.Campaign{}
				s.mockRepo.On("GetByID", campaign.GetCampaignDTO(dto)).Return(c, nil)
				s.mockRepo.On("Delete", dto).Return(nil)
				s.mockLogger.On("Info", "Campaign deleted successfully", "id", dto.ID).Return()
			},
			wantErr: false,
		},
		{
			name: "campaign not found",
			dto: campaign.DeleteCampaignDTO{
				ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
			},
			mock: func(dto campaign.DeleteCampaignDTO) {
				s.mockRepo.On("GetByID", campaign.GetCampaignDTO(dto)).Return(nil, campaign.ErrCampaignNotFound)
				s.mockLogger.On("Debug", "Campaign not found for deletion", "id", dto.ID).Return()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.mockRepo.ExpectedCalls = nil
			s.mockRepo.Calls = nil

			tt.mock(tt.dto)
			err := s.service.DeleteCampaign(tt.dto)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
			s.mockRepo.AssertExpectations(s.T())
		})
	}
}
