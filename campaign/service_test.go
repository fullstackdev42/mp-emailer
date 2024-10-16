package campaign

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetCampaignByID(t *testing.T) {
	mockRepo := NewMockRepositoryInterface(t)
	service := NewService(mockRepo)

	campaign := &Campaign{ID: 1, Name: "Test Campaign"}
	mockRepo.EXPECT().GetByID(1).Return(campaign, nil)

	result, err := service.GetCampaignByID(1)

	assert.NoError(t, err)
	assert.Equal(t, campaign, result)
}

func TestService_GetAllCampaigns(t *testing.T) {
	mockRepo := NewMockRepositoryInterface(t)
	service := NewService(mockRepo)

	campaigns := []Campaign{{ID: 1, Name: "Campaign 1"}, {ID: 2, Name: "Campaign 2"}}
	mockRepo.EXPECT().GetAll().Return(campaigns, nil)

	result, err := service.GetAllCampaigns()

	assert.NoError(t, err)
	assert.Equal(t, campaigns, result)
}

func TestService_CreateCampaign(t *testing.T) {
	mockRepo := NewMockRepositoryInterface(t)
	service := NewService(mockRepo)

	campaign := &Campaign{Name: "New Campaign", Template: "Test Template"}
	mockRepo.EXPECT().Create(mock.AnythingOfType("*campaign.Campaign")).Return(nil)

	err := service.CreateCampaign(campaign)

	assert.NoError(t, err)
	assert.NotZero(t, campaign.CreatedAt)
	assert.NotZero(t, campaign.UpdatedAt)
}

func TestService_CreateCampaign_ValidationError(t *testing.T) {
	mockRepo := NewMockRepositoryInterface(t)
	service := NewService(mockRepo)

	campaign := &Campaign{Name: "", Template: "Short"}

	err := service.CreateCampaign(campaign)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "campaign name cannot be empty")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestService_DeleteCampaign(t *testing.T) {
	mockRepo := NewMockRepositoryInterface(t)
	service := NewService(mockRepo)

	mockRepo.EXPECT().GetByID(1).Return(&Campaign{ID: 1}, nil)
	mockRepo.EXPECT().Delete(1).Return(nil)

	err := service.DeleteCampaign(1)

	assert.NoError(t, err)
}

func TestService_DeleteCampaign_NotFound(t *testing.T) {
	mockRepo := NewMockRepositoryInterface(t)
	service := NewService(mockRepo)

	mockRepo.EXPECT().GetByID(1).Return((*Campaign)(nil), nil)

	err := service.DeleteCampaign(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "campaign not found")
}

func TestService_UpdateCampaign(t *testing.T) {
	mockRepo := NewMockRepositoryInterface(t)
	service := NewService(mockRepo)

	campaign := &Campaign{ID: 1, Name: "Updated Campaign", Template: "Updated Template"}
	mockRepo.EXPECT().GetByID(1).Return(campaign, nil)
	mockRepo.EXPECT().Update(mock.AnythingOfType("*campaign.Campaign")).Return(nil)

	err := service.UpdateCampaign(campaign)

	assert.NoError(t, err)
	assert.NotZero(t, campaign.UpdatedAt)
}

func TestService_UpdateCampaign_ValidationError(t *testing.T) {
	mockRepo := NewMockRepositoryInterface(t)
	service := NewService(mockRepo)

	campaign := &Campaign{ID: 1, Name: "", Template: "Short"}

	err := service.UpdateCampaign(campaign)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "campaign name cannot be empty")
	mockRepo.AssertNotCalled(t, "GetByID")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestService_FetchCampaign(t *testing.T) {
	mockRepo := NewMockRepositoryInterface(t)
	service := NewService(mockRepo)

	campaign := &Campaign{ID: 1, Name: "Test Campaign"}
	mockRepo.EXPECT().GetCampaign(1).Return(campaign, nil)

	result, err := service.FetchCampaign(1)

	assert.NoError(t, err)
	assert.Equal(t, campaign, result)
}

func TestService_FetchCampaign_NotFound(t *testing.T) {
	mockRepo := NewMockRepositoryInterface(t)
	service := NewService(mockRepo)

	mockRepo.EXPECT().GetCampaign(1).Return((*Campaign)(nil), errors.New("sql: no rows in result set"))

	result, err := service.FetchCampaign(1)

	assert.Error(t, err)
	assert.Equal(t, ErrCampaignNotFound, err)
	assert.Nil(t, result)
}

func TestService_ComposeEmail(t *testing.T) {
	service := NewService(nil) // Repository not needed for this test

	mp := Representative{Name: "John Doe", Email: "john@example.com"}
	campaign := &Campaign{Template: "Hello {{MP's Name}}, This is a test email. Today is {{Date}}. {{First Name}} {{Last Name}}"}
	userData := map[string]string{
		"First Name": "Jane",
		"Last Name":  "Smith",
	}

	result := service.ComposeEmail(mp, campaign, userData)

	expectedDate := time.Now().Format("2006-01-02")
	expected := "Hello John Doe, This is a test email. Today is " + expectedDate + ". Jane Smith"
	assert.Equal(t, expected, result)
}
