package campaign

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// CreateCampaignDTO defines the data structure for creating a campaign
type CreateCampaignDTO struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	PostalCode  string `json:"postal_code" validate:"required"`
	Template    string `json:"template" validate:"required"`
	OwnerID     string `json:"owner_id" validate:"required,uuid4"`
}

type Params struct {
	ID          int
	Name        string
	Description string
	PostalCode  string
	Template    string
	OwnerID     string
}

type ServiceInterface interface {
	ComposeEmail(mp Representative, c *Campaign, userData map[string]string) string
	CreateCampaign(c *CreateCampaignDTO) error
	DeleteCampaign(id int) error
	FetchCampaign(id int) (*Campaign, error)
	GetAllCampaigns() ([]Campaign, error)
	GetCampaignByID(id int) (*Campaign, error)
	UpdateCampaign(c *Campaign) error
}

type Service struct {
	repo     RepositoryInterface
	validate *validator.Validate
}

// NewService creates a new campaign service
func NewService(repo RepositoryInterface) (*Service, error) {
	validate := validator.New()
	service := &Service{
		repo:     repo,
		validate: validate,
	}
	return service, nil
}

func (s *Service) CreateCampaign(dto *CreateCampaignDTO) error {
	// Validate DTO fields here
	err := s.validate.Struct(dto)
	if err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}

	campaign := &Campaign{
		Name:        dto.Name,
		Description: dto.Description,
		PostalCode:  dto.PostalCode,
		Template:    dto.Template,
		OwnerID:     dto.OwnerID,
	}

	return s.repo.Create(campaign)
}

func (s *Service) UpdateCampaign(c *Campaign) error {
	return s.repo.Update(c)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCampaignNotFound
		}
		return nil, fmt.Errorf("failed to fetch campaign: %w", err)
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
