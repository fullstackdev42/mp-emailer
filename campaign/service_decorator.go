package campaign

import (
	"context"

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

// Info logs an info message with the given parameters
func (d *LoggingDecorator) ComposeEmail(ctx context.Context, params ComposeEmailParams) (string, error) {
	d.logger.Info("Composing email", "params", params)
	email, err := d.service.ComposeEmail(ctx, params)
	if err != nil {
		d.logger.Error("Failed to compose email", err, "params", params)
	}
	return email, err
}

// Info logs an info message with the given parameters
func (d *LoggingDecorator) Info(message string, params ...interface{}) {
	d.logger.Info(message, params...)
}

// Warn logs a warning message with the given parameters
func (d *LoggingDecorator) Warn(message string, params ...interface{}) {
	d.logger.Warn(message, params...)
}

// Error logs an error message with the given parameters
func (d *LoggingDecorator) Error(message string, err error, params ...interface{}) {
	d.logger.Error(message, err, params...)
}

// CreateCampaign creates a new campaign
func (d *LoggingDecorator) CreateCampaign(ctx context.Context, dto *CreateCampaignDTO) (*Campaign, error) {
	d.logger.Info("Creating campaign", "dto", dto)
	campaign, err := d.service.CreateCampaign(ctx, dto)
	if err != nil {
		d.logger.Error("Failed to create campaign", err, "dto", dto)
	}
	return campaign, err
}

// UpdateCampaign updates an existing campaign
func (d *LoggingDecorator) UpdateCampaign(ctx context.Context, dto *UpdateCampaignDTO) error {
	d.logger.Info("Updating campaign", "dto", dto)
	err := d.service.UpdateCampaign(ctx, dto)
	if err != nil {
		d.logger.Error("Failed to update campaign", err, "dto", dto)
	}
	return err
}

// GetCampaignByID gets a campaign by its ID
func (d *LoggingDecorator) GetCampaignByID(ctx context.Context, params GetCampaignParams) (*Campaign, error) {
	d.logger.Info("Getting campaign by ID", "params", params)
	campaign, err := d.service.GetCampaignByID(ctx, params)
	if err != nil {
		d.logger.Error("Failed to get campaign", err, "params", params)
	}
	return campaign, err
}

// GetCampaigns gets all campaigns
func (d *LoggingDecorator) GetCampaigns(ctx context.Context) ([]Campaign, error) {
	d.logger.Info("Fetching all campaigns")
	campaigns, err := d.service.GetCampaigns(ctx)
	if err != nil {
		d.logger.Error("Failed to fetch campaigns", err)
	}
	return campaigns, err
}

// DeleteCampaign deletes a campaign
func (d *LoggingDecorator) DeleteCampaign(ctx context.Context, params DeleteCampaignDTO) error {
	d.logger.Info("Deleting campaign", "params", params)
	err := d.service.DeleteCampaign(ctx, params)
	if err != nil {
		d.logger.Error("Failed to delete campaign", err, "params", params)
	}
	return err
}

// FetchCampaign fetches a campaign
func (d *LoggingDecorator) FetchCampaign(ctx context.Context, params GetCampaignParams) (*Campaign, error) {
	d.logger.Info("Fetching campaign", "params", params)
	campaign, err := d.service.FetchCampaign(ctx, params)
	if err != nil {
		d.logger.Error("Failed to fetch campaign", err, "params", params)
	}
	return campaign, err
}
