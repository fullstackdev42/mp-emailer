package campaign

import "html/template"

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
	Name        string `form:"name"`
	Description string `form:"description"`
	Template    string `form:"template"`
	OwnerID     string // This will be set from the session
}

// EditParams defines the parameters for editing a campaign
type EditParams struct {
	ID       int
	Name     string
	Template string
}

// SendCampaignParams defines the parameters for sending a campaign
type SendCampaignParams struct {
	ID         int    `param:"id"`
	PostalCode string `form:"postal_code"`
}
