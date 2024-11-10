package middleware_test

import (
	"encoding/gob"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	"github.com/fullstackdev42/mp-emailer/middleware"
	"github.com/fullstackdev42/mp-emailer/mocks"
	mocksMiddleware "github.com/fullstackdev42/mp-emailer/mocks/middleware"
)

type MiddlewareTestSuite struct {
	suite.Suite
	echo       *echo.Echo
	mockLogger *mocks.MockLoggerInterface
	mockStore  *mocksMiddleware.MockSessionStore
	manager    *middleware.Manager
}

func (s *MiddlewareTestSuite) SetupTest() {
	s.echo = echo.New()
	s.mockLogger = mocks.NewMockLoggerInterface(s.T())
	s.mockStore = mocksMiddleware.NewMockSessionStore(s.T())

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
