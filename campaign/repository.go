package campaign

import (
	"fmt"

	"gorm.io/gorm"
)

// RepositoryInterface defines the methods that a campaign repository must implement
type RepositoryInterface interface {
	Create(dto *CreateCampaignDTO) (*Campaign, error)
	GetAll() ([]Campaign, error)
	Update(dto *UpdateCampaignDTO) error
	Delete(dto DeleteCampaignDTO) error
	GetByID(dto GetCampaignDTO) (*Campaign, error)
}

// Ensure that Repository implements RepositoryInterface
var _ RepositoryInterface = (*Repository)(nil)

// Repository handles CRUD operations for campaigns
type Repository struct {
	db *gorm.DB
}

// Create creates a new campaign in the database
func (r *Repository) Create(dto *CreateCampaignDTO) (*Campaign, error) {
	campaign := &Campaign{
		Name:        dto.Name,
		Description: dto.Description,
		Template:    dto.Template,
		OwnerID:     dto.OwnerID,
	}

	if err := r.db.Create(campaign).Error; err != nil {
		return nil, fmt.Errorf("error creating campaign: %w", err)
	}
	return campaign, nil
}

func (r *Repository) GetAll() ([]Campaign, error) {
	var campaigns []Campaign
	if err := r.db.Find(&campaigns).Error; err != nil {
		return nil, fmt.Errorf("error querying campaigns: %w", err)
	}
	return campaigns, nil
}

func (r *Repository) Update(dto *UpdateCampaignDTO) error {
	result := r.db.Model(&Campaign{}).Where("id = ?", dto.ID).
		Updates(map[string]interface{}{
			"name":        dto.Name,
			"description": dto.Description,
			"template":    dto.Template,
		})

	if result.Error != nil {
		return fmt.Errorf("error updating campaign: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrCampaignNotFound
	}
	return nil
}

func (r *Repository) Delete(dto DeleteCampaignDTO) error {
	result := r.db.Delete(&Campaign{}, dto.ID)
	if result.Error != nil {
		return fmt.Errorf("error deleting campaign: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrCampaignNotFound
	}
	return nil
}

func (r *Repository) GetByID(dto GetCampaignDTO) (*Campaign, error) {
	var campaign Campaign
	if err := r.db.First(&campaign, dto.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrCampaignNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}
	return &campaign, nil
}
