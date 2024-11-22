package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/mp-emailer/config"
	"github.com/jonesrussell/mp-emailer/middleware"
	mocksLogger "github.com/jonesrussell/mp-emailer/mocks/logger"
	mocksSession "github.com/jonesrussell/mp-emailer/mocks/session"
	mocksShared "github.com/jonesrussell/mp-emailer/mocks/shared"
)

type MiddlewareTestSuite struct {
	suite.Suite
	echo               *echo.Echo
	mockLogger         *mocksLogger.MockInterface
	mockErrHandler     *mocksShared.MockErrorHandlerInterface
	manager            *middleware.Manager
	config             *config.Config
	mockSessionManager *mocksSession.MockManager
}

func (s *MiddlewareTestSuite) SetupTest() {
	s.echo = echo.New()
	s.mockLogger = mocksLogger.NewMockInterface(s.T())
	s.mockErrHandler = mocksShared.NewMockErrorHandlerInterface(s.T())
	s.mockSessionManager = mocksSession.NewMockManager(s.T())

	// Create mock dependencies with rate limiting configuration
	s.config = &config.Config{
		App: config.AppConfig{
			Env: "test",
		},
		Server: config.ServerConfig{
			RateLimiting: struct {
				RequestsPerSecond float64 `yaml:"requests_per_second" env:"RATE_LIMIT_RPS" envDefault:"20"`
				BurstSize         int     `yaml:"burst_size" env:"RATE_LIMIT_BURST" envDefault:"50"`
			}{
				RequestsPerSecond: 10.0,
				BurstSize:         20,
			},
		},
		Auth: config.AuthConfig{
			JWTSecret: "test-secret",
		},
	}

	var err error
	s.manager, err = middleware.NewManager(middleware.ManagerParams{
		Logger:         s.mockLogger,
		Cfg:            s.config,
		ErrorHandler:   s.mockErrHandler,
		SessionManager: s.mockSessionManager,
	})
	s.Require().NoError(err)
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

func (s *MiddlewareTestSuite) TestRegisterMiddleware() {
	// Configure error handler mock for rate limiter
	s.mockErrHandler.On("HandleHTTPError",
		mock.Anything,
		mock.Anything,
		"rate_limit",
		http.StatusTooManyRequests).
		Return(nil).
		Maybe()

	// Configure session manager expectations
	session := &sessions.Session{
		Values: make(map[interface{}]interface{}),
	}

	// Mock GetSession
	s.mockSessionManager.On("GetSession",
		mock.Anything,
		mock.Anything).
		Return(session, nil).
		Maybe()

	// Mock GetFlashes
	s.mockSessionManager.On("GetFlashes",
		mock.AnythingOfType("*sessions.Session")).
		Return([]interface{}{}, nil).
		Maybe()

	// Mock SaveSession
	s.mockSessionManager.On("SaveSession",
		mock.Anything,
		mock.AnythingOfType("*sessions.Session")).
		Return(nil).
		Maybe()

	// Configure logger expectations for both session debug calls
	s.mockLogger.On("Debug",
		"Session state",
		"session_id", mock.Anything,
		"is_new", mock.Anything,
		"user_id", mock.Anything,
		"is_authenticated", mock.Anything).
		Return().
		Maybe()

	s.mockLogger.On("Debug",
		"Session state after handler",
		"session_id", mock.Anything,
		"user_id", mock.Anything,
		"is_authenticated", mock.Anything).
		Return().
		Maybe()

	// Test that all middleware are registered without error
	s.manager.Register(s.echo)

	// Create a test request with required headers
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.Header.Set("Host", "test-host")
	rec := httptest.NewRecorder()

	// Add a handler that we can use to verify middleware execution
	handlerCalled := false
	s.echo.GET("/", func(c echo.Context) error {
		handlerCalled = true
		return c.String(http.StatusOK, "OK")
	})

	// Execute the request through the Echo instance
	s.echo.ServeHTTP(rec, req)

	// Check if we got a successful response
	s.True(rec.Code == http.StatusOK || rec.Code == http.StatusTooManyRequests,
		"Expected either 200 OK or 429 Too Many Requests")

	// Only verify handler was called if we got a 200 OK
	if rec.Code == http.StatusOK {
		s.True(handlerCalled, "Expected handler to be called through middleware chain")
	}
}

func (s *MiddlewareTestSuite) TestSessionMiddleware() {
	// Configure session manager expectations
	session := &sessions.Session{
		Values: make(map[interface{}]interface{}),
	}

	// Mock GetSession
	s.mockSessionManager.On("GetSession",
		mock.Anything,
		mock.Anything).
		Return(session, nil)

	// Mock GetFlashes
	s.mockSessionManager.On("GetFlashes",
		mock.AnythingOfType("*sessions.Session")).
		Return([]interface{}{}, nil)

	// Mock SaveSession
	s.mockSessionManager.On("SaveSession",
		mock.Anything,
		mock.AnythingOfType("*sessions.Session")).
		Return(nil)

	// Configure logger expectations for both session debug calls
	s.mockLogger.On("Debug",
		"Session state",
		"session_id", mock.Anything,
		"is_new", mock.Anything,
		"user_id", mock.Anything,
		"is_authenticated", mock.Anything).
		Return()

	s.mockLogger.On("Debug",
		"Session state after handler",
		"session_id", mock.Anything,
		"user_id", mock.Anything,
		"is_authenticated", mock.Anything).
		Return()

	// Create test request and response
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)

	// Create test handler
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}

	// Create middleware
	middleware := s.manager.SessionMiddleware(s.mockSessionManager)

	// Execute middleware
	err := middleware(handler)(c)

	// Assert
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
}
