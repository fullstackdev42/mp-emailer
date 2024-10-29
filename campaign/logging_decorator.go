package campaign

import (
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/jonesrussell/loggo"
)

// LoggingServiceDecorator adds logging functionality to the Campaign ServiceInterface
type LoggingServiceDecorator struct {
	shared.LoggingServiceDecorator
	service ServiceInterface
}

// NewLoggingServiceDecorator creates a new instance of LoggingServiceDecorator
func NewLoggingServiceDecorator(service ServiceInterface, logger loggo.LoggerInterface) *LoggingServiceDecorator {
	return &LoggingServiceDecorator{
		LoggingServiceDecorator: *shared.NewLoggingServiceDecorator(service, logger),
		service:                 service,
	}
}

func (d *LoggingServiceDecorator) CreateCampaign(dto *CreateCampaignDTO) (*Campaign, error) {
	d.Logger.Info("Creating new campaign", "name", dto.Name)
	campaign, err := d.service.CreateCampaign(dto)
	if err != nil {
		d.Logger.Error("Failed to create campaign", err)
	}
	return campaign, err
}

func (d *LoggingServiceDecorator) UpdateCampaign(dto *UpdateCampaignDTO) error {
	d.Logger.Info("Updating campaign", "id", dto.ID)
	err := d.service.UpdateCampaign(dto)
	if err != nil {
		d.Logger.Error("Failed to update campaign", err)
	}
	return err
}

func (d *LoggingServiceDecorator) GetCampaignByID(params GetCampaignParams) (*Campaign, error) {
	d.Logger.Info("Fetching campaign", "id", params.ID)
	campaign, err := d.service.GetCampaignByID(params)
	if err != nil {
		d.Logger.Error("Failed to fetch campaign", err)
	}
	return campaign, err
}

func (d *LoggingServiceDecorator) GetAllCampaigns() ([]Campaign, error) {
	d.Logger.Info("Fetching all campaigns")
	campaigns, err := d.service.GetAllCampaigns()
	if err != nil {
		d.Logger.Error("Failed to fetch all campaigns", err)
	}
	return campaigns, err
}

func (d *LoggingServiceDecorator) DeleteCampaign(params DeleteCampaignParams) error {
	d.Logger.Info("Deleting campaign", "id", params.ID)
	err := d.service.DeleteCampaign(params)
	if err != nil {
		d.Logger.Error("Failed to delete campaign", err)
	}
	return err
}

func (d *LoggingServiceDecorator) FetchCampaign(params GetCampaignParams) (*Campaign, error) {
	d.Logger.Info("Fetching campaign", "id", params.ID)
	campaign, err := d.service.FetchCampaign(params)
	if err != nil {
		d.Logger.Error("Failed to fetch campaign", err)
	}
	return campaign, err
}

func (d *LoggingServiceDecorator) ComposeEmail(params ComposeEmailParams) string {
	d.Logger.Info("Composing email for campaign", "campaignID", params.Campaign.ID)
	email := d.service.ComposeEmail(params)
	d.Logger.Debug("Email composed", "length", len(email))
	return email
}

func (d *LoggingServiceDecorator) Error(message string, err error, params ...interface{}) {
	d.Logger.Error(message, err, params...)
}

func (d *LoggingServiceDecorator) Info(message string, params ...interface{}) {
	d.Logger.Info(message, params...)
}

func (d *LoggingServiceDecorator) Warn(message string, params ...interface{}) {
	d.Logger.Warn(message, params...)
}
