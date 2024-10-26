package campaign

// CreateCampaignDTO represents the data structure for creating a new campaign
type CreateCampaignDTO struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Template    string `json:"template" validate:"required"`
	OwnerID     string `json:"owner_id" validate:"required,uuid4"`
}

// UpdateCampaignDTO represents the data structure for updating an existing campaign
type UpdateCampaignDTO struct {
	ID          int    `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Template    string `json:"template" validate:"required"`
}

// GetCampaignDTO represents the data structure for getting a campaign
type GetCampaignDTO struct {
	ID int `json:"id" validate:"required"`
}

// ComposeEmailDTO represents the data structure for composing an email
type ComposeEmailDTO struct {
	MP       Representative    `json:"mp" validate:"required"`
	Campaign *Campaign         `json:"campaign" validate:"required"`
	UserData map[string]string `json:"user_data" validate:"required"`
}

// DeleteCampaignDTO represents the data structure for deleting a campaign
type DeleteCampaignDTO struct {
	ID int `json:"id" validate:"required"`
}
