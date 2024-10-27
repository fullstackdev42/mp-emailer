package campaign

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// GetCampaignParams defines parameters for fetching a campaign
type GetCampaignParams struct {
	ID int
}

// ComposeEmailParams defines parameters for composing an email
type ComposeEmailParams struct {
	MP       Representative
	Campaign *Campaign
	UserData map[string]string
}

// ServiceInterface defines the methods of the campaign service
type ServiceInterface interface {
	CreateCampaign(dto *CreateCampaignDTO) (*Campaign, error)
	UpdateCampaign(dto *UpdateCampaignDTO) error
	GetCampaignByID(params GetCampaignParams) (*Campaign, error)
	GetAllCampaigns() ([]Campaign, error)
	DeleteCampaign(params DeleteCampaignParams) error
	FetchCampaign(params GetCampaignParams) (*Campaign, error)
	ComposeEmail(params ComposeEmailParams) string
}

// Error implements shared.ServiceInterface.
func (s *Service) Error(message string, err error) {
	panic("unimplemented")
}

// Info implements shared.ServiceInterface.
func (s *Service) Info(message string) {
	panic("unimplemented")
}

// Warn implements shared.ServiceInterface.
func (s *Service) Warn(message string, err error) {
	panic("unimplemented")
}

// Service implements the campaign service
type Service struct {
	repo     RepositoryInterface
	validate *validator.Validate
}

// CreateCampaign creates a new campaign
func (s *Service) CreateCampaign(dto *CreateCampaignDTO) (*Campaign, error) {
	err := s.validate.Struct(dto)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	campaign, err := s.repo.Create(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	return campaign, nil
}

// UpdateCampaign updates an existing campaign
func (s *Service) UpdateCampaign(dto *UpdateCampaignDTO) error {
	err := s.validate.Struct(dto)
	if err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}

	return s.repo.Update(dto)
}

// GetCampaignByID retrieves a campaign by ID
func (s *Service) GetCampaignByID(params GetCampaignParams) (*Campaign, error) {
	return s.repo.GetByID(GetCampaignDTO{ID: params.ID})
}

// GetAllCampaigns retrieves all campaigns
func (s *Service) GetAllCampaigns() ([]Campaign, error) {
	campaigns, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get campaigns: %w", err)
	}
	return campaigns, nil
}

// DeleteCampaign deletes a campaign by ID
func (s *Service) DeleteCampaign(params DeleteCampaignParams) error {
	campaign, err := s.repo.GetByID(GetCampaignDTO{ID: params.ID})
	if err != nil {
		return fmt.Errorf("failed to get campaign for deletion: %w", err)
	}
	if campaign == nil {
		return fmt.Errorf("campaign not found")
	}

	err = s.repo.Delete(DeleteCampaignDTO{ID: params.ID})
	if err != nil {
		return fmt.Errorf("failed to delete campaign: %w", err)
	}
	return nil
}

// FetchCampaign retrieves a campaign by parameters
func (s *Service) FetchCampaign(params GetCampaignParams) (*Campaign, error) {
	campaign, err := s.repo.GetCampaign(GetCampaignDTO{ID: params.ID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCampaignNotFound
		}
		return nil, fmt.Errorf("failed to fetch campaign: %w", err)
	}
	return campaign, nil
}

// ComposeEmail composes an email using campaign data and user data
func (s *Service) ComposeEmail(params ComposeEmailParams) string {
	emailTemplate := params.Campaign.Template
	for key, value := range params.UserData {
		placeholder := fmt.Sprintf("{{%s}}", key)
		emailTemplate = strings.ReplaceAll(emailTemplate, placeholder, value)
	}
	emailTemplate = strings.ReplaceAll(emailTemplate, "{{MP's Name}}", params.MP.Name)
	emailTemplate = strings.ReplaceAll(emailTemplate, "{{MPEmail}}", params.MP.Email)
	emailTemplate = strings.ReplaceAll(emailTemplate, "{{Date}}", time.Now().Format("2006-01-02"))
	return emailTemplate
}
