package campaign

import (
	"html/template"

	"github.com/google/uuid"
)

// TemplateData provides a consistent structure for all template rendering
type TemplateData struct {
	Campaign        *Campaign
	Campaigns       []Campaign
	Email           string
	Content         template.HTML
	Error           error
	Representatives []Representative
}

// CreateCampaignParams defines the parameters for creating a campaign
type CreateCampaignParams struct {
	Name        string    `form:"name"`
	Description string    `form:"description"`
	Template    string    `form:"template"`
	OwnerID     uuid.UUID `param:"owner_id"`
}

// EditParams defines the parameters for editing a campaign
type EditParams struct {
	ID          uuid.UUID `param:"id"`
	Name        string    `param:"name"`
	Description string    `param:"description"`
	Template    string    `param:"template"`
}

// SendCampaignParams defines the parameters for sending a campaign
type SendCampaignParams struct {
	ID         uuid.UUID `param:"id"`
	PostalCode string    `param:"postal_code"`
}

// GetCampaignParams represents parameters for fetching a campaign
type GetCampaignParams struct {
	ID uuid.UUID `param:"id"`
}
