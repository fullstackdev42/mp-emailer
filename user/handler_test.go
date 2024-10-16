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
	usermocks "github.com/fullstackdev42/mp-emailer/mocks/user"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testSessionName = "test_session"

func TestNewHandler(t *testing.T) {
	// Create mock dependencies
	mockService := usermocks.NewMockServiceInterface(t)
	mockLogger := new(loggo.MockLogger)
	mockConfig := &config.Config{}

	// Call NewHandler
	handler := NewHandler(mockService, mockLogger, mockConfig)

	// Assert that the handler is not nil
	assert.NotNil(t, handler)

	// Assert that the handler has the correct type
	assert.IsType(t, &Handler{}, handler)

	// Assert that the handler's fields are set correctly
	assert.Equal(t, mockService, handler.service)
	assert.Equal(t, mockLogger, handler.logger)
	assert.Equal(t, mockConfig, handler.config)
}

func TestHandler_RegisterPOST(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*usermocks.MockServiceInterface, *loggo.MockLogger, *mocks.MockSessionStore)
		inputBody      string
		wantStatusCode int
		wantRedirect   string
	}{
		{
			name: "Successful registration",
			setupMock: func(ms *usermocks.MockServiceInterface, ml *loggo.MockLogger, mss *mocks.MockSessionStore) {
				ms.EXPECT().RegisterUser("newuser", "password123", "newuser@example.com").Return(nil)
				ml.On("Debug", "User registered successfully", mock.Anything).Return()
				session := sessions.NewSession(mss, testSessionName)
				mss.On("Get", mock.Anything, testSessionName).Return(session, nil)
				// We don't expect Save to be called in the successful case
			},
			inputBody:      `username=newuser&password=password123&email=newuser@example.com`,
			wantStatusCode: http.StatusSeeOther,
			wantRedirect:   "/login",
		},
		{
			name: "Missing required fields",
			setupMock: func(_ *usermocks.MockServiceInterface, ml *loggo.MockLogger, mss *mocks.MockSessionStore) {
				ml.On("Warn", "Missing required fields", mock.Anything).Return()
				session := sessions.NewSession(mss, testSessionName)
				mss.On("Get", mock.Anything, testSessionName).Return(session, nil)
				mss.On("Save", mock.Anything, mock.Anything, mock.AnythingOfType("*sessions.Session")).Return(nil)
			},
			inputBody:      `username=newuser&password=`,
			wantStatusCode: http.StatusSeeOther,
			wantRedirect:   "/register",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := usermocks.NewMockServiceInterface(t)
			mockLogger := new(loggo.MockLogger)
			mockSessionStore := new(mocks.MockSessionStore)

			if tt.setupMock != nil {
				tt.setupMock(mockService, mockLogger, mockSessionStore)
			}

			handler := NewHandler(mockService, mockLogger, &config.Config{
				SessionName: testSessionName,
			})

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.Set("_session_store", mockSessionStore)

			err := handler.RegisterPOST(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatusCode, rec.Code)
			assert.Equal(t, tt.wantRedirect, rec.Header().Get("Location"))

			mockService.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockSessionStore.AssertExpectations(t)
		})
	}
}

func TestHandler_LoginPOST(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*usermocks.MockServiceInterface, *loggo.MockLogger, *mocks.MockSessionStore)
		method         string
		username       string
		password       string
		wantStatusCode int
		wantBody       string
		wantRedirect   string
	}{
		{
			name: "Successful login",
			setupMock: func(ms *usermocks.MockServiceInterface, ml *loggo.MockLogger, mss *mocks.MockSessionStore) {
				ms.EXPECT().VerifyUser("validuser", "validpass").Return("123", nil)
				ml.On("Debug", mock.Anything, mock.Anything).Return()
				session := sessions.NewSession(mss, testSessionName)
				mss.On("Get", mock.Anything, testSessionName).Return(session, nil)
				mss.On("Save", mock.Anything, mock.Anything, mock.MatchedBy(func(s *sessions.Session) bool {
					return s.Values["userID"] == "123" && s.Values["username"] == "validuser"
				})).Return(nil)
			},
			method:         http.MethodPost,
			username:       "validuser",
			password:       "validpass",
			wantStatusCode: http.StatusSeeOther,
			wantRedirect:   "/campaigns",
		},
		{
			name: "Invalid credentials",
			setupMock: func(ms *usermocks.MockServiceInterface, ml *loggo.MockLogger, mss *mocks.MockSessionStore) {
				ms.EXPECT().VerifyUser("testuser", "wrongpassword").Return("", fmt.Errorf("invalid username or password"))
				ml.On("Debug", mock.Anything, mock.Anything).Return()
				ml.On("Warn", mock.Anything, mock.Anything).Return()
				session := sessions.NewSession(mss, testSessionName)
				mss.On("Get", mock.Anything, testSessionName).Return(session, nil)
				mss.On("Save", mock.Anything, mock.Anything, mock.MatchedBy(func(s *sessions.Session) bool {
					return s.Values["error"] != nil
				})).Return(nil)
			},
			method:         http.MethodPost,
			username:       "testuser",
			password:       "wrongpassword",
			wantStatusCode: http.StatusSeeOther,
			wantRedirect:   "/login",
		},
		{
			name: "Empty username",
			setupMock: func(_ *usermocks.MockServiceInterface, ml *loggo.MockLogger, _ *mocks.MockSessionStore) {
				ml.On("Debug", mock.Anything, mock.Anything).Return()
				ml.On("Warn", mock.Anything, mock.Anything).Return()
			},
			method:         http.MethodPost,
			username:       "",
			password:       "somepassword",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "Username and password are required",
		},
		{
			name: "Empty password",
			setupMock: func(_ *usermocks.MockServiceInterface, ml *loggo.MockLogger, _ *mocks.MockSessionStore) {
				ml.On("Debug", mock.Anything, mock.Anything).Return()
				ml.On("Warn", mock.Anything, mock.Anything).Return()
				ml.On("Debug", mock.Anything, mock.Anything).Return()
			},
			method:         http.MethodPost,
			username:       "someuser",
			password:       "",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "Username and password are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := usermocks.NewMockServiceInterface(t)
			mockLogger := new(loggo.MockLogger)
			mockSessionStore := new(mocks.MockSessionStore)

			if tt.setupMock != nil {
				tt.setupMock(mockService, mockLogger, mockSessionStore)
			}

			config := &config.Config{
				SessionName: testSessionName,
			}
			handler := NewHandler(mockService, mockLogger, config)

			e := echo.New()
			req := httptest.NewRequest(tt.method, "/login", strings.NewReader(url.Values{
				"username": {tt.username},
				"password": {tt.password},
			}.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.Set("_session_store", mockSessionStore)

			err := handler.LoginPOST(c)

			t.Logf("Response status: %d", rec.Code)
			t.Logf("Response headers: %+v", rec.Header())
			t.Logf("Response body: %s", rec.Body.String())

			if tt.wantStatusCode == http.StatusSeeOther {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatusCode, rec.Code)
				assert.Equal(t, tt.wantRedirect, rec.Header().Get("Location"))
			} else {
				if assert.Error(t, err) {
					httpError, ok := err.(*echo.HTTPError)
					if assert.True(t, ok) {
						assert.Equal(t, tt.wantStatusCode, httpError.Code)
						assert.Equal(t, tt.wantBody, httpError.Message)
					}
				}
			}

			mockService.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockSessionStore.AssertExpectations(t)
		})
	}
}

func TestHandler_LogoutGET(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*usermocks.MockServiceInterface, *loggo.MockLogger, *mocks.MockSessionStore)
		wantStatusCode int
		wantRedirect   string
		wantErr        bool
	}{
		{
			name: "Successful logout",
			setupMock: func(_ *usermocks.MockServiceInterface, ml *loggo.MockLogger, mss *mocks.MockSessionStore) {
				ml.On("Debug", "Handling logout request", map[string]interface{}{"path": "/logout"}).Return()
				session := &sessions.Session{Values: make(map[interface{}]interface{}), Options: &sessions.Options{}}
				mss.On("Get", mock.Anything, testSessionName).Return(session, nil)
				mss.On("Save", mock.Anything, mock.Anything, session).Return(nil)
			},
			wantStatusCode: http.StatusSeeOther,
			wantRedirect:   "/",
			wantErr:        false,
		},
		{
			name: "Error getting session",
			setupMock: func(_ *usermocks.MockServiceInterface, ml *loggo.MockLogger, mss *mocks.MockSessionStore) {
				ml.On("Debug", "Handling logout request", map[string]interface{}{"path": "/logout"}).Return()
				mss.On("Get", mock.Anything, testSessionName).Return(nil, fmt.Errorf("failed to get session"))
				ml.On("Error", "Failed to get session", mock.AnythingOfType("error")).Return()
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name: "Error saving session",
			setupMock: func(_ *usermocks.MockServiceInterface, ml *loggo.MockLogger, mss *mocks.MockSessionStore) {
				ml.On("Debug", "Handling logout request", map[string]interface{}{"path": "/logout"}).Return()
				session := &sessions.Session{Values: make(map[interface{}]interface{}), Options: &sessions.Options{}}
				mss.On("Get", mock.Anything, testSessionName).Return(session, nil)
				mss.On("Save", mock.Anything, mock.Anything, session).Return(fmt.Errorf("failed to save session"))
				ml.On("Error", "Failed to save session", mock.AnythingOfType("error")).Return()
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := usermocks.NewMockServiceInterface(t)
			mockLogger := new(loggo.MockLogger)
			mockSessionStore := new(mocks.MockSessionStore)

			tt.setupMock(mockService, mockLogger, mockSessionStore)

			config := &config.Config{SessionName: testSessionName}
			handler := NewHandler(mockService, mockLogger, config)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/logout", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("_session_store", mockSessionStore)

			err := handler.LogoutGET(c)

			if tt.wantErr {
				assert.Error(t, err)
				httpErr, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.wantStatusCode, httpErr.Code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatusCode, rec.Code)
				assert.Equal(t, tt.wantRedirect, rec.Header().Get("Location"))
			}

			mockService.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockSessionStore.AssertExpectations(t)
		})
	}
}
