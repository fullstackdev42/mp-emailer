package campaign_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fullstackdev42/mp-emailer/config"

	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/mocks"
	mocksCampaign "github.com/fullstackdev42/mp-emailer/mocks/campaign"
	mocksEmail "github.com/fullstackdev42/mp-emailer/mocks/email"
	mocksShared "github.com/fullstackdev42/mp-emailer/mocks/shared"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type handlerTestSuite struct {
	handler              *campaign.Handler
	mockService          *mocksCampaign.MockServiceInterface
	mockLookupService    *mocksCampaign.MockRepresentativeLookupServiceInterface
	mockEmailService     *mocksEmail.MockService
	mockClient           *mocksCampaign.MockClientInterface
	mockLogger           *mocks.MockLoggerInterface
	mockErrorHandler     *mocksShared.MockErrorHandlerInterface
	mockTemplateRenderer *mocksShared.MockTemplateRendererInterface
}

func setupHandlerTest(t *testing.T) *handlerTestSuite {
	suite := &handlerTestSuite{
		mockService:          mocksCampaign.NewMockServiceInterface(t),
		mockLookupService:    mocksCampaign.NewMockRepresentativeLookupServiceInterface(t),
		mockEmailService:     mocksEmail.NewMockService(t),
		mockClient:           mocksCampaign.NewMockClientInterface(t),
		mockLogger:           mocks.NewMockLoggerInterface(t),
		mockErrorHandler:     mocksShared.NewMockErrorHandlerInterface(t),
		mockTemplateRenderer: mocksShared.NewMockTemplateRendererInterface(t),
	}

	params := campaign.HandlerParams{
		BaseHandlerParams: shared.BaseHandlerParams{
			Logger:           suite.mockLogger,
			ErrorHandler:     suite.mockErrorHandler,
			TemplateRenderer: suite.mockTemplateRenderer,
			Store:            sessions.NewCookieStore([]byte("test")),
			Config:           &config.Config{},
		},
		Service:                     suite.mockService,
		Logger:                      suite.mockLogger,
		RepresentativeLookupService: suite.mockLookupService,
		EmailService:                suite.mockEmailService,
		Client:                      suite.mockClient,
	}

	result, err := campaign.NewHandler(params)
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	suite.handler = result.Handler
	return suite
}

func TestCampaignGET(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*handlerTestSuite)
		campaignID     string
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful campaign fetch",
			setupMocks: func(s *handlerTestSuite) {
				s.mockLogger.EXPECT().Debug("CampaignGET: Starting")
				s.mockLogger.EXPECT().Debug("CampaignGET: Parsed ID", "id", 1)
				s.mockLogger.EXPECT().Debug("CampaignGET: Campaign fetched successfully", "id", 1)

				campaignTest := &campaign.Campaign{ID: 1, Name: "Test Campaign"}
				s.mockService.EXPECT().FetchCampaign(
					campaign.GetCampaignParams{ID: 1},
				).Return(campaignTest, nil)

				s.mockTemplateRenderer.EXPECT().Render(
					mock.Anything,
					"campaign",
					mock.MatchedBy(func(data map[string]interface{}) bool {
						_, hasCampaign := data["Campaign"]
						_, hasPageName := data["PageName"]
						_, hasTitle := data["Title"]
						return hasCampaign && hasPageName && hasTitle
					}),
					mock.Anything,
				).Return(nil)
			},
			campaignID:     "1",
			expectedStatus: http.StatusOK,
		},
		{
			name:       "invalid campaign ID",
			campaignID: "invalid",
			setupMocks: func(s *handlerTestSuite) {
				s.mockLogger.EXPECT().Debug("CampaignGET: Starting")

				s.mockErrorHandler.EXPECT().HandleHTTPError(
					mock.Anything,
					mock.Anything,
					"Invalid campaign ID",
					http.StatusBadRequest,
				).Return(echo.NewHTTPError(http.StatusBadRequest, "Invalid campaign ID"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid campaign ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setupHandlerTest(t)
			tt.setupMocks(suite)

			req := httptest.NewRequest(http.MethodGet, "/campaign/"+tt.campaignID, nil)
			rec := httptest.NewRecorder()
			e := echo.New()
			e.Renderer = suite.mockTemplateRenderer
			c := e.NewContext(req, rec)
			c.SetPath("/campaign/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.campaignID)

			err := suite.handler.CampaignGET(c)

			if tt.expectedError != "" {
				assert.Error(t, err)
				httpErr, ok := err.(*echo.HTTPError)
				assert.True(t, ok, "Expected HTTP error")
				assert.Equal(t, tt.expectedStatus, httpErr.Code)
				assert.Contains(t, httpErr.Message, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestGetCampaigns(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*handlerTestSuite)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful campaigns fetch",
			setupMocks: func(s *handlerTestSuite) {
				campaigns := []campaign.Campaign{{ID: 1, Name: "Test Campaign"}}

				s.mockLogger.EXPECT().Debug("GetCampaigns: Starting")
				s.mockLogger.EXPECT().Debug("GetCampaigns: Campaigns fetched", "count", 1)
				s.mockLogger.EXPECT().Debug("GetCampaigns: Template rendered successfully")

				s.mockService.EXPECT().GetCampaigns().Return(campaigns, nil)
				s.mockTemplateRenderer.EXPECT().Render(
					mock.Anything,
					"campaigns",
					mock.Anything,
					mock.Anything,
				).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setupHandlerTest(t)
			tt.setupMocks(suite)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/campaigns", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := suite.handler.GetCampaigns(c)

			if tt.expectedError != "" {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, he.Code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}
