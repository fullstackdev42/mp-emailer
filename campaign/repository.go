package campaign

import (
	"fmt"

	"github.com/fullstackdev42/mp-emailer/database/core"
)

// RepositoryInterface defines the contract for campaign repository operations
type RepositoryInterface interface {
	Create(dto *CreateCampaignDTO) (*Campaign, error)
	GetAll() ([]Campaign, error)
	Update(dto *UpdateCampaignDTO) error
	Delete(dto DeleteCampaignDTO) error
	GetByID(dto GetCampaignDTO) (*Campaign, error)
}

// Repository implements the RepositoryInterface
type Repository struct {
	db core.Interface
}

// NewRepository creates a new instance of Repository
func NewRepository(db core.Interface) RepositoryInterface {
	return &Repository{db: db}
}

// Create creates a new campaign in the database
func (r *Repository) Create(dto *CreateCampaignDTO) (*Campaign, error) {
	campaign := &Campaign{
		Name:        dto.Name,
		Description: dto.Description,
		Template:    dto.Template,
		OwnerID:     dto.OwnerID,
	}

	if err := r.db.Create(campaign); err != nil {
		return nil, fmt.Errorf("error creating campaign: %w", err)
	}
	return campaign, nil
}

// GetAll retrieves all campaigns from the database
func (r *Repository) GetAll() ([]Campaign, error) {
	var campaigns []Campaign
	result := r.db.Query("SELECT * FROM campaigns")
	if err := result.Scan(&campaigns).Error(); err != nil {
		return nil, fmt.Errorf("error querying campaigns: %w", err)
	}
	return campaigns, nil
}

// Update updates an existing campaign in the database
func (r *Repository) Update(dto *UpdateCampaignDTO) error {
	exists, err := r.db.Exists(&Campaign{}, "id = ?", dto.ID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrCampaignNotFound
	}

	err = r.db.Exec("UPDATE campaigns SET name = ?, description = ?, template = ? WHERE id = ?", dto.Name, dto.Description, dto.Template, dto.ID)
	if err != nil {
		return fmt.Errorf("error updating campaign: %w", err)
	}
	return nil
}

// Delete removes a campaign from the database
func (r *Repository) Delete(dto DeleteCampaignDTO) error {
	exists, err := r.db.Exists(&Campaign{}, "id = ?", dto.ID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrCampaignNotFound
	}

	err = r.db.Exec("DELETE FROM campaigns WHERE id = ?", dto.ID)
	if err != nil {
		return fmt.Errorf("error deleting campaign: %w", err)
	}
	return nil
}

// GetByID retrieves a campaign by its ID
func (r *Repository) GetByID(dto GetCampaignDTO) (*Campaign, error) {
	var campaign Campaign
	err := r.db.FindOne(&campaign, "id = ?", dto.ID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving campaign: %w", err)
	}
	return &campaign, nil
}
