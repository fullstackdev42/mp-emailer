package campaign_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/fullstackdev42/mp-emailer/campaign"
	mocksCampaign "github.com/fullstackdev42/mp-emailer/mocks/campaign"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTest() (*campaign.Service, *mocksCampaign.MockRepositoryInterface) {
	mockRepo := new(mocksCampaign.MockRepositoryInterface)
	validate := validator.New()

	// Register the UUID validator
	err := validate.RegisterValidation("uuid4", func(fl validator.FieldLevel) bool {
		// Simple UUID4 format check
		uuid := fl.Field().String()
		if len(uuid) != 36 {
			return false
		}
		// Check for UUID4 format: 8-4-4-4-12 characters
		parts := strings.Split(uuid, "-")
		if len(parts) != 5 {
			return false
		}
		lengths := []int{8, 4, 4, 4, 12}
		for i, part := range parts {
			if len(part) != lengths[i] {
				return false
			}
		}
		return true
	})

	if err != nil {
		panic(fmt.Sprintf("failed to register uuid4 validator: %v", err))
	}

	service := campaign.NewService(mockRepo, validate)
	return service.(*campaign.Service), mockRepo
}

func TestCreateCampaign(t *testing.T) {
	service, mockRepo := setupTest()

	tests := []struct {
		name    string
		input   *campaign.CreateCampaignDTO
		mock    func()
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
			mock: func() {
				expectedCampaign := &campaign.Campaign{
					Name:        "Test Campaign",
					Description: "Test Description",
					Template:    "Test Template",
					OwnerID:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Tokens:      []string{},
				}
				mockRepo.On("Create", &campaign.CreateCampaignDTO{
					Name:        "Test Campaign",
					Description: "Test Description",
					Template:    "Test Template",
					OwnerID:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				}).Return(expectedCampaign, nil)
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
		{
			name: "repository error",
			input: &campaign.CreateCampaignDTO{
				Name:        "Test Campaign",
				Description: "Test Description",
				Template:    "Test Template",
				OwnerID:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			},
			mock: func() {
				mockRepo.On("Create", mock.AnythingOfType("*campaign.CreateCampaignDTO")).
					Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "validation error",
			input: &campaign.CreateCampaignDTO{
				// Missing required fields
				Name: "Test Campaign",
			},
			mock:    func() {}, // Validation will fail before repository is called
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			tt.mock()
			got, err := service.CreateCampaign(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestComposeEmail(t *testing.T) {
	service, _ := setupTest()

	tests := []struct {
		name   string
		params campaign.ComposeEmailParams
		want   string
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
			want: "Dear John Doe, This is a test email. Your email is john@example.com. Date: " + time.Now().Format("2006-01-02"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.ComposeEmail(tt.params)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetCampaignByID(t *testing.T) {
	service, mockRepo := setupTest()

	testID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name    string
		params  campaign.GetCampaignParams
		mock    func(uuid.UUID)
		want    *campaign.Campaign
		wantErr bool
	}{
		{
			name: "successful retrieval",
			params: campaign.GetCampaignParams{
				ID: testID,
			},
			mock: func(id uuid.UUID) {
				mockRepo.On("GetByID", campaign.GetCampaignDTO{ID: id}).
					Return(&campaign.Campaign{
						BaseModel: shared.BaseModel{
							ID: id,
						},
						Name:        "Test Campaign",
						Description: "Test Description",
						Template:    "Test Template",
						OwnerID:     id,
						Tokens:      []string{},
					}, nil)
			},
			want: &campaign.Campaign{
				BaseModel: shared.BaseModel{
					ID: testID,
				},
				Name:        "Test Campaign",
				Description: "Test Description",
				Template:    "Test Template",
				OwnerID:     testID,
				Tokens:      []string{},
			},
			wantErr: false,
		},
		{
			name: "campaign not found",
			params: campaign.GetCampaignParams{
				ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			},
			mock: func(id uuid.UUID) {
				mockRepo.On("GetByID", campaign.GetCampaignDTO{ID: id}).
					Return(nil, campaign.ErrCampaignNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear previous mock calls
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			// Setup mock with the same UUID as params
			tt.mock(tt.params.ID)

			got, err := service.GetCampaignByID(tt.params)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteCampaign(t *testing.T) {
	service, mockRepo := setupTest()
	testID := uuid.New()

	tests := []struct {
		name    string
		params  campaign.DeleteCampaignDTO
		mock    func()
		wantErr bool
	}{
		{
			name:   "successful deletion",
			params: campaign.DeleteCampaignDTO{ID: testID},
			mock: func() {
				mockRepo.On("GetByID", campaign.GetCampaignDTO{ID: testID}).
					Return(&campaign.Campaign{}, nil)
				mockRepo.On("Delete", campaign.DeleteCampaignDTO{ID: testID}).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "campaign not found",
			params: campaign.DeleteCampaignDTO{ID: testID},
			mock: func() {
				mockRepo.On("GetByID", campaign.GetCampaignDTO{ID: testID}).
					Return(nil, campaign.ErrCampaignNotFound).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear previous mock calls and expectations
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			tt.mock()
			err := service.DeleteCampaign(tt.params)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
