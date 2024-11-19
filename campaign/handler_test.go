package campaign_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jonesrussell/mp-emailer/campaign"
	"github.com/jonesrussell/mp-emailer/internal/testutil"
	"github.com/jonesrussell/mp-emailer/shared"
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
	const invalidUUID = "invalid-uuid"

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
		{
			name: "invalid campaign ID",
			setupMocks: func() {
				s.Logger.EXPECT().Debug("CampaignGET: Starting")
				s.ErrorHandler.EXPECT().HandleHTTPError(
					mock.Anything,
					mock.Anything,
					"Invalid campaign ID",
					http.StatusBadRequest,
				).Return(echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID"))
			},
			campaignID:     invalidUUID,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid campaign ID",
		},
		{
			name: "campaign not found",
			setupMocks: func() {
				s.Logger.EXPECT().Debug("CampaignGET: Starting")
				s.Logger.EXPECT().Debug("CampaignGET: Parsed ID", "id", campaignID)

				s.CampaignService.EXPECT().FetchCampaign(
					campaign.GetCampaignParams{ID: campaignID},
				).Return(nil, campaign.ErrCampaignNotFound)

				httpError := echo.NewHTTPError(http.StatusNotFound, "Campaign not found")

				s.ErrorHandler.EXPECT().HandleHTTPError(
					mock.Anything,
					campaign.ErrCampaignNotFound,
					"Campaign not found",
					http.StatusNotFound,
				).Return(httpError).Once()
			},
			campaignID:     campaignID.String(),
			expectedStatus: http.StatusNotFound,
			expectedError:  "Campaign not found",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()

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
				s.Equal(tt.expectedError, httpErr.Message)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedStatus, s.Recorder.Code)
			}
		})
	}
}

func (s *HandlerTestSuite) TestGetCampaigns() {
	campaigns := []campaign.Campaign{
		{
			Name: "Campaign 1",
			BaseModel: shared.BaseModel{
				ID: uuid.New(),
			},
		},
		{
			Name: "Campaign 2",
			BaseModel: shared.BaseModel{
				ID: uuid.New(),
			},
		},
	}

	tests := []struct {
		name           string
		setupMocks     func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful campaigns fetch",
			setupMocks: func() {
				s.Logger.EXPECT().Debug("Handling GetCampaigns request")
				s.Logger.EXPECT().Debug("Rendering all campaigns", "count", len(campaigns))

				s.CampaignService.EXPECT().GetCampaigns().Return(campaigns, nil)

				s.TemplateRenderer.EXPECT().Render(
					mock.Anything,
					"campaigns",
					mock.MatchedBy(func(data shared.Data) bool {
						content, ok := data.Content.(map[string]interface{})
						return ok && data.Title == "Campaigns" &&
							data.PageName == "campaigns" &&
							len(content["Campaigns"].([]campaign.Campaign)) == len(campaigns)
					}),
					mock.Anything,
				).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "service error",
			setupMocks: func() {
				s.Logger.EXPECT().Debug("Handling GetCampaigns request")

				dbErr := errors.New("database error")
				s.CampaignService.EXPECT().GetCampaigns().Return(nil, dbErr)

				// Create the HTTP error that will be returned
				httpError := echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")

				s.ErrorHandler.EXPECT().HandleHTTPError(
					mock.Anything, // echo.Context
					dbErr,         // original error
					"Internal server error",
					http.StatusInternalServerError,
				).Return(httpError).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Internal server error",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Clear any previous mock calls
			s.SetupTest()

			tt.setupMocks()

			c := s.NewContext(http.MethodGet, "/campaigns", nil)
			err := s.handler.GetCampaigns(c)

			if tt.expectedError != "" {
				s.Error(err)
				httpErr, ok := err.(*echo.HTTPError)
				s.True(ok)
				s.Equal(tt.expectedStatus, httpErr.Code)
				s.Equal(tt.expectedError, httpErr.Message)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedStatus, s.Recorder.Code)
			}
		})
	}
}

func (s *HandlerTestSuite) TestCreateCampaignForm() {
	s.Run("successful form render", func() {
		s.Logger.EXPECT().Debug("Handling CreateCampaignForm request")

		s.TemplateRenderer.EXPECT().Render(
			mock.Anything,
			"campaign_create",
			mock.MatchedBy(func(data shared.Data) bool {
				return data.Title == "Create Campaign" &&
					data.PageName == "campaign_create" &&
					data.Content == nil
			}),
			mock.Anything,
		).Return(nil)

		c := s.NewContext(http.MethodGet, "/campaign/create", nil)
		err := s.handler.CreateCampaignForm(c)

		s.NoError(err)
		s.Equal(http.StatusOK, s.Recorder.Code)
	})
}
