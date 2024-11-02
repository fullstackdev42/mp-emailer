package campaign

import (
	"github.com/jonesrussell/loggo"
)

type LoggingDecorator struct {
	service ServiceInterface
	logger  loggo.LoggerInterface
}

func NewLoggingDecorator(service ServiceInterface, logger loggo.LoggerInterface) ServiceInterface {
	return &LoggingDecorator{
		service: service,
		logger:  logger,
	}
}

// Implement LoggableService methods
func (d *LoggingDecorator) Info(message string, params ...interface{}) {
	d.logger.Info(message, params...)
}

func (d *LoggingDecorator) Warn(message string, params ...interface{}) {
	d.logger.Warn(message, params...)
}

func (d *LoggingDecorator) Error(message string, err error, params ...interface{}) {
	d.logger.Error(message, err, params...)
}

// Implement ServiceInterface methods with logging
func (d *LoggingDecorator) CreateCampaign(dto *CreateCampaignDTO) (*Campaign, error) {
	d.logger.Info("Creating campaign", "dto", dto)
	campaign, err := d.service.CreateCampaign(dto)
	if err != nil {
		d.logger.Error("Failed to create campaign", err, "dto", dto)
	}
	return campaign, err
}

func (d *LoggingDecorator) UpdateCampaign(dto *UpdateCampaignDTO) error {
	d.logger.Info("Updating campaign", "dto", dto)
	err := d.service.UpdateCampaign(dto)
	if err != nil {
		d.logger.Error("Failed to update campaign", err, "dto", dto)
	}
	return err
}

func (d *LoggingDecorator) GetCampaignByID(params GetCampaignParams) (*Campaign, error) {
	d.logger.Info("Getting campaign by ID", "params", params)
	campaign, err := d.service.GetCampaignByID(params)
	if err != nil {
		d.logger.Error("Failed to get campaign", err, "params", params)
	}
	return campaign, err
}

func (d *LoggingDecorator) GetCampaigns() ([]Campaign, error) {
	d.logger.Info("Fetching all campaigns")
	campaigns, err := d.service.GetCampaigns()
	if err != nil {
		d.logger.Error("Failed to fetch campaigns", err)
	}
	return campaigns, err
}

func (d *LoggingDecorator) DeleteCampaign(params DeleteCampaignParams) error {
	d.logger.Info("Deleting campaign", "params", params)
	err := d.service.DeleteCampaign(params)
	if err != nil {
		d.logger.Error("Failed to delete campaign", err, "params", params)
	}
	return err
}

func (d *LoggingDecorator) FetchCampaign(params GetCampaignParams) (*Campaign, error) {
	d.logger.Info("Fetching campaign", "params", params)
	campaign, err := d.service.FetchCampaign(params)
	if err != nil {
		d.logger.Error("Failed to fetch campaign", err, "params", params)
	}
	return campaign, err
}

func (d *LoggingDecorator) ComposeEmail(params ComposeEmailParams) string {
	d.logger.Info("Composing email", "params", params)
	return d.service.ComposeEmail(params)
}
