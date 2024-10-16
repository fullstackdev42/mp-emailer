package user

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/mocks"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testSessionName = "test_session"

func TestNewHandler(t *testing.T) {
	mockService := NewMockServiceInterface(t)
	mockConfig := &config.Config{}

	handler := NewHandler(mockService, nil, mockConfig)

	assert.NotNil(t, handler)
	assert.IsType(t, &Handler{}, handler)
	assert.Equal(t, mockService, handler.service)
	assert.Equal(t, mockConfig, handler.config)
}

func TestHandler_RegisterPOST(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockServiceInterface, mocks.MockLoggerInterface)
		inputBody      string
		wantStatusCode int
		wantRedirect   string
	}{
		{
			name: "Successful registration",
			setupMock: func(ms *MockServiceInterface, ml mocks.MockLoggerInterface) {
				ms.EXPECT().RegisterUser("newuser", "password123", "newuser@example.com").Return(nil)
				ml.EXPECT().Info("User registered successfully").Return()
			},
			inputBody:      `username=newuser&password=password123&email=newuser@example.com`,
			wantStatusCode: http.StatusSeeOther,
			wantRedirect:   "/login",
		},
		{
			name: "Missing required fields",
			setupMock: func(_ *MockServiceInterface, ml mocks.MockLoggerInterface) {
				ml.EXPECT().Warn("Missing required fields in registration form").Return()
			},
			inputBody:      `username=newuser&password=`,
			wantStatusCode: http.StatusSeeOther,
			wantRedirect:   "/register",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockServiceInterface(t)
			mockLogger := mocks.NewMockLoggerInterface(t)

			tt.setupMock(mockService, mockLogger)

			handler := NewHandler(mockService, mockLogger, &config.Config{SessionName: testSessionName})

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.RegisterPOST(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatusCode, rec.Code)
			assert.Equal(t, tt.wantRedirect, rec.Header().Get("Location"))

			mockService.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

func TestHandler_LoginPOST(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockServiceInterface, *MockSessionStore)
		username       string
		password       string
		wantStatusCode int
		wantBody       string
		wantRedirect   string
	}{
		{
			name: "Successful login",
			setupMock: func(ms *MockServiceInterface, mss *MockSessionStore) {
				ms.EXPECT().VerifyUser("validuser", "validpass").Return("123", nil)
				session := sessions.NewSession(mss, testSessionName)
				mss.On("Get", mock.Anything, testSessionName).Return(session, nil)
				mss.On("Save", mock.Anything, mock.Anything, mock.MatchedBy(func(s *sessions.Session) bool {
					return s.Values["userID"] == "123" && s.Values["username"] == "validuser"
				})).Return(nil)
			},
			username:       "validuser",
			password:       "validpass",
			wantStatusCode: http.StatusSeeOther,
			wantRedirect:   "/campaigns",
		},
		{
			name: "Invalid credentials",
			setupMock: func(ms *MockServiceInterface, mss *MockSessionStore) {
				ms.EXPECT().VerifyUser("testuser", "wrongpassword").Return("", fmt.Errorf("invalid username or password"))
				session := sessions.NewSession(mss, testSessionName)
				mss.On("Get", mock.Anything, testSessionName).Return(session, nil)
				mss.On("Save", mock.Anything, mock.Anything, mock.MatchedBy(func(s *sessions.Session) bool {
					return s.Values["error"] != nil
				})).Return(nil)
			},
			username:       "testuser",
			password:       "wrongpassword",
			wantStatusCode: http.StatusSeeOther,
			wantRedirect:   "/login",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockServiceInterface(t)
			mockSessionStore := new(MockSessionStore)
			tt.setupMock(mockService, mockSessionStore)

			handler := NewHandler(mockService, nil, &config.Config{SessionName: testSessionName})

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(url.Values{
				"username": {tt.username},
				"password": {tt.password},
			}.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("_session_store", mockSessionStore)

			err := handler.LoginPOST(c)

			if tt.wantStatusCode == http.StatusSeeOther {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatusCode, rec.Code)
				assert.Equal(t, tt.wantRedirect, rec.Header().Get("Location"))
			} else {
				assert.Error(t, err)
				httpError, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.wantStatusCode, httpError.Code)
				assert.Equal(t, tt.wantBody, httpError.Message)
			}

			mockService.AssertExpectations(t)
			mockSessionStore.AssertExpectations(t)
		})
	}
}

func TestHandler_RegisterGET(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockServiceInterface, mocks.MockLoggerInterface)
		wantStatusCode int
		wantErr        bool
	}{
		{
			name: "Successful GET request",
			setupMock: func(_ *MockServiceInterface, ml mocks.MockLoggerInterface) {
				ml.EXPECT().Debug("RegisterGET: Starting").Return()
			},
			wantStatusCode: http.StatusOK,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockServiceInterface(t)
			mockLogger := mocks.MockLoggerInterface(t)

			tt.setupMock(mockService, mockLogger)

			handler := NewHandler(mockService, mockLogger, &config.Config{SessionName: testSessionName})

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/register", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.RegisterGET(c)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)

			mockService.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}
