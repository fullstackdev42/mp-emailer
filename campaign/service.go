package campaign

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jonesrussell/loggo"
)

// NewService creates a new campaign service
func NewService(repo RepositoryInterface, validate *validator.Validate, logger loggo.LoggerInterface) ServiceInterface {
	return &Service{
		repo:     repo,
		validate: validate,
		logger:   logger,
	}
}

// ServiceInterface defines the methods of the campaign service
type ServiceInterface interface {
	CreateCampaign(ctx context.Context, dto *CreateCampaignDTO) (*Campaign, error)
	UpdateCampaign(ctx context.Context, dto *UpdateCampaignDTO) error
	GetCampaignByID(ctx context.Context, params GetCampaignParams) (*Campaign, error)
	GetCampaigns(ctx context.Context) ([]Campaign, error)
	DeleteCampaign(ctx context.Context, params DeleteCampaignDTO) error
	FetchCampaign(ctx context.Context, params GetCampaignParams) (*Campaign, error)
	ComposeEmail(ctx context.Context, params ComposeEmailParams) (string, error)
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
func (s *Service) CreateCampaign(ctx context.Context, dto *CreateCampaignDTO) (*Campaign, error) {
	if dto == nil {
		return nil, fmt.Errorf("campaign data is required")
	}

	if err := s.validate.Struct(dto); err != nil {
		s.logger.Debug("Invalid campaign data", "error", err)
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	campaign, err := s.repo.Create(ctx, dto)
	if err != nil {
		s.logger.Error("Failed to create campaign", err)
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	s.logger.Info("Campaign created successfully", "id", campaign.ID)
	return campaign, nil
}

// UpdateCampaign updates an existing campaign
func (s *Service) UpdateCampaign(ctx context.Context, dto *UpdateCampaignDTO) error {
	err := s.validate.Struct(dto)
	if err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}
	return s.repo.Update(ctx, dto)
}

// GetCampaignByID retrieves a campaign by ID
func (s *Service) GetCampaignByID(ctx context.Context, params GetCampaignParams) (*Campaign, error) {
	campaign, err := s.repo.GetByID(ctx, GetCampaignDTO{ID: params.ID})
	if err != nil {
		if errors.Is(err, ErrCampaignNotFound) {
			s.logger.Debug("Campaign not found", "id", params.ID)
			return nil, err
		}
		s.logger.Error("Failed to get campaign", err, "id", params.ID)
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}

	return campaign, nil
}

// GetCampaigns retrieves all campaigns
func (s *Service) GetCampaigns(ctx context.Context) ([]Campaign, error) {
	campaigns, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaigns: %w", err)
	}
	return campaigns, nil
}

// DeleteCampaign deletes a campaign by ID
func (s *Service) DeleteCampaign(ctx context.Context, params DeleteCampaignDTO) error {
	_, err := s.repo.GetByID(ctx, GetCampaignDTO(params))
	if err != nil {
		if errors.Is(err, ErrCampaignNotFound) {
			s.logger.Debug("Campaign not found for deletion", "id", params.ID)
			return err
		}
		s.logger.Error("Failed to fetch campaign for deletion", err, "id", params.ID)
		return fmt.Errorf("failed to fetch campaign: %w", err)
	}

	if err := s.repo.Delete(ctx, params); err != nil {
		s.logger.Error("Failed to delete campaign", err, "id", params.ID)
		return fmt.Errorf("failed to delete campaign: %w", err)
	}

	s.logger.Info("Campaign deleted successfully", "id", params.ID)
	return nil
}

// FetchCampaign retrieves a campaign by parameters
func (s *Service) FetchCampaign(ctx context.Context, params GetCampaignParams) (*Campaign, error) {
	campaign, err := s.repo.GetByID(ctx, GetCampaignDTO{ID: params.ID})
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
func (s *Service) ComposeEmail(ctx context.Context, params ComposeEmailParams) (string, error) {
	if params.Campaign == nil {
		return "", fmt.Errorf("campaign is required")
	}
	if params.Campaign.Template == "" {
		return "", fmt.Errorf("campaign template is required")
	}

	emailTemplate := params.Campaign.Template

	// Replace user data placeholders
	for key, value := range params.UserData {
		placeholder := fmt.Sprintf("{{%s}}", key)
		emailTemplate = strings.ReplaceAll(emailTemplate, placeholder, value)
	}

	// Replace standard placeholders
	replacements := map[string]string{
		"{{MP's Name}}": params.MP.Name,
		"{{MPEmail}}":   params.MP.Email,
		"{{Date}}":      time.Now().Format("2006-01-02"),
	}

	for placeholder, value := range replacements {
		emailTemplate = strings.ReplaceAll(emailTemplate, placeholder, value)
	}

	return emailTemplate, nil
}
