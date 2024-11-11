package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fullstackdev42/mp-emailer/api"
	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/mocks"
	mocksCampaign "github.com/fullstackdev42/mp-emailer/mocks/campaign"
	mocksShared "github.com/fullstackdev42/mp-emailer/mocks/shared"
	mocksUser "github.com/fullstackdev42/mp-emailer/mocks/user"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type APITestSuite struct {
	handler          *api.Handler
	mockCampaign     *mocksCampaign.MockServiceInterface
	mockUser         *mocksUser.MockServiceInterface
	mockLogger       *mocks.MockLoggerInterface
	mockErrorHandler *mocksShared.MockErrorHandlerInterface
	echo             *echo.Echo
}

func setupAPITest(t *testing.T) *APITestSuite {
	suite := &APITestSuite{
		mockCampaign:     mocksCampaign.NewMockServiceInterface(t),
		mockUser:         mocksUser.NewMockServiceInterface(t),
		mockLogger:       mocks.NewMockLoggerInterface(t),
		mockErrorHandler: mocksShared.NewMockErrorHandlerInterface(t),
		echo:             echo.New(),
	}

	suite.handler = api.NewHandler(api.HandlerParams{
		CampaignService: suite.mockCampaign,
		UserService:     suite.mockUser,
		Logger:          suite.mockLogger,
		ErrorHandler:    suite.mockErrorHandler,
		JWTExpiry:       3600,
	})

	return suite
}

func (s *APITestSuite) tearDown() {
	s.mockCampaign = nil
	s.mockUser = nil
	s.mockLogger = nil
	s.mockErrorHandler = nil
	s.handler = nil
}

func TestGetCampaigns(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*APITestSuite)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful fetch",
			setupMocks: func(s *APITestSuite) {
				campaigns := []campaign.Campaign{{
					BaseModel: shared.BaseModel{ID: uuid.New()},
					Name:      "Test Campaign",
				}}

				s.mockCampaign.EXPECT().GetCampaigns().Return(campaigns, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "service error",
			setupMocks: func(s *APITestSuite) {
				s.mockCampaign.EXPECT().GetCampaigns().Return(nil, assert.AnError)
				s.mockErrorHandler.EXPECT().HandleHTTPError(
					mock.Anything,
					assert.AnError,
					"Error fetching campaigns",
					http.StatusInternalServerError,
				).Return(echo.NewHTTPError(http.StatusInternalServerError, "Error fetching campaigns"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setupAPITest(t)
			defer suite.tearDown()

			tt.setupMocks(suite)

			req := httptest.NewRequest(http.MethodGet, "/api/campaigns", nil)
			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)

			err := suite.handler.GetCampaigns(c)

			if tt.expectedStatus != http.StatusOK {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, he.Code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				if tt.expectedBody != "" {
					assert.JSONEq(t, tt.expectedBody, rec.Body.String())
				}
			}
		})
	}
}
