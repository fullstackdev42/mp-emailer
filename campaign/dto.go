package campaign

import (
	"github.com/google/uuid"
)

// CreateCampaignDTO represents the data structure for creating a new campaign
type CreateCampaignDTO struct {
	Name        string    `validate:"required,min=3"`
	Description string    `validate:"required"`
	Template    string    `validate:"required"`
	OwnerID     uuid.UUID `validate:"required"`
}

// UpdateCampaignDTO represents the data structure for updating an existing campaign
type UpdateCampaignDTO struct {
	ID          uuid.UUID `param:"id"`
	Name        string    `param:"name"`
	Description string    `param:"description"`
	Template    string    `param:"template"`
}

// GetCampaignDTO represents the data structure for getting a campaign
type GetCampaignDTO struct {
	ID uuid.UUID `param:"id"`
}

// ComposeEmailDTO represents the data structure for composing an email
type ComposeEmailDTO struct {
	MP       Representative    `param:"mp" validate:"required"`        // Representative details
	Campaign *Campaign         `param:"campaign" validate:"required"`  // Campaign details
	UserData map[string]string `param:"user_data" validate:"required"` // User data for email customization
}

// DeleteCampaignDTO represents the data structure for deleting a campaign
type DeleteCampaignDTO struct {
	ID uuid.UUID `param:"id"`
}
