package campaign

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jonesrussell/loggo"
)

// NewService creates a new campaign service
func NewService(repo RepositoryInterface, validate *validator.Validate) ServiceInterface {
	return &Service{
		repo:     repo,
		validate: validate,
	}
}

// ServiceInterface defines the methods of the campaign service
type ServiceInterface interface {
	CreateCampaign(dto *CreateCampaignDTO) (*Campaign, error)
	UpdateCampaign(dto *UpdateCampaignDTO) error
	GetCampaignByID(params GetCampaignParams) (*Campaign, error)
	GetCampaigns() ([]Campaign, error)
	DeleteCampaign(params DeleteCampaignDTO) error
	FetchCampaign(params GetCampaignParams) (*Campaign, error)
	ComposeEmail(params ComposeEmailParams) string
	Error(message string, err error, params ...interface{})
	Info(message string, params ...interface{})
	Warn(message string, params ...interface{})
}

// Service implements the campaign service
type Service struct {
	repo     RepositoryInterface
	validate *validator.Validate
	logger   loggo.LoggerInterface
}

// Ensure Service implements ServiceInterface
var _ ServiceInterface = &Service{}

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

// GetCampaignParams defines parameters for fetching a campaign
type GetCampaignParams struct {
	ID int
}

// GetCampaignByID retrieves a campaign by ID
func (s *Service) GetCampaignByID(params GetCampaignParams) (*Campaign, error) {
	campaign, err := s.repo.GetByID(GetCampaignDTO{ID: params.ID})
	if err != nil {
		if errors.Is(err, ErrCampaignNotFound) {
			return nil, err // Pass through standard errors
		}
		return nil, fmt.Errorf("failed to get campaign with ID %d: %w", params.ID, err)
	}
	return campaign, nil
}

// GetCampaigns retrieves all campaigns
func (s *Service) GetCampaigns() ([]Campaign, error) {
	campaigns, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get campaigns: %w", err)
	}
	return campaigns, nil
}

// DeleteCampaign deletes a campaign by ID
func (s *Service) DeleteCampaign(params DeleteCampaignDTO) error {
	// Convert DeleteCampaignDTO to GetCampaignDTO
	getCampaignDTO := GetCampaignDTO(params)

	campaign, err := s.repo.GetByID(getCampaignDTO)
	if err != nil {
		return fmt.Errorf("failed to get campaign for deletion: %w", err)
	}
	if campaign == nil {
		return fmt.Errorf("campaign not found")
	}
	err = s.repo.Delete(params)
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

// ComposeEmailParams defines parameters for composing an email
type ComposeEmailParams struct {
	MP       Representative
	Campaign *Campaign
	UserData map[string]string
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

// Error implements the ServiceInterface for logging errors
func (s *Service) Error(message string, err error, params ...interface{}) {
	s.logger.Error(message, err, params...)
}

// Info implements the ServiceInterface for logging information
func (s *Service) Info(message string, params ...interface{}) {
	s.logger.Info(message, params...)
}

// Warn implements the ServiceInterface for logging warnings
func (s *Service) Warn(message string, params ...interface{}) {
	s.logger.Warn(message, params...)
}
