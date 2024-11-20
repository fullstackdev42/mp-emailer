package campaign

import (
	"context"
	"fmt"

	"github.com/jonesrussell/mp-emailer/database/core"
)

// RepositoryInterface defines the contract for campaign repository operations
type RepositoryInterface interface {
	Create(ctx context.Context, dto *CreateCampaignDTO) (*Campaign, error)
	GetAll(ctx context.Context) ([]Campaign, error)
	Update(ctx context.Context, dto *UpdateCampaignDTO) error
	Delete(ctx context.Context, dto DeleteCampaignDTO) error
	GetByID(ctx context.Context, dto GetCampaignDTO) (*Campaign, error)
}

// Repository implements the RepositoryInterface
type Repository struct {
	db core.Interface
}

// NewRepository creates a new instance of Repository
func NewRepository(params RepositoryParams) RepositoryInterface {
	return &Repository{db: params.DB}
}

// Create creates a new campaign in the database
func (r *Repository) Create(ctx context.Context, dto *CreateCampaignDTO) (*Campaign, error) {
	campaign := &Campaign{
		Name:        dto.Name,
		Description: dto.Description,
		Template:    dto.Template,
		OwnerID:     dto.OwnerID,
	}

	if err := r.db.Create(ctx, campaign); err != nil {
		return nil, fmt.Errorf("error creating campaign: %w", err)
	}
	return campaign, nil
}

// GetAll retrieves all campaigns from the database
func (r *Repository) GetAll(ctx context.Context) ([]Campaign, error) {
	var campaigns []Campaign
	result := r.db.Query(ctx, "SELECT * FROM campaigns")
	if err := result.Scan(&campaigns).Error(); err != nil {
		return nil, fmt.Errorf("error querying campaigns: %w", err)
	}
	return campaigns, nil
}

// Update updates an existing campaign in the database
func (r *Repository) Update(ctx context.Context, dto *UpdateCampaignDTO) error {
	exists, err := r.db.Exists(ctx, &Campaign{}, "id = ?", dto.ID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrCampaignNotFound
	}

	err = r.db.Exec(ctx, "UPDATE campaigns SET name = ?, description = ?, template = ? WHERE id = ?", dto.Name, dto.Description, dto.Template, dto.ID)
	if err != nil {
		return fmt.Errorf("error updating campaign: %w", err)
	}
	return nil
}

// Delete removes a campaign from the database
func (r *Repository) Delete(ctx context.Context, dto DeleteCampaignDTO) error {
	exists, err := r.db.Exists(ctx, &Campaign{}, "id = ?", dto.ID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrCampaignNotFound
	}

	err = r.db.Exec(ctx, "DELETE FROM campaigns WHERE id = ?", dto.ID)
	if err != nil {
		return fmt.Errorf("error deleting campaign: %w", err)
	}
	return nil
}

// GetByID retrieves a campaign by its ID
func (r *Repository) GetByID(ctx context.Context, dto GetCampaignDTO) (*Campaign, error) {
	var campaign Campaign
	err := r.db.FindOne(ctx, &campaign, "id = ?", dto.ID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving campaign: %w", err)
	}
	return &campaign, nil
}
