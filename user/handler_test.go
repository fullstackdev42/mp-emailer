package user

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) VerifyUser(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

// Update the RegisterUser method to match the ServiceInterface
func (m *MockService) RegisterUser(username, password, email string) error {
	args := m.Called(username, password, email)
	return args.Error(0)
}

func TestHandler_HandleLogin(t *testing.T) {
	type fields struct {
		service ServiceInterface
		logger  loggo.LoggerInterface
	}
	tests := []struct {
		name           string
		fields         fields
		setupMock      func(*MockService)
		wantErr        bool
		wantStatusCode int
		wantBody       string
	}{
		{
			name: "Invalid credentials",
			fields: fields{
				service: func() ServiceInterface {
					ms := new(MockService)
					return ms
				}(),
				logger: loggo.NewMockLogger(gomock.NewController(t)),
			},
			setupMock: func(ms *MockService) {
				ms.On("VerifyUser", "testuser", "wrongpassword").Return("", echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password"))
			},
			wantErr:        false,
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       "Invalid username or password",
		},
		// TODO: Add more test cases (e.g., successful login, server error, etc.)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(url.Values{"username": {"testuser"}, "password": {"wrongpassword"}}.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := &Handler{
				service: tt.fields.service, // Ensure this is of type ServiceInterface
				logger:  tt.fields.logger.(*loggo.Logger),
			}
			// Setup mock expectations
			if mockService, ok := h.service.(*MockService); ok && tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			// Call the method
			err := h.HandleLogin(c)

			// Assertions
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.HandleLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.wantBody)

			// Verify mock expectations
			if ms, ok := h.service.(*MockService); ok {
				ms.AssertExpectations(t)
			}
		})
	}
}
