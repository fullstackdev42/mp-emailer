package middleware_test

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/fullstackdev42/mp-emailer/middleware"
	"github.com/fullstackdev42/mp-emailer/mocks"
)

type MiddlewareTestSuite struct {
	suite.Suite
	echo       *echo.Echo
	mockLogger *mocks.MockLoggerInterface
	mockStore  *sessions.CookieStore
	manager    *middleware.Manager
}

func (s *MiddlewareTestSuite) SetupTest() {
	s.echo = echo.New()
	s.mockLogger = mocks.NewMockLoggerInterface(s.T())
	s.mockStore = sessions.NewCookieStore([]byte("test-key"))

	// Add this line to register UUID type with gob
	gob.Register(uuid.UUID{})

	s.manager = middleware.NewManager(middleware.ManagerParams{
		SessionStore: s.mockStore,
		Logger:       s.mockLogger,
	})
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}
func (s *MiddlewareTestSuite) TestGetUserIDFromSession() {
	testCases := []struct {
		name          string
		setupSession  func(*sessions.Session)
		setupMocks    func()
		expectedError error
	}{
		{
			name: "valid UUID in session",
			setupSession: func(session *sessions.Session) {
				testUUID := uuid.New()
				session.Values["user_id"] = testUUID
				fmt.Printf("Setting session user_id: %v\n", testUUID)
			},
			setupMocks: func() {
				s.mockLogger.On("Debug", "Attempting to get userID from session", "session_name", "test-session").Return()
				s.mockLogger.On("Debug", "Getting session", "session_name", "test-session").Return()
				s.mockLogger.On("Debug", "Session retrieved successfully", "session_id", mock.Anything, "is_new", mock.Anything, "values_count", mock.Anything).Return()
				s.mockLogger.On("Debug", "Raw user_id value", "type", "<nil>").Return()
				s.mockLogger.On("Debug", "Unexpected type for user_id", "type", "<nil>").Return()
			},
			expectedError: middleware.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			fmt.Printf("Running test case: %s\n", tc.name)

			// Setup
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)

			// Create and setup session
			session, _ := s.mockStore.New(req, "test-session")
			if tc.setupSession != nil {
				tc.setupSession(session)
				err := session.Save(req, rec)
				s.NoError(err, "Failed to save session")
			}

			// Set store in context
			c.Set("store", s.mockStore)

			// Setup mocks
			tc.setupMocks()

			// Execute
			userID, err := s.manager.GetUserIDFromSession(c, "test-session")

			// Assert
			if tc.expectedError != nil {
				s.Error(err)
				s.Equal(tc.expectedError, err)
			} else {
				s.NoError(err)
				s.NotEqual(uuid.Nil, userID)
			}

			s.mockLogger.AssertExpectations(s.T())
		})
	}
}
