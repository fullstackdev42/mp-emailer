package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jonesrussell/mp-emailer/api"
	"github.com/jonesrussell/mp-emailer/campaign"
	"github.com/jonesrussell/mp-emailer/logger"
	mocksCampaign "github.com/jonesrussell/mp-emailer/mocks/campaign"
	mocksLogger "github.com/jonesrussell/mp-emailer/mocks/logger"
	mocksShared "github.com/jonesrussell/mp-emailer/mocks/shared"
	mocksUser "github.com/jonesrussell/mp-emailer/mocks/user"
	"github.com/jonesrussell/mp-emailer/shared"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type APITestSuite struct {
	handler          *api.Handler
	mockCampaign     *mocksCampaign.MockServiceInterface
	mockUser         *mocksUser.MockServiceInterface
	mockLogger       logger.Interface
	mockErrorHandler *mocksShared.MockErrorHandlerInterface
	echo             *echo.Echo
}

func setupAPITest(t *testing.T) *APITestSuite {
	suite := &APITestSuite{
		mockCampaign:     mocksCampaign.NewMockServiceInterface(t),
		mockUser:         mocksUser.NewMockServiceInterface(t),
		mockLogger:       mocksLogger.NewMockInterface(t),
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
				s.mockCampaign.EXPECT().
					GetCampaigns(mock.Anything).
					Return([]campaign.Campaign{{
						BaseModel: shared.BaseModel{ID: uuid.New()},
						Name:      "Test Campaign",
					}}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "service error",
			setupMocks: func(s *APITestSuite) {
				s.mockCampaign.EXPECT().
					GetCampaigns(mock.Anything).
					Return(nil, assert.AnError)
				s.mockErrorHandler.EXPECT().
					HandleHTTPError(
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

			ctx := context.Background()
			req = req.WithContext(ctx)
			c.SetRequest(req)

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
