package campaign

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock of ServiceInterface
type MockService struct {
	mock.Mock
}

// Update the GetAllCampaigns method to match the interface
func (m *MockService) GetAllCampaigns() ([]*Campaign, error) {
	args := m.Called()
	return args.Get(0).([]*Campaign), args.Error(1)
}

func (m *MockService) GetCampaignByID(id string) (*Campaign, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Campaign), args.Error(1)
}

// Update the ComposeEmail method to match the interface
func (m *MockService) ComposeEmail(representative Representative, campaign *Campaign, data map[string]string) string {
	args := m.Called(representative, campaign, data)
	return args.String(0)
}

// Add the missing CreateCampaign method
func (m *MockService) CreateCampaign(campaign *Campaign) error {
	args := m.Called(campaign)
	return args.Error(0)
}

func (m *MockService) DeleteCampaign(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// Add the missing ExtractAndValidatePostalCode method
func (m *MockService) ExtractAndValidatePostalCode(c echo.Context) (string, error) {
	args := m.Called(c)
	return args.String(0), args.Error(1)
}

// Add the missing UpdateCampaign method
func (m *MockService) UpdateCampaign(campaign *Campaign) error {
	args := m.Called(campaign)
	return args.Error(0)
}

// MockRepresentativeLookupService is a mock of RepresentativeLookupService
type MockRepresentativeLookupService struct {
	mock.Mock
}

// Add the missing FetchRepresentatives method
func (m *MockRepresentativeLookupService) FetchRepresentatives(postalCode string) ([]Representative, error) {
	args := m.Called(postalCode)
	return args.Get(0).([]Representative), args.Error(1)
}

func (m *MockRepresentativeLookupService) FilterRepresentatives(representatives []Representative, filters map[string]string) []Representative {
	args := m.Called(representatives, filters)
	return args.Get(0).([]Representative)
}

// MockEmailService implements email.Service for testing
type MockEmailService struct {
	mock.Mock
}

// SendEmail implements the SendEmail method of the email.Service interface
func (m *MockEmailService) SendEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

// MockClient implements ClientInterface for testing
type MockClient struct {
	mock.Mock
}

// FetchRepresentatives implements the FetchRepresentatives method of the ClientInterface
func (m *MockClient) FetchRepresentatives(postalCode string) ([]Representative, error) {
	args := m.Called(postalCode)
	return args.Get(0).([]Representative), args.Error(1)
}

// MockRenderer is a mock of echo.Renderer
type MockRenderer struct{}

func (m *MockRenderer) Render(w io.Writer, name string, _ interface{}, _ echo.Context) error {
	if name == "error.html" {
		_, err := w.Write([]byte("Error page"))
		if err != nil {
			return err
		}
	}
	return nil
}

func TestNewHandler(t *testing.T) {
	// Mock dependencies
	mockService := new(MockService)
	mockLogger := &loggo.MockLogger{}
	mockRepLookupService := new(MockRepresentativeLookupService)
	mockEmailService := &MockEmailService{}
	mockClient := &MockClient{}

	type args struct {
		service                     ServiceInterface
		logger                      loggo.LoggerInterface
		representativeLookupService RepresentativeLookupServiceInterface
		emailService                email.Service
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

			// Compare individual fields
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

func TestHandler_GetCampaign(t *testing.T) {
	tests := []struct {
		name         string
		campaignID   string
		mockReturn   interface{}
		mockError    error
		expectedCode int
	}{
		{
			name:         "Successful campaign retrieval",
			campaignID:   "1",
			mockReturn:   &Campaign{ID: 1, Name: "Test Campaign"},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Campaign not found",
			campaignID:   "999",
			mockError:    echo.ErrNotFound, // Ensure this error is returned for not found
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Internal server error",
			campaignID:   "1",
			mockError:    errors.New("database error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			mockLogger := new(loggo.MockLogger)
			e := echo.New()
			e.Renderer = &MockRenderer{}

			// Set up logger expectations
			mockLogger.On("Info", mock.Anything, mock.Anything).Return()
			mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

			if tt.mockError != nil {
				mockService.On("GetCampaignByID", tt.campaignID).Return(nil, tt.mockError)
			} else {
				mockService.On("GetCampaignByID", tt.campaignID).Return(tt.mockReturn, nil)
			}

			h := &Handler{
				service:                     mockService,
				logger:                      mockLogger,
				representativeLookupService: new(MockRepresentativeLookupService),
				emailService:                &MockEmailService{},
				client:                      &MockClient{},
			}

			req := httptest.NewRequest(http.MethodGet, "/campaigns/"+tt.campaignID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.campaignID)

			err := h.GetCampaign(c)

			if err != nil && err.Error() != tt.mockError.Error() {
				t.Errorf("Handler.GetCampaign() unexpected error = %v, want %v", err, tt.mockError)
			}

			if rec.Code != tt.expectedCode {
				t.Errorf("Handler.GetCampaign() status code = %v, want %v", rec.Code, tt.expectedCode)
			}
		})
	}
}
