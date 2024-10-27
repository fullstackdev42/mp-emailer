package campaign

import (
	"github.com/jonesrussell/loggo"
)

// LoggingServiceDecorator adds logging functionality to the ServiceInterface
type LoggingServiceDecorator struct {
	service ServiceInterface
	logger  loggo.LoggerInterface
}

// NewLoggingServiceDecorator creates a new instance of LoggingServiceDecorator
func NewLoggingServiceDecorator(service ServiceInterface, logger loggo.LoggerInterface) *LoggingServiceDecorator {
	return &LoggingServiceDecorator{
		service: service,
		logger:  logger,
	}
}

// CreateCampaign logs and delegates the creation of a new campaign
func (d *LoggingServiceDecorator) CreateCampaign(dto *CreateCampaignDTO) (*Campaign, error) {
	d.logger.Info("Creating new campaign", "name", dto.Name)
	campaign, err := d.service.CreateCampaign(dto)
	if err != nil {
		d.logger.Error("Failed to create campaign", err)
	}
	return campaign, err
}

// UpdateCampaign logs and delegates the update of an existing campaign
func (d *LoggingServiceDecorator) UpdateCampaign(dto *UpdateCampaignDTO) error {
	d.logger.Info("Updating campaign", "id", dto.ID)
	err := d.service.UpdateCampaign(dto)
	if err != nil {
		d.logger.Error("Failed to update campaign", err, "id", dto.ID)
	}
	return err
}

// GetCampaignByID logs and delegates fetching a campaign by ID
func (d *LoggingServiceDecorator) GetCampaignByID(params GetCampaignParams) (*Campaign, error) {
	d.logger.Info("Fetching campaign", "id", params.ID)
	campaign, err := d.service.GetCampaignByID(params)
	if err != nil {
		d.logger.Error("Failed to fetch campaign", err, "id", params.ID)
	}
	return campaign, err
}

// GetAllCampaigns logs and delegates fetching all campaigns
func (d *LoggingServiceDecorator) GetAllCampaigns() ([]Campaign, error) {
	d.logger.Info("Fetching all campaigns")
	campaigns, err := d.service.GetAllCampaigns()
	if err != nil {
		d.logger.Error("Failed to fetch all campaigns", err)
	}
	return campaigns, err
}

// DeleteCampaign logs and delegates deleting a campaign by ID
func (d *LoggingServiceDecorator) DeleteCampaign(params DeleteCampaignParams) error {
	d.logger.Info("Deleting campaign", "id", params.ID)
	err := d.service.DeleteCampaign(params)
	if err != nil {
		d.logger.Error("Failed to delete campaign", err, "id", params.ID)
	}
	return err
}

// FetchCampaign logs and delegates fetching a campaign by parameters
func (d *LoggingServiceDecorator) FetchCampaign(params GetCampaignParams) (*Campaign, error) {
	d.logger.Info("Fetching campaign", "id", params.ID)
	campaign, err := d.service.FetchCampaign(params)
	if err != nil {
		d.logger.Error("Failed to fetch campaign", err, "id", params.ID)
	}
	return campaign, err
}

// ComposeEmail logs and delegates composing an email for a campaign
func (d *LoggingServiceDecorator) ComposeEmail(params ComposeEmailParams) string {
	d.logger.Info("Composing email for campaign", "campaignID", params.Campaign.ID)
	email := d.service.ComposeEmail(params)
	d.logger.Debug("Email composed", "length", len(email))
	return email
}
