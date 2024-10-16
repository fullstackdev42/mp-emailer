package campaign

import (
	"fmt"
	"strings"
	"time"
)

// ServiceInterface defines the interface for the campaign service
type ServiceInterface interface {
	ComposeEmail(mp Representative, campaign *Campaign, userData map[string]string) string
	CreateCampaign(campaign *Campaign) error
	DeleteCampaign(id int) error
	FetchCampaign(id int) (*Campaign, error)
	GetAllCampaigns() ([]Campaign, error)
	GetCampaignByID(id int) (*Campaign, error)
	UpdateCampaign(campaign *Campaign) error
}

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) ServiceInterface {
	return &Service{repo: repo}
}

func (s *Service) GetCampaignByID(id int) (*Campaign, error) {
	return s.repo.GetByID(id)
}

func (s *Service) GetAllCampaigns() ([]Campaign, error) {
	campaigns, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get campaigns: %w", err)
	}
	return campaigns, nil
}

func (s *Service) CreateCampaign(campaign *Campaign) error {
	if campaign.Name == "" {
		return fmt.Errorf("campaign name cannot be empty")
	}
	if len(campaign.Template) < 10 {
		return fmt.Errorf("campaign template must be at least 10 characters long")
	}
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()
	err := s.repo.Create(campaign)
	if err != nil {
		return fmt.Errorf("failed to create campaign: %w", err)
	}
	return nil
}

func (s *Service) DeleteCampaign(id int) error {
	campaign, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get campaign for deletion: %w", err)
	}
	if campaign == nil {
		return fmt.Errorf("campaign not found")
	}
	err = s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete campaign: %w", err)
	}
	return nil
}

func (s *Service) UpdateCampaign(campaign *Campaign) error {
	if campaign.Name == "" {
		return fmt.Errorf("campaign name cannot be empty")
	}
	if len(campaign.Template) < 10 {
		return fmt.Errorf("campaign template must be at least 10 characters long")
	}
	existingCampaign, err := s.repo.GetByID(campaign.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing campaign: %w", err)
	}
	if existingCampaign == nil {
		return fmt.Errorf("campaign not found")
	}
	campaign.UpdatedAt = time.Now()
	err = s.repo.Update(campaign)
	if err != nil {
		return fmt.Errorf("failed to update campaign: %w", err)
	}
	return nil
}

func (s *Service) FetchCampaign(id int) (*Campaign, error) {
	campaign, err := s.repo.GetCampaign(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, ErrCampaignNotFound
		}
		return nil, err
	}
	return campaign, nil
}

func (s *Service) ComposeEmail(mp Representative, campaign *Campaign, userData map[string]string) string {
	emailTemplate := campaign.Template
	for key, value := range userData {
		placeholder := fmt.Sprintf("{{%s}}", key)
		emailTemplate = strings.ReplaceAll(emailTemplate, placeholder, value)
	}
	emailTemplate = strings.ReplaceAll(emailTemplate, "{{MP's Name}}", mp.Name)
	emailTemplate = strings.ReplaceAll(emailTemplate, "{{MPEmail}}", mp.Email)
	emailTemplate = strings.ReplaceAll(emailTemplate, "{{Date}}", time.Now().Format("2006-01-02"))
	return emailTemplate
}
