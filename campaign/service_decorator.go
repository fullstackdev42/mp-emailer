package campaign

import (
	"context"

	"github.com/jonesrussell/mp-emailer/logger"
)

type LoggingDecorator struct {
	service ServiceInterface
	Logger  logger.Interface
}

func NewLoggingDecorator(service ServiceInterface, log logger.Interface) ServiceInterface {
	return &LoggingDecorator{
		service: service,
		Logger:  log,
	}
}

// Info logs an info message with the given parameters
func (d *LoggingDecorator) ComposeEmail(ctx context.Context, params ComposeEmailParams) (string, error) {
	d.Logger.Info("Composing email", "params", params)
	email, err := d.service.ComposeEmail(ctx, params)
	if err != nil {
		d.Logger.Error("Failed to compose email", err, "params", params)
	}
	return email, err
}

// Info logs an info message with the given parameters
func (d *LoggingDecorator) Info(message string, params ...interface{}) {
	d.Logger.Info(message, params...)
}

// Warn logs a warning message with the given parameters
func (d *LoggingDecorator) Warn(message string, params ...interface{}) {
	d.Logger.Warn(message, params...)
}

// Error logs an error message with the given parameters
func (d *LoggingDecorator) Error(message string, err error, params ...interface{}) {
	d.Logger.Error(message, err, params...)
}

// CreateCampaign creates a new campaign
func (d *LoggingDecorator) CreateCampaign(ctx context.Context, dto *CreateCampaignDTO) (*Campaign, error) {
	d.Logger.Info("Creating campaign", "dto", dto)
	campaign, err := d.service.CreateCampaign(ctx, dto)
	if err != nil {
		d.Logger.Error("Failed to create campaign", err, "dto", dto)
	}
	return campaign, err
}

// UpdateCampaign updates an existing campaign
func (d *LoggingDecorator) UpdateCampaign(ctx context.Context, dto *UpdateCampaignDTO) error {
	d.Logger.Info("Updating campaign", "dto", dto)
	err := d.service.UpdateCampaign(ctx, dto)
	if err != nil {
		d.Logger.Error("Failed to update campaign", err, "dto", dto)
	}
	return err
}

// GetCampaignByID gets a campaign by its ID
func (d *LoggingDecorator) GetCampaignByID(ctx context.Context, params GetCampaignParams) (*Campaign, error) {
	d.Logger.Info("Getting campaign by ID", "params", params)
	campaign, err := d.service.GetCampaignByID(ctx, params)
	if err != nil {
		d.Logger.Error("Failed to get campaign", err, "params", params)
	}
	return campaign, err
}

// GetCampaigns gets all campaigns
func (d *LoggingDecorator) GetCampaigns(ctx context.Context) ([]Campaign, error) {
	d.Logger.Info("Fetching all campaigns")
	campaigns, err := d.service.GetCampaigns(ctx)
	if err != nil {
		d.Logger.Error("Failed to fetch campaigns", err)
	}
	return campaigns, err
}

// DeleteCampaign deletes a campaign
func (d *LoggingDecorator) DeleteCampaign(ctx context.Context, params DeleteCampaignDTO) error {
	d.Logger.Info("Deleting campaign", "params", params)
	err := d.service.DeleteCampaign(ctx, params)
	if err != nil {
		d.Logger.Error("Failed to delete campaign", err, "params", params)
	}
	return err
}

// FetchCampaign fetches a campaign
func (d *LoggingDecorator) FetchCampaign(ctx context.Context, params GetCampaignParams) (*Campaign, error) {
	d.Logger.Info("Fetching campaign", "params", params)
	campaign, err := d.service.FetchCampaign(ctx, params)
	if err != nil {
		d.Logger.Error("Failed to fetch campaign", err, "params", params)
	}
	return campaign, err
}
