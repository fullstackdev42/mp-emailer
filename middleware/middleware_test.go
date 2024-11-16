package middleware_test

import (
	"encoding/gob"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/middleware"
	"github.com/fullstackdev42/mp-emailer/mocks"
	mocksMiddleware "github.com/fullstackdev42/mp-emailer/mocks/middleware"
	mocksShared "github.com/fullstackdev42/mp-emailer/mocks/shared"
	"github.com/fullstackdev42/mp-emailer/shared"
)

type MiddlewareTestSuite struct {
	suite.Suite
	echo           *echo.Echo
	mockLogger     *mocks.MockLoggerInterface
	mockStore      *mocksMiddleware.MockSessionStore
	mockErrHandler *mocksShared.MockErrorHandlerInterface
	manager        *middleware.Manager
	config         *config.Config
}

func (s *MiddlewareTestSuite) SetupTest() {
	s.echo = echo.New()
	s.mockLogger = mocks.NewMockLoggerInterface(s.T())
	s.mockStore = mocksMiddleware.NewMockSessionStore(s.T())
	s.mockErrHandler = mocksShared.NewMockErrorHandlerInterface(s.T())

	// Add this line to register UUID type with gob
	gob.Register(uuid.UUID{})

	// Create mock dependencies
	s.config = &config.Config{
		Auth: config.AuthConfig{
			SessionName: "test_session",
		},
		App: config.AppConfig{
			Env: "test",
		},
	}

	var err error
	s.manager, err = middleware.NewManager(middleware.ManagerParams{
		SessionStore: s.mockStore,
		Logger:       s.mockLogger,
		Cfg:          s.config,
		ErrorHandler: s.mockErrHandler,
	})
	s.Require().NoError(err)
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

func (s *MiddlewareTestSuite) TestGetUserIDFromSession() {
	s.Run("valid UUID in session", func() {
		// Set up all expected logger calls in the correct order
		s.mockLogger.EXPECT().
			Debug("Attempting to get userID from session", "session_name", "test-session")

		s.mockLogger.EXPECT().
			Debug("Getting session", "session_name", "test-session")

		s.mockLogger.EXPECT().
			Debug("Session retrieved successfully",
				"session_id", "", // The session ID might be empty in test
				"is_new", false,
				"values_count", 1)

		s.mockLogger.EXPECT().
			Debug("User ID (UUID) found in session",
				"user_id", "e513302d-4563-47c4-932f-d22af5c07e62",
				"session_name", "test-session")

		// Setup test session
		userID := uuid.MustParse("e513302d-4563-47c4-932f-d22af5c07e62")
		session := sessions.NewSession(s.mockStore, "test-session")
		session.Values["user_id"] = userID

		// Create test context
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Mock the session store Get method
		s.mockStore.EXPECT().
			Get(req, "test-session").
			Return(session, nil)

		// Test the method
		result, err := s.manager.GetUserIDFromSession(c, "test-session")

		s.NoError(err)
		s.Equal(userID.String(), result)
	})
}

// Test session middleware
func (s *MiddlewareTestSuite) TestSessionsMiddleware() {
	s.Run("successful session handling", func() {
		// Setup
		mockStore := new(mocksMiddleware.MockSessionStore)
		mockLogger := new(mocks.MockLoggerInterface)
		mockErrHandler := new(mocksShared.MockErrorHandlerInterface)

		session := sessions.NewSession(mockStore, "test-session")

		mockStore.On("Get", mock.Anything, "test-session").Return(session, nil)
		mockStore.On("Save", mock.Anything, mock.Anything, session).Return(nil)
		mockLogger.On("Debug", "Session middleware processing request", "path", mock.Anything).Return()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		middleware := middleware.NewSessionsMiddleware(
			mockStore,
			mockLogger,
			"test-session",
			mockErrHandler,
		)

		// Execute
		err := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		})(c)

		// Assert
		s.NoError(err)
		s.Equal(http.StatusOK, rec.Code)
		mockStore.AssertExpectations(s.T())
		mockLogger.AssertExpectations(s.T())
	})

	s.Run("session error", func() {
		// Setup
		mockStore := new(mocksMiddleware.MockSessionStore)
		mockLogger := new(mocks.MockLoggerInterface)
		mockErrHandler := new(mocksShared.MockErrorHandlerInterface)

		expectedErr := errors.New("session error")
		mockLogger.On("Debug", "Session middleware processing request", "path", mock.Anything).Return()
		mockStore.On("Get", mock.Anything, "test-session").Return(nil, expectedErr)

		mockErrHandler.On("HandleHTTPError",
			mock.Anything,
			expectedErr,
			"Error getting session",
			http.StatusInternalServerError,
		).Return(echo.NewHTTPError(http.StatusInternalServerError))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		middleware := middleware.NewSessionsMiddleware(
			mockStore,
			mockLogger,
			"test-session",
			mockErrHandler,
		)

		// Execute
		handlerCalled := false
		err := middleware(func(c echo.Context) error {
			handlerCalled = true
			return c.String(http.StatusOK, "success")
		})(c)

		// Assert
		s.Error(err)
		s.Equal(http.StatusInternalServerError, err.(*echo.HTTPError).Code)
		s.False(handlerCalled, "Handler should not be called when session error occurs")
		mockStore.AssertExpectations(s.T())
		mockLogger.AssertExpectations(s.T())
	})

	s.Run("session save error", func() {
		// Setup
		mockStore := new(mocksMiddleware.MockSessionStore)
		mockLogger := new(mocks.MockLoggerInterface)
		mockErrHandler := new(mocksShared.MockErrorHandlerInterface)

		session := sessions.NewSession(mockStore, "test-session")
		expectedErr := errors.New("save error")

		mockStore.On("Get", mock.Anything, "test-session").Return(session, nil)
		mockStore.On("Save", mock.Anything, mock.Anything, session).Return(expectedErr)
		mockLogger.On("Debug", "Session middleware processing request", "path", mock.Anything).Return()

		mockErrHandler.On("HandleHTTPError",
			mock.Anything,
			expectedErr,
			"Error saving session",
			http.StatusInternalServerError,
		).Return(echo.NewHTTPError(http.StatusInternalServerError))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		middleware := middleware.NewSessionsMiddleware(
			mockStore,
			mockLogger,
			"test-session",
			mockErrHandler,
		)

		// Execute
		handlerCalled := false
		err := middleware(func(c echo.Context) error {
			handlerCalled = true
			return c.String(http.StatusOK, "success")
		})(c)

		// Assert
		s.Error(err)
		s.Equal(http.StatusInternalServerError, err.(*echo.HTTPError).Code)
		s.True(handlerCalled, "Handler should be called before session save error")
		mockStore.AssertExpectations(s.T())
		mockLogger.AssertExpectations(s.T())
	})
}

// Test JWT middleware
func (s *MiddlewareTestSuite) TestJWTMiddleware() {
	testCases := []struct {
		name         string
		setupHeader  func(*http.Request)
		expectError  bool
		errorMessage string
	}{
		{
			name: "valid JWT token",
			setupHeader: func(r *http.Request) {
				// Generate token with correct parameters
				token, _ := shared.GenerateToken("testuser", s.config.Auth.JWTSecret, 60) // 60 minutes expiration
				r.Header.Set("Authorization", "Bearer "+token)
			},
			expectError: false,
		},
		{
			name:         "missing authorization header",
			setupHeader:  func(_ *http.Request) {},
			expectError:  true,
			errorMessage: "Missing authorization header",
		},
		{
			name: "invalid authorization format",
			setupHeader: func(r *http.Request) {
				r.Header.Set("Authorization", "InvalidFormat token")
			},
			expectError:  true,
			errorMessage: "Invalid authorization header",
		},
		{
			name: "invalid token",
			setupHeader: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid.token.here")
			},
			expectError:  true,
			errorMessage: "Invalid token",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tc.setupHeader(req)

			// Execute middleware
			handler := s.manager.JWTMiddleware()(func(_ echo.Context) error {
				return nil
			})

			err := handler(c)

			// Assertions
			if tc.expectError {
				s.Error(err)
				httpErr, ok := err.(*echo.HTTPError)
				s.True(ok)
				s.Equal(http.StatusUnauthorized, httpErr.Code)
				s.Contains(httpErr.Message.(map[string]string)["error"], tc.errorMessage)
			} else {
				s.NoError(err)
				// Verify user info was set in context
				userID := c.Get("user_id")
				s.NotNil(userID)
				s.Equal("testuser", userID)
			}
		})
	}
}
