package campaign

// CreateCampaignDTO represents the data structure for creating a new campaign
type CreateCampaignDTO struct {
	Name        string `json:"name" validate:"required"`           // Name of the campaign
	Description string `json:"description" validate:"required"`    // Description of the campaign
	Template    string `json:"template" validate:"required"`       // Email template for the campaign
	OwnerID     string `json:"owner_id" validate:"required,uuid4"` // ID of the campaign owner (UUID)
}

// UpdateCampaignDTO represents the data structure for updating an existing campaign
type UpdateCampaignDTO struct {
	ID          int    `json:"id" validate:"required"`          // Unique identifier for the campaign
	Name        string `json:"name" validate:"required"`        // Updated name of the campaign
	Description string `json:"description" validate:"required"` // Updated description of the campaign
	Template    string `json:"template" validate:"required"`    // Updated email template for the campaign
}

// GetCampaignDTO represents the data structure for getting a campaign
type GetCampaignDTO struct {
	ID int `json:"id" validate:"required"` // Unique identifier for the campaign
}

// ComposeEmailDTO represents the data structure for composing an email
type ComposeEmailDTO struct {
	MP       Representative    `json:"mp" validate:"required"`        // Representative details
	Campaign *Campaign         `json:"campaign" validate:"required"`  // Campaign details
	UserData map[string]string `json:"user_data" validate:"required"` // User data for email customization
}

// DeleteCampaignDTO represents the data structure for deleting a campaign
type DeleteCampaignDTO struct {
	ID int `json:"id" validate:"required"` // Unique identifier for the campaign
}
