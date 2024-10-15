package user

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/gorilla/sessions"
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

type mockSessionStore struct {
	sessions  map[string]*sessions.Session
	saveError error
	getError  error
}

func newMockSessionStore() *mockSessionStore {
	return &mockSessionStore{
		sessions: make(map[string]*sessions.Session),
	}
}

func (m *mockSessionStore) Get(_ *http.Request, name string) (*sessions.Session, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	session, ok := m.sessions[name]
	if !ok {
		session = sessions.NewSession(m, name)
		m.sessions[name] = session
	}
	return session, nil
}

func (m *mockSessionStore) New(_ *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(m, name)
	m.sessions[name] = session
	return session, nil
}

func (m *mockSessionStore) Save(_ *http.Request, _ http.ResponseWriter, s *sessions.Session) error {
	if m.saveError != nil {
		return m.saveError
	}
	m.sessions[s.Name()] = s
	return nil
}

const testSessionName = "test_session"

func TestNewHandler(t *testing.T) {
	// Create mock dependencies
	mockService := new(MockService)
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

func TestHandler_HandleRegister(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockService, *loggo.MockLogger)
		inputBody      string
		wantStatusCode int
		wantBody       string
	}{
		{
			name: "Successful registration",
			setupMock: func(ms *MockService, ml *loggo.MockLogger) {
				ms.On("RegisterUser", "newuser", "password123", "newuser@example.com").Return(nil)
				ml.On("Debug", "User registered successfully", mock.AnythingOfType("[]interface {}")).Return()
			},
			inputBody:      `{"username":"newuser","password":"password123","email":"newuser@example.com"}`,
			wantStatusCode: http.StatusCreated,
			wantBody:       `{"message":"User registered successfully"}`,
		},
		{
			name: "Invalid JSON",
			setupMock: func(_ *MockService, ml *loggo.MockLogger) {
				ml.On("Warn", "Invalid request body", mock.Anything).Return()
			},
			inputBody:      `{"username":"newuser","password":"password123","email":}`,
			wantStatusCode: http.StatusBadRequest,
			wantBody:       `{"message":"Invalid request body"}`,
		},
		{
			name: "Missing required fields",
			setupMock: func(_ *MockService, ml *loggo.MockLogger) {
				ml.On("Warn", "Missing required fields", mock.AnythingOfType("[]interface {}")).Return()
			},
			inputBody:      `{"username":"newuser","password":""}`,
			wantStatusCode: http.StatusBadRequest,
			wantBody:       `{"message":"Username, password, and email are required"}`,
		},
		{
			name: "User already exists",
			setupMock: func(ms *MockService, ml *loggo.MockLogger) {
				ms.On("RegisterUser", "existinguser", "password123", "existing@example.com").Return(fmt.Errorf("user already exists"))
				ml.On("Error", "Failed to register user", mock.AnythingOfType("*errors.errorString"), mock.AnythingOfType("[]interface {}")).Return()
			},
			inputBody:      `{"username":"existinguser","password":"password123","email":"existing@example.com"}`,
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       `{"message":"Failed to register user"}`,
		},
		{
			name: "Internal server error",
			setupMock: func(ms *MockService, ml *loggo.MockLogger) {
				ms.On("RegisterUser", "newuser", "password123", "newuser@example.com").Return(fmt.Errorf("database error"))
				ml.On("Error", "Failed to register user", mock.AnythingOfType("*errors.errorString"), mock.AnythingOfType("[]interface {}")).Return()
			},
			inputBody:      `{"username":"newuser","password":"password123","email":"newuser@example.com"}`,
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       `{"message":"Failed to register user"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			mockLogger := new(loggo.MockLogger)

			if tt.setupMock != nil {
				tt.setupMock(mockService, mockLogger)
			}

			handler := NewHandler(mockService, mockLogger, &config.Config{})

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.RegisterGET(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatusCode, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())

			mockService.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

func TestHandler_HandleLogin(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockService, *loggo.MockLogger)
		method         string
		username       string
		password       string
		wantStatusCode int
		wantBody       string
		wantRedirect   string
		checkSession   func(*testing.T, *http.Response, *mockSessionStore, *config.Config)
	}{
		{
			name: "Successful login",
			setupMock: func(ms *MockService, ml *loggo.MockLogger) {
				ms.On("VerifyUser", "validuser", "validpass").Return("123", nil)
				ml.On("Debug", mock.Anything, mock.Anything).Return().Twice()
			},
			method:         http.MethodPost,
			username:       "validuser",
			password:       "validpass",
			wantStatusCode: http.StatusSeeOther,
			wantRedirect:   "/campaigns",
			checkSession: func(t *testing.T, _ *http.Response, mss *mockSessionStore, cfg *config.Config) {
				sess, ok := mss.sessions[cfg.SessionName]
				assert.True(t, ok, "Session should be created")
				assert.Equal(t, "123", sess.Values["userID"])
				assert.Equal(t, "validuser", sess.Values["username"])
			},
		},
		{
			name: "Invalid credentials",
			setupMock: func(ms *MockService, ml *loggo.MockLogger) {
				ms.On("VerifyUser", "testuser", "wrongpassword").Return("", fmt.Errorf("invalid username or password"))
				ml.On("Debug", mock.Anything, mock.Anything).Return().Twice()
				ml.On("Warn", mock.Anything, mock.Anything).Return()
			},
			method:         http.MethodPost,
			username:       "testuser",
			password:       "wrongpassword",
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       "Invalid username or password",
		},
		{
			name: "Empty username",
			setupMock: func(_ *MockService, ml *loggo.MockLogger) {
				ml.On("Debug", mock.Anything, mock.Anything).Return().Twice()
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
			setupMock: func(_ *MockService, ml *loggo.MockLogger) {
				ml.On("Debug", mock.Anything, mock.Anything).Return().Twice()
				ml.On("Warn", mock.Anything, mock.Anything).Return()
			},
			method:         http.MethodPost,
			username:       "someuser",
			password:       "",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "Username and password are required",
		},
		// Session expiration
		{
			name: "Session Expired",
			setupMock: func(ms *MockService, ml *loggo.MockLogger) {
				ms.On("VerifyUser", "validuser", "validpass").Return("123", nil)
				ml.On("Debug", mock.Anything, mock.Anything).Return().Twice()
			},
			method:         http.MethodPost,
			username:       "validuser",
			password:       "validpass",
			wantStatusCode: http.StatusSeeOther,
			wantRedirect:   "/campaigns",
			checkSession: func(t *testing.T, _ *http.Response, mss *mockSessionStore, cfg *config.Config) {
				_, ok := mss.sessions[cfg.SessionName]
				assert.True(t, ok, "Session should be created")
				// Simulate session expiration
				delete(mss.sessions, cfg.SessionName)
				_, ok = mss.sessions[cfg.SessionName]
				assert.False(t, ok, "Session should be expired")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			mockLogger := new(loggo.MockLogger)
			mockSessionStore := newMockSessionStore()

			if tt.setupMock != nil {
				tt.setupMock(mockService, mockLogger)
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

			// Set up the session store
			c.Set("_session_store", mockSessionStore)

			err := handler.HandleLogin(c)

			t.Logf("Response status: %d", rec.Code)
			t.Logf("Response headers: %+v", rec.Header())
			t.Logf("Response body: %s", rec.Body.String())

			if tt.wantStatusCode == http.StatusSeeOther {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatusCode, rec.Code)
				assert.Equal(t, tt.wantRedirect, rec.Header().Get("Location"))
				if tt.checkSession != nil {
					tt.checkSession(t, rec.Result(), mockSessionStore, config)
				}
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
		})
	}
}

func TestHandler_HandleLogout(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*mockSessionStore, *loggo.MockLogger)
		wantStatusCode int
		wantRedirect   string
		wantErr        bool
	}{
		{
			name: "Successful logout",
			setupMock: func(mss *mockSessionStore, ml *loggo.MockLogger) {
				sess := sessions.NewSession(mss, testSessionName)
				sess.Values["userID"] = "123"
				sess.Values["username"] = "testuser"
				mss.sessions[testSessionName] = sess
				ml.On("Debug", "Handling logout request", mock.Anything).Return()
			},
			wantStatusCode: http.StatusSeeOther,
			wantRedirect:   "/",
			wantErr:        false,
		},
		{
			name: "Error getting session",
			setupMock: func(mss *mockSessionStore, ml *loggo.MockLogger) {
				mss.sessions = make(map[string]*sessions.Session)
				mss.getError = fmt.Errorf("failed to get session")
				ml.On("Debug", "Handling logout request", mock.Anything).Return()
				ml.On("Error", "Failed to get session", mock.AnythingOfType("*errors.errorString"), mock.Anything).Return()
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name: "Error saving session",
			setupMock: func(mss *mockSessionStore, ml *loggo.MockLogger) {
				sess := sessions.NewSession(mss, testSessionName)
				mss.sessions[testSessionName] = sess
				mss.saveError = fmt.Errorf("failed to save session")
				ml.On("Debug", "Handling logout request", mock.Anything).Return()
				ml.On("Error", "Failed to save session", mock.AnythingOfType("*errors.errorString"), mock.Anything).Return()
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSessionStore := newMockSessionStore()
			mockLogger := new(loggo.MockLogger)

			if tt.setupMock != nil {
				tt.setupMock(mockSessionStore, mockLogger)
			}

			config := &config.Config{
				SessionName: testSessionName,
			}
			handler := NewHandler(&MockService{}, mockLogger, config)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/logout", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set up the session store
			c.Set("_session_store", mockSessionStore)

			err := handler.HandleLogout(c)

			if tt.wantErr {
				assert.Error(t, err)
				if httpErr, ok := err.(*echo.HTTPError); ok {
					assert.Equal(t, tt.wantStatusCode, httpErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatusCode, rec.Code)
				assert.Equal(t, tt.wantRedirect, rec.Header().Get("Location"))

				// Check if session values are cleared
				sess, _ := mockSessionStore.Get(req, testSessionName)
				assert.Nil(t, sess.Values["userID"])
				assert.Nil(t, sess.Values["username"])
				assert.Equal(t, -1, sess.Options.MaxAge)
			}

			mockLogger.AssertExpectations(t)
		})
	}
}
