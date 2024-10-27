package campaign

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// Parameter objects (for internal service use)
type GetCampaignParams struct {
	ID int
}

type ComposeEmailParams struct {
	MP       Representative
	Campaign *Campaign
	UserData map[string]string
}

type ServiceInterface interface {
	CreateCampaign(dto *CreateCampaignDTO) (*Campaign, error)
	UpdateCampaign(dto *UpdateCampaignDTO) error
	GetCampaignByID(params GetCampaignParams) (*Campaign, error)
	GetAllCampaigns() ([]Campaign, error)
	DeleteCampaign(params DeleteCampaignParams) error
	FetchCampaign(params GetCampaignParams) (*Campaign, error)
	ComposeEmail(params ComposeEmailParams) string
}

type Service struct {
	repo     RepositoryInterface
	validate *validator.Validate
}

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

func (s *Service) UpdateCampaign(dto *UpdateCampaignDTO) error {
	err := s.validate.Struct(dto)
	if err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}

	return s.repo.Update(dto)
}

func (s *Service) GetCampaignByID(params GetCampaignParams) (*Campaign, error) {
	return s.repo.GetByID(GetCampaignDTO{ID: params.ID})
}

func (s *Service) GetAllCampaigns() ([]Campaign, error) {
	campaigns, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get campaigns: %w", err)
	}
	return campaigns, nil
}

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
