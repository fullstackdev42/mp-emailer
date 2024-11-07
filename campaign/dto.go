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
	ID          uuid.UUID `validate:"required"`
	Name        string    `validate:"required,min=3"`
	Description string    `validate:"required"`
	Template    string    `validate:"required"`
}

// GetCampaignDTO represents the data structure for getting a campaign
type GetCampaignDTO struct {
	ID uuid.UUID `validate:"required"`
}

// ComposeEmailDTO represents the data structure for composing an email
type ComposeEmailDTO struct {
	MP       Representative    `validate:"required"`
	Campaign *Campaign         `validate:"required"`
	UserData map[string]string `validate:"required"`
}

// DeleteCampaignDTO represents the data structure for deleting a campaign
type DeleteCampaignDTO struct {
	ID uuid.UUID `validate:"required"`
}
