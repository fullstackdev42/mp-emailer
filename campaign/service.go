package campaign

import (
	"fmt"
	"strconv"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCampaignByID(id string) (*Campaign, error) {
	campaignID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign ID: %w", err)
	}
	return s.repo.GetByID(campaignID)
}

func (s *Service) GetAllCampaigns() ([]Campaign, error) {
	campaigns, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get campaigns: %w", err)
	}
	return campaigns, nil
}

func (s *Service) CreateCampaign(campaign *Campaign) error {
	// Add business logic, e.g., validation
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
	// Add business logic, e.g., check if user has permission to delete
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
	// Add business logic, e.g., validation
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
