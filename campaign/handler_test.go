package campaign

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/mocks"
	mocksEmail "github.com/fullstackdev42/mp-emailer/mocks/email"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	mockService := new(MockServiceInterface)
	mockLogger := mocks.NewMockLoggerInterface(t)
	mockRepLookupService := new(MockRepresentativeLookupServiceInterface)
	mockEmailService := mocksEmail.NewMockService(t)
	mockClient := new(MockClientInterface)

	type args struct {
		service                     ServiceInterface
		logger                      loggo.LoggerInterface
		representativeLookupService RepresentativeLookupServiceInterface
		emailService                *mocksEmail.MockService
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
		mockReturn   *Campaign
		mockError    error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Successful campaign retrieval",
			campaignID:   "1",
			mockReturn:   &Campaign{ID: 1, Name: "Test Campaign"},
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: "campaign_details.html",
		},
		{
			name:         "Campaign not found",
			campaignID:   "2",
			mockReturn:   nil,
			mockError:    ErrCampaignNotFound,
			expectedCode: http.StatusNotFound,
			expectedBody: "error.html",
		},
		{
			name:         "Internal server error",
			campaignID:   "3",
			mockReturn:   nil,
			mockError:    errors.New("internal server error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "error.html",
		},
		{
			name:         "Invalid campaign ID",
			campaignID:   "invalid",
			mockReturn:   nil,
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "error.html",
		},
		{
			name:         "Zero campaign ID",
			campaignID:   "0",
			mockReturn:   nil,
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "error.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockServiceInterface(t)
			mockLogger := mocks.NewMockLoggerInterface(t)
			mockRepLookupService := mocks.NewMockRepresentativeLookupServiceInterface(t)
			mockMailgunClient := mocksEmail.NewMockMailgunClient(t)
			mockClient := mocks.NewMockClientInterface(t)

			// Create a new email.Service with the mock MailgunClient
			emailService := &email.Service{MailgunClient: mockMailgunClient}
			e := echo.New()
			e.Renderer = &MockRenderer{}

			if tt.mockReturn != nil || tt.mockError != nil {
				mockService.EXPECT().FetchCampaign(mock.AnythingOfType("int")).Return(tt.mockReturn, tt.mockError)
			}

			h := NewHandler(mockService, mockLogger, mockRepLookupService, *emailService, mockClient)
			req := httptest.NewRequest(http.MethodGet, "/campaigns/"+tt.campaignID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.campaignID)

			// Example: Setting expectations on MockMailgunClient
			mockMailgunClient.EXPECT().NewMessage(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mailgun.Message{})
			mockMailgunClient.EXPECT().Send(mock.Anything, mock.Anything).Return("message-id", "response", nil)

			err := h.CampaignGET(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.Equal(t, tt.expectedBody, rec.Body.String())

			if tt.expectedCode >= 400 {
				assert.Contains(t, rec.Body.String(), "error")
			}
		})
	}
}
