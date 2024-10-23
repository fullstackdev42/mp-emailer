package campaign

import (
	"fmt"
	"strings"
	"time"
)

type Params struct {
	ID          int
	Name        string
	Description string
	PostalCode  string
	Template    string
	OwnerID     string // Change this from int to string
}

type ServiceInterface interface {
	ComposeEmail(mp Representative, c *Campaign, userData map[string]string) string
	CreateCampaign(c *Campaign) error
	DeleteCampaign(id int) error
	FetchCampaign(id int) (*Campaign, error)
	GetAllCampaigns() ([]Campaign, error)
	GetCampaignByID(id int) (*Campaign, error)
	UpdateCampaign(c *Campaign) error
}

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) ServiceInterface {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateCampaign(c *Campaign) error {
	return s.repo.Create(&Campaign{
		Name:        c.Name,
		Description: c.Description,
		PostalCode:  c.PostalCode,
		Template:    c.Template,
		OwnerID:     c.OwnerID, // This should now be a string (UUID)
	})
}

func (s *Service) UpdateCampaign(c *Campaign) error {
	return s.repo.Update(&Campaign{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		PostalCode:  c.PostalCode,
		Template:    c.Template,
		OwnerID:     c.OwnerID, // This should now be a string (UUID)
	})
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
