package campaign_test

import (
	"net/http"
	"testing"

	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/internal/testutil"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	testutil.BaseTestSuite
	handler *campaign.Handler
}

func (s *HandlerTestSuite) SetupTest() {
	s.BaseTestSuite.SetupTest()

	// Register renderer with Echo
	s.Echo.Renderer = s.TemplateRenderer

	params := campaign.HandlerParams{
		BaseHandlerParams: shared.BaseHandlerParams{
			Logger:           s.Logger,
			ErrorHandler:     s.ErrorHandler,
			TemplateRenderer: s.TemplateRenderer,
			Store:            s.Store,
			Config:           s.Config,
		},
		Service:                     s.CampaignService,
		RepresentativeLookupService: s.RepresentativeLookupService,
		EmailService:                s.EmailService,
		Client:                      s.CampaignClient,
	}

	result, err := campaign.NewHandler(params)
	s.NoError(err)
	s.handler = result.Handler
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) TestCampaignGET() {
	campaignID := uuid.New()

	tests := []struct {
		name           string
		setupMocks     func()
		campaignID     string
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful campaign fetch",
			setupMocks: func() {
				s.Logger.EXPECT().Debug("CampaignGET: Starting")
				s.Logger.EXPECT().Debug("CampaignGET: Parsed ID", "id", campaignID)
				s.Logger.EXPECT().Debug("CampaignGET: Campaign fetched successfully", "id", campaignID)

				campaignTest := &campaign.Campaign{
					Name: "Test Campaign",
					BaseModel: shared.BaseModel{
						ID: campaignID,
					},
				}

				s.CampaignService.EXPECT().FetchCampaign(
					campaign.GetCampaignParams{ID: campaignID},
				).Return(campaignTest, nil)

				s.TemplateRenderer.EXPECT().Render(
					mock.Anything,
					"campaign",
					mock.MatchedBy(func(data map[string]interface{}) bool {
						return data["Campaign"] != nil &&
							data["PageName"] == "campaign" &&
							data["Title"] == "Campaign Details"
					}),
					mock.Anything,
				).Return(nil)
			},
			campaignID:     campaignID.String(),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMocks()

			c := s.NewContext(http.MethodGet, "/campaign/"+tt.campaignID, nil)
			c.SetPath("/campaign/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.campaignID)

			err := s.handler.CampaignGET(c)

			if tt.expectedError != "" {
				s.Error(err)
				httpErr, ok := err.(*echo.HTTPError)
				s.True(ok)
				s.Equal(tt.expectedStatus, httpErr.Code)
				s.Contains(httpErr.Message, tt.expectedError)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedStatus, s.Recorder.Code)
			}
		})
	}
}
