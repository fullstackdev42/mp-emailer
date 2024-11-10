package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/mocks"
	mocksCampaign "github.com/fullstackdev42/mp-emailer/mocks/campaign"
	mocksEmail "github.com/fullstackdev42/mp-emailer/mocks/email"
	mocksMiddleware "github.com/fullstackdev42/mp-emailer/mocks/middleware"
	mocksShared "github.com/fullstackdev42/mp-emailer/mocks/shared"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	handler        HandlerInterface
	mockStore      *mocksMiddleware.MockSessionStore
	mockTemplate   *mocksShared.MockTemplateRendererInterface
	mockCampaign   *mocksCampaign.MockServiceInterface
	mockEmail      *mocksEmail.MockService
	mockLogger     *mocks.MockLoggerInterface
	mockErrHandler *mocksShared.MockErrorHandlerInterface
	echoContext    echo.Context
}

func (suite *HandlerTestSuite) SetupTest() {
	// Initialize mocks
	suite.mockStore = mocksMiddleware.NewMockSessionStore(suite.T())
	suite.mockTemplate = new(mocksShared.MockTemplateRendererInterface)
	suite.mockCampaign = new(mocksCampaign.MockServiceInterface)
	suite.mockEmail = new(mocksEmail.MockService)
	suite.mockLogger = new(mocks.MockLoggerInterface)
	suite.mockErrHandler = new(mocksShared.MockErrorHandlerInterface)

	// Create handler with mocked dependencies
	suite.handler = NewHandler(HandlerParams{
		Store:           suite.mockStore,
		TemplateManager: suite.mockTemplate,
		CampaignService: suite.mockCampaign,
		ErrorHandler:    suite.mockErrHandler,
		EmailService:    suite.mockEmail,
		Logger:          suite.mockLogger,
	})

	// Setup echo context for tests
	e := echo.New()
	e.Renderer = suite.mockTemplate
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	suite.echoContext = e.NewContext(req, rec)
}

// Add TearDownTest to clean up after each test
func (suite *HandlerTestSuite) TearDownTest() {
	// Clear all expectations from mocks
	suite.mockStore = nil
	suite.mockTemplate = nil
	suite.mockCampaign = nil
	suite.mockEmail = nil
	suite.mockLogger = nil
	suite.mockErrHandler = nil
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestHandleIndex() {
	testCases := []struct {
		name           string
		setupMocks     func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "successful_index_render",
			setupMocks: func() {
				campaigns := []campaign.Campaign{{
					BaseModel: shared.BaseModel{
						ID: uuid.MustParse("e513302d-4563-47c4-932f-d22af5c07e62"),
					},
					Name:        "Test Campaign",
					Description: "Test Description",
					Template:    "Test Template",
					OwnerID:     uuid.MustParse("f623302d-4563-47c4-932f-d22af5c07e62"),
				}}
				suite.mockCampaign.On("GetCampaigns").Return(campaigns, nil)

				suite.mockTemplate.On("Render",
					mock.AnythingOfType("*bytes.Buffer"),
					"home",
					mock.MatchedBy(func(data *shared.Data) bool {
						return data.Error == ""
					}),
					mock.AnythingOfType("*echo.context"),
				).Return(nil)

			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "campaign_service_error",
			setupMocks: func() {
				dbError := errors.New("database error")
				suite.mockCampaign.On("GetCampaigns").Return(nil, dbError)

				// Expect logger to be called
				suite.mockLogger.On("Error", "Error fetching campaigns", dbError).Return()

				// Expect error template to be rendered
				suite.mockTemplate.On("Render",
					mock.AnythingOfType("*bytes.Buffer"),
					"error",
					mock.MatchedBy(func(data *shared.Data) bool {
						return data.Error == "Error fetching campaigns" &&
							data.StatusCode == http.StatusInternalServerError
					}),
					mock.AnythingOfType("*echo.context"),
				).Return(nil)

				// Expect error handler to be called
				suite.mockErrHandler.On("HandleHTTPError",
					suite.echoContext,
					dbError,
					"Error fetching campaigns",
					http.StatusInternalServerError,
				).Return(nil)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Setup
			suite.SetupTest() // Reset mocks before each test case
			tc.setupMocks()

			// Execute
			err := suite.handler.HandleIndex(suite.echoContext)

			// Assert
			if tc.expectedError {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}

			// Verify all mocked expectations were met
			suite.mockCampaign.AssertExpectations(suite.T())
			suite.mockTemplate.AssertExpectations(suite.T())
			suite.mockLogger.AssertExpectations(suite.T())
			suite.mockErrHandler.AssertExpectations(suite.T())
		})
	}
}

func (suite *HandlerTestSuite) TestLogging() {
	testMessage := "test message"
	testError := errors.New("test error")

	suite.Run("Info logging", func() {
		suite.mockLogger.On("Info", testMessage, mock.Anything).Return()
		suite.handler.Info(testMessage)
		suite.mockLogger.AssertExpectations(suite.T())
	})

	suite.Run("Warn logging", func() {
		suite.mockLogger.On("Warn", testMessage, mock.Anything).Return()
		suite.handler.Warn(testMessage)
		suite.mockLogger.AssertExpectations(suite.T())
	})

	suite.Run("Error logging", func() {
		suite.mockLogger.On("Error", testMessage, testError, mock.Anything).Return()
		suite.handler.Error(testMessage, testError)
		suite.mockLogger.AssertExpectations(suite.T())
	})
}
