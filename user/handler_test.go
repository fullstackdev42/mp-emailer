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

func (m *MockService) RegisterUser(username, password, email string) error {
	args := m.Called(username, password, email)
	return args.Error(0)
}

func TestHandler_HandleLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := loggo.NewMockLogger(ctrl)
	mockService := new(MockService)

	h := &Handler{
		service: mockService,
		logger:  mockLogger,
	}

	tests := []struct {
		name           string
		setupMock      func()
		wantStatusCode int
		wantBody       string
	}{
		{
			name: "Invalid credentials",
			setupMock: func() {
				mockService.On("VerifyUser", "testuser", "wrongpassword").Return("", echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password"))
				mockLogger.EXPECT().Debug("HandleLogin called with method: POST").Times(1)
				mockLogger.EXPECT().Debug("Login attempt for username: testuser").Times(1)
				mockLogger.EXPECT().Warn("Login failed for user: testuser").Times(1)
			},
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       "Invalid username or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(url.Values{"username": {"testuser"}, "password": {"wrongpassword"}}.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.HandleLogin(c)

			if err != nil {
				httpError, ok := err.(*echo.HTTPError)
				if ok {
					assert.Equal(t, tt.wantStatusCode, httpError.Code)
					assert.Contains(t, httpError.Message, tt.wantBody)
				} else {
					t.Fatalf("expected HTTPError, got %v", err)
				}
			} else {
				assert.Equal(t, tt.wantStatusCode, rec.Code)
				assert.Contains(t, rec.Body.String(), tt.wantBody)
			}

			mockService.AssertExpectations(t)
		})
	}
}
