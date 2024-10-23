package campaign

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/fullstackdev42/mp-emailer/mocks"
	"github.com/fullstackdev42/mp-emailer/mocks/email"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

// MockRenderer is a mock of echo.Renderer
type MockRenderer struct{}

func (m *MockRenderer) Render(w io.Writer, name string, _ interface{}, _ echo.Context) error {
	// For simplicity, we'll just write the template name to the response writer
	_, err := w.Write([]byte(name))
	return err
}

func TestNewHandler(t *testing.T) {
	// Mock dependencies
	mockLogger := mocks.NewMockLoggerInterface(t)
	mockRepLookupService := NewMockRepresentativeLookupServiceInterface(t)
	mockService := NewMockServiceInterface(t)
	mockEmailService := email.NewMockService(t)
	mockClient := NewMockClientInterface(t)

	type args struct {
		service                     ServiceInterface
		logger                      loggo.LoggerInterface
		representativeLookupService RepresentativeLookupServiceInterface
		emailService                *email.MockService
		client                      ClientInterface
	}

	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{
			name: "Create handler with all dependencies",
			args: args{
				service:                     mockService,
				logger:                      mockLogger,
				representativeLookupService: mockRepLookupService,
				emailService:                mockEmailService,
				client:                      mockClient,
			},
			want: &Handler{
				service:                     mockService,
				logger:                      mockLogger,
				representativeLookupService: mockRepLookupService,
				emailService:                mockEmailService,
				client:                      mockClient,
			},
		},
		{
			name: "Create handler with nil logger",
			args: args{
				service:                     mockService,
				logger:                      nil,
				representativeLookupService: mockRepLookupService,
				emailService:                mockEmailService,
				client:                      mockClient,
			},
			want: &Handler{
				service:                     mockService,
				logger:                      nil,
				representativeLookupService: mockRepLookupService,
				emailService:                mockEmailService,
				client:                      mockClient,
			},
		},
		{
			name: "Create handler with nil client",
			args: args{
				service:                     mockService,
				logger:                      mockLogger,
				representativeLookupService: mockRepLookupService,
				emailService:                mockEmailService,
				client:                      nil,
			},
			want: &Handler{
				service:                     mockService,
				logger:                      mockLogger,
				representativeLookupService: mockRepLookupService,
				emailService:                mockEmailService,
				client:                      nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHandler(tt.args.service, tt.args.logger, tt.args.representativeLookupService, tt.args.emailService, tt.args.client)
			if got.service != tt.want.service {
				t.Errorf("NewHandler().service = %v, want %v", got.service, tt.want.service)
			}
			if got.logger != tt.want.logger {
				t.Errorf("NewHandler().logger = %v, want %v", got.logger, tt.want.logger)
			}
			if got.representativeLookupService != tt.want.representativeLookupService {
				t.Errorf("NewHandler().representativeLookupService = %v, want %v", got.representativeLookupService, tt.want.representativeLookupService)
			}
			if got.emailService != tt.want.emailService {
				t.Errorf("NewHandler().emailService = %v, want %v", got.emailService, tt.want.emailService)
			}
			if got.client != tt.want.client {
				t.Errorf("NewHandler().client = %v, want %v", got.client, tt.want.client)
			}
		})
	}
}

func TestHandler_CampaignGET(t *testing.T) {
	tests := []struct {
		name         string
		campaignID   string
		mockCampaign *Campaign
		mockReps     []Representative
		mockFilters  map[string]string
		mockError    error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Successful campaign retrieval",
			campaignID:   "1",
			mockCampaign: &Campaign{ID: 1, Name: "Test Campaign"},
			mockReps:     []Representative{},
			mockFilters:  map[string]string{},
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: "campaign_details.html",
		},
		{
			name:         "Campaign not found",
			campaignID:   "2",
			mockCampaign: nil,
			mockReps:     []Representative{},
			mockFilters:  map[string]string{},
			mockError:    ErrCampaignNotFound,
			expectedCode: http.StatusNotFound,
			expectedBody: "error.html",
		},
		{
			name:         "Internal server error",
			campaignID:   "3",
			mockCampaign: nil,
			mockReps:     []Representative{},
			mockFilters:  map[string]string{},
			mockError:    errors.New("internal server error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "error.html",
		},
		{
			name:         "Invalid campaign ID",
			campaignID:   "invalid",
			mockCampaign: nil,
			mockReps:     []Representative{},
			mockFilters:  map[string]string{},
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "error.html",
		},
		{
			name:         "Zero campaign ID",
			campaignID:   "0",
			mockCampaign: nil,
			mockReps:     []Representative{},
			mockFilters:  map[string]string{},
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "error.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockServiceInterface(t)
			mockLogger := mocks.NewMockLoggerInterface(t)
			mockRepLookupService := NewMockRepresentativeLookupServiceInterface(t)
			mockEmailService := email.NewMockService(t)
			mockClient := NewMockClientInterface(t)

			e := echo.New()
			e.Renderer = &MockRenderer{}

			// Set up expectations for the mock logger
			mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()

			// Set up expectations for the mock service
			if tt.campaignID != "invalid" {
				campaignID, _ := strconv.Atoi(tt.campaignID)
				if campaignID == 0 {
					mockService.EXPECT().FetchCampaign(campaignID).Return(nil, errors.New("invalid campaign ID"))
				} else {
					mockService.EXPECT().FetchCampaign(campaignID).Return(tt.mockCampaign, tt.mockError)
				}

				// If the campaign has a postal code, set up expectation for FetchRepresentatives
				if tt.mockCampaign != nil && tt.mockCampaign.PostalCode != "" {
					mockRepLookupService.EXPECT().FetchRepresentatives(tt.mockCampaign.PostalCode).Return(tt.mockReps, nil)

					// If filters are provided, set up expectation for FilterRepresentatives
					if len(tt.mockFilters) > 0 {
						mockRepLookupService.EXPECT().FilterRepresentatives(tt.mockReps, tt.mockFilters).Return(tt.mockReps)
					}
				}
			}

			h := NewHandler(mockService, mockLogger, mockRepLookupService, mockEmailService, mockClient)
			req := httptest.NewRequest(http.MethodGet, "/campaigns/"+tt.campaignID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.campaignID)

			// Call the handler
			err := h.CampaignGET(c)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.expectedBody)

			if tt.expectedCode == http.StatusOK {
				assert.Contains(t, rec.Body.String(), "campaign_details.html")
			} else {
				assert.Contains(t, rec.Body.String(), "error.html")
			}

			// Verify that all expected calls were made
			mockService.AssertExpectations(t)
			mockRepLookupService.AssertExpectations(t)
		})
	}
}
