package campaign

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jonesrussell/mp-emailer/database"
	"github.com/jonesrussell/mp-emailer/shared"
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
	db database.Database
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

	// Fetch the created campaign to ensure all fields are properly set
	var created Campaign
	if err := r.db.FindOne(ctx, &created, "id = ?", campaign.ID); err != nil {
		return nil, fmt.Errorf("error retrieving created campaign: %w", err)
	}

	return &created, nil
}

// GetAll retrieves all campaigns from the database
func (r *Repository) GetAll(ctx context.Context) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.FindAll(ctx, &campaigns, "1=1")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get all campaigns: %w", err)
	}
	// Return empty slice if no records found
	if len(campaigns) == 0 {
		return []Campaign{}, nil
	}
	return campaigns, nil
}

// Update updates an existing campaign in the database
func (r *Repository) Update(ctx context.Context, dto *UpdateCampaignDTO) error {
	campaign := &Campaign{
		BaseModel:   shared.BaseModel{ID: dto.ID},
		Name:        dto.Name,
		Description: dto.Description,
		Template:    dto.Template,
	}

	if err := r.db.Update(ctx, campaign); err != nil {
		return fmt.Errorf("error updating campaign: %w", err)
	}
	return nil
}

// Delete removes a campaign from the database
func (r *Repository) Delete(ctx context.Context, dto DeleteCampaignDTO) error {
	campaign := &Campaign{
		BaseModel: shared.BaseModel{ID: dto.ID},
	}
	if err := r.db.Delete(ctx, campaign); err != nil {
		return fmt.Errorf("error deleting campaign: %w", err)
	}
	return nil
}

// GetByID retrieves a campaign by its ID
func (r *Repository) GetByID(ctx context.Context, dto GetCampaignDTO) (*Campaign, error) {
	var campaign Campaign
	if err := r.db.FindOne(ctx, &campaign, "id = ?", dto.ID); err != nil {
		return nil, fmt.Errorf("error retrieving campaign: %w", err)
	}
	return &campaign, nil
}
