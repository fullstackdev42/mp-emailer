package campaign

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/mocks"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock of ServiceInterface
type MockService struct {
	mock.Mock
}

func (m *MockService) GetAllCampaigns() ([]Campaign, error) {
	args := m.Called()
	return args.Get(0).([]Campaign), args.Error(1)
}

func (m *MockService) GetCampaignByID(id int) (*Campaign, error) {
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

// Add the missing UpdateCampaign method
func (m *MockService) UpdateCampaign(campaign *Campaign) error {
	args := m.Called(campaign)
	return args.Error(0)
}

func (m *MockService) FetchCampaign(id int) (*Campaign, error) {
	args := m.Called(id)
	return args.Get(0).(*Campaign), args.Error(1)
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
	mockLogger := &mocks.MockLoggerInterface{}
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
		campaignID   int
		mockReturn   *Campaign
		mockError    error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Successful campaign retrieval",
			campaignID:   1,
			mockReturn:   &Campaign{ID: 1, Name: "Test Campaign"},
			expectedCode: http.StatusOK,
			expectedBody: `{"id":1,"name":"Test Campaign","template":"","owner_id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","tokens":null}`,
		},
		{
			name:         "Campaign not found",
			campaignID:   2,
			mockReturn:   nil,
			mockError:    echo.ErrNotFound,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"message":"campaign not found"}`,
		},
		{
			name:         "Internal server error",
			campaignID:   3,
			mockReturn:   nil,
			mockError:    errors.New("internal server error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockServiceInterface(t)
			mockLogger := mocks.NewMockLoggerInterface(t)
			e := echo.New()
			e.Renderer = &MockRenderer{}

			mockService.EXPECT().FetchCampaign(tt.campaignID).Return(tt.mockReturn, tt.mockError)

			h := NewHandler(mockService, mockLogger, nil, nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/campaigns/"+strconv.Itoa(tt.campaignID), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(strconv.Itoa(tt.campaignID))

			err := h.GetCampaign(c)

			if tt.mockError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.mockError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, rec.Code)
				assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			}
		})
	}
}
