package campaign

import (
	"github.com/jonesrussell/loggo"
)

type LoggingServiceDecorator struct {
	service ServiceInterface
	logger  loggo.LoggerInterface
}

func NewLoggingServiceDecorator(service ServiceInterface, logger loggo.LoggerInterface) *LoggingServiceDecorator {
	return &LoggingServiceDecorator{
		service: service,
		logger:  logger,
	}
}

func (d *LoggingServiceDecorator) CreateCampaign(dto *CreateCampaignDTO) (*Campaign, error) {
	d.logger.Info("Creating new campaign", "name", dto.Name)
	campaign, err := d.service.CreateCampaign(dto)
	if err != nil {
		d.logger.Error("Failed to create campaign", err)
	}
	return campaign, err
}

func (d *LoggingServiceDecorator) UpdateCampaign(dto *UpdateCampaignDTO) error {
	d.logger.Info("Updating campaign", "id", dto.ID)
	err := d.service.UpdateCampaign(dto)
	if err != nil {
		d.logger.Error("Failed to update campaign", err)
	}
	return err
}

func (d *LoggingServiceDecorator) GetCampaignByID(params GetCampaignParams) (*Campaign, error) {
	d.logger.Info("Fetching campaign", "id", params.ID)
	campaign, err := d.service.GetCampaignByID(params)
	if err != nil {
		d.logger.Error("Failed to fetch campaign", err)
	}
	return campaign, err
}

func (d *LoggingServiceDecorator) GetAllCampaigns() ([]Campaign, error) {
	d.logger.Info("Fetching all campaigns")
	campaigns, err := d.service.GetAllCampaigns()
	if err != nil {
		d.logger.Error("Failed to fetch all campaigns", err)
	}
	return campaigns, err
}

func (d *LoggingServiceDecorator) DeleteCampaign(params DeleteCampaignParams) error {
	d.logger.Info("Deleting campaign", "id", params.ID)
	err := d.service.DeleteCampaign(params)
	if err != nil {
		d.logger.Error("Failed to delete campaign", err)
	}
	return err
}

func (d *LoggingServiceDecorator) FetchCampaign(params GetCampaignParams) (*Campaign, error) {
	d.logger.Info("Fetching campaign", "id", params.ID)
	campaign, err := d.service.FetchCampaign(params)
	if err != nil {
		d.logger.Error("Failed to fetch campaign", err)
	}
	return campaign, err
}

func (d *LoggingServiceDecorator) ComposeEmail(params ComposeEmailParams) string {
	d.logger.Info("Composing email for campaign", "campaignID", params.Campaign.ID)
	email := d.service.ComposeEmail(params)
	d.logger.Debug("Email composed", "length", len(email))
	return email
}
