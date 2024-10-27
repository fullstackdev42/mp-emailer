package campaign

import (
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/jonesrussell/loggo"
)

// CampaignLoggingServiceDecorator adds logging functionality to the Campaign ServiceInterface
type CampaignLoggingServiceDecorator struct {
	shared.LoggingServiceDecorator
	service ServiceInterface
}

// NewCampaignLoggingServiceDecorator creates a new instance of CampaignLoggingServiceDecorator
func NewCampaignLoggingServiceDecorator(service ServiceInterface, logger loggo.LoggerInterface) *CampaignLoggingServiceDecorator {
	return &CampaignLoggingServiceDecorator{
		LoggingServiceDecorator: shared.LoggingServiceDecorator{
			Logger: logger,
		},
		service: service,
	}
}

func (d *CampaignLoggingServiceDecorator) CreateCampaign(dto *CreateCampaignDTO) (*Campaign, error) {
	d.Logger.Info("Creating new campaign", "name", dto.Name)
	campaign, err := d.service.CreateCampaign(dto)
	if err != nil {
		d.Logger.Error("Failed to create campaign", err)
	}
	return campaign, err
}

func (d *CampaignLoggingServiceDecorator) UpdateCampaign(dto *UpdateCampaignDTO) error {
	d.Logger.Info("Updating campaign", "id", dto.ID)
	err := d.service.UpdateCampaign(dto)
	if err != nil {
		d.Logger.Error("Failed to update campaign", err)
	}
	return err
}

func (d *CampaignLoggingServiceDecorator) GetCampaignByID(params GetCampaignParams) (*Campaign, error) {
	d.Logger.Info("Fetching campaign", "id", params.ID)
	campaign, err := d.service.GetCampaignByID(params)
	if err != nil {
		d.Logger.Error("Failed to fetch campaign", err)
	}
	return campaign, err
}

func (d *CampaignLoggingServiceDecorator) GetAllCampaigns() ([]Campaign, error) {
	d.Logger.Info("Fetching all campaigns")
	campaigns, err := d.service.GetAllCampaigns()
	if err != nil {
		d.Logger.Error("Failed to fetch all campaigns", err)
	}
	return campaigns, err
}

func (d *CampaignLoggingServiceDecorator) DeleteCampaign(params DeleteCampaignParams) error {
	d.Logger.Info("Deleting campaign", "id", params.ID)
	err := d.service.DeleteCampaign(params)
	if err != nil {
		d.Logger.Error("Failed to delete campaign", err)
	}
	return err
}

func (d *CampaignLoggingServiceDecorator) FetchCampaign(params GetCampaignParams) (*Campaign, error) {
	d.Logger.Info("Fetching campaign", "id", params.ID)
	campaign, err := d.service.FetchCampaign(params)
	if err != nil {
		d.Logger.Error("Failed to fetch campaign", err)
	}
	return campaign, err
}

func (d *CampaignLoggingServiceDecorator) ComposeEmail(params ComposeEmailParams) string {
	d.Logger.Info("Composing email for campaign", "campaignID", params.Campaign.ID)
	email := d.service.ComposeEmail(params)
	d.Logger.Debug("Email composed", "length", len(email))
	return email
}
