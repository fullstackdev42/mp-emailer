package user

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/mocks"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testSessionName = "test_session"

func TestNewHandler(t *testing.T) {
	mockService := NewMockServiceInterface(t)
	mockLogger := mocks.NewMockLoggerInterface(t)
	mockConfig := &config.Config{}

	handler := NewHandler(mockService, mockLogger, mockConfig)

	assert.NotNil(t, handler)
	assert.IsType(t, &Handler{}, handler)
	assert.Equal(t, mockService, handler.service)
	assert.Equal(t, mockConfig, handler.config)
}

func TestHandler_RegisterPOST(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockServiceInterface, *mocks.MockLoggerInterface, *MockSessionStore)
		username       string
		email          string
		password       string
		wantStatusCode int
		wantRedirect   string
	}{
		{
			name: "Successful registration",
			setupMock: func(ms *MockServiceInterface, ml *mocks.MockLoggerInterface, mss *MockSessionStore) {
				ms.EXPECT().RegisterUser("testuser", "test@example.com", "password123").Return(nil)
				ml.EXPECT().Info("User registered successfully", "username", "testuser").Times(1)

				// Setup session store mock
				session := sessions.NewSession(mss, "test_session")
				mss.On("Get", mock.Anything, "test_session").Return(session, nil)
				mss.On("Save", mock.Anything, mock.Anything, session).Return(nil)
			},
			username:       "testuser",
			email:          "test@example.com",
			password:       "password123",
			wantStatusCode: http.StatusSeeOther,
			wantRedirect:   "/login",
		},
		// {
		// 	name: "Missing required fields",
		// 	setupMock: func(_ *MockServiceInterface, ml *mocks.MockLoggerInterface) {
		// 		ml.EXPECT().Warn("Missing required fields in registration form").Return()
		// 	},
		// 	inputBody:      `username=newuser&password=`,
		// 	wantStatusCode: http.StatusSeeOther,
		// 	wantRedirect:   "/register",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockServiceInterface(t)
			mockLogger := mocks.NewMockLoggerInterface(t)
			mockSessionStore := new(MockSessionStore)
			tt.setupMock(mockService, mockLogger, mockSessionStore)

			handler := NewHandler(mockService, mockLogger, &config.Config{SessionName: "test_session"})

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(url.Values{
				"username": {tt.username},
				"email":    {tt.email},
				"password": {tt.password},
			}.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set the session store in the context
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
	t.Run("Successful login", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(url.Values{
			"username": {"testuser"},
			"password": {"testpass"},
		}.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService := NewMockServiceInterface(t)
		mockLogger := mocks.NewMockLoggerInterface(t)
		mockConfig := &config.Config{SessionName: "test_session"}

		handler := NewHandler(mockService, mockLogger, mockConfig)

		// Expectations
		mockService.EXPECT().VerifyUser("testuser", "testpass").Return("123", nil)
		mockLogger.EXPECT().Debug(mock.Anything, mock.Anything).Return()
		mockLogger.EXPECT().Warn(mock.Anything, mock.Anything).Maybe()
		mockLogger.EXPECT().Error(mock.Anything, mock.Anything).Maybe()

		// Create a mock session store
		mockStore := sessions.NewCookieStore([]byte("secret"))
		e.Use(session.Middleware(mockStore))

		// Perform the request
		err := handler.LoginPOST(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.Equal(t, "/campaigns", rec.Header().Get("Location"))

		// Check session
		sess, _ := session.Get("test_session", c)
		assert.True(t, sess.Values["authenticated"].(bool))
		assert.Equal(t, "123", sess.Values["userID"])
		assert.Equal(t, "testuser", sess.Values["username"])
	})
}

func TestHandler_RegisterGET(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*mocks.MockLoggerInterface)
		wantStatusCode int
		wantErr        bool
	}{
		{
			name: "Successful GET request",
			setupMock: func(ml *mocks.MockLoggerInterface) {
				ml.EXPECT().Debug("RegisterGET: Starting").Return()
			},
			wantStatusCode: http.StatusOK,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockServiceInterface(t)
			mockLogger := *mocks.NewMockLoggerInterface(t)

			tt.setupMock(&mockLogger)

			handler := NewHandler(mockService, &mockLogger, &config.Config{SessionName: testSessionName})

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

func TestHandler_LoginGET(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := NewMockServiceInterface(t)
	mockLogger := mocks.NewMockLoggerInterface(t)
	mockConfig := &config.Config{SessionName: "test_session"}

	handler := NewHandler(mockService, mockLogger, mockConfig)

	// Perform the request
	err := handler.LoginGET(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Check that the response body contains expected content
	// This assumes that your login page contains the text "Login"
	assert.Contains(t, rec.Body.String(), "Login")

	mockLogger.AssertExpectations(t)
}
