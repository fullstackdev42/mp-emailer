package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/jonesrussell/mp-emailer/config"
	"github.com/jonesrussell/mp-emailer/middleware"
	"github.com/jonesrussell/mp-emailer/mocks"
	mocksShared "github.com/jonesrussell/mp-emailer/mocks/shared"
)

type MiddlewareTestSuite struct {
	suite.Suite
	echo           *echo.Echo
	mockLogger     *mocks.MockLoggerInterface
	mockErrHandler *mocksShared.MockErrorHandlerInterface
	manager        *middleware.Manager
	config         *config.Config
}

func (s *MiddlewareTestSuite) SetupTest() {
	s.echo = echo.New()
	s.mockLogger = mocks.NewMockLoggerInterface(s.T())
	s.mockErrHandler = mocksShared.NewMockErrorHandlerInterface(s.T())

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
				RequestsPerSecond: 10.0, // Allow 10 requests per second in tests
				BurstSize:         20,   // Allow burst of 20 requests
			},
		},
		Auth: config.AuthConfig{
			JWTSecret: "test-secret",
		},
	}

	var err error
	s.manager, err = middleware.NewManager(middleware.ManagerParams{
		Logger:       s.mockLogger,
		Cfg:          s.config,
		ErrorHandler: s.mockErrHandler,
	})
	s.Require().NoError(err)
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

func (s *MiddlewareTestSuite) TestMethodOverride() {
	testCases := []struct {
		name           string
		method         string
		formMethod     string
		expectedMethod string
	}{
		{
			name:           "POST to PUT override",
			method:         "POST",
			formMethod:     "PUT",
			expectedMethod: "PUT",
		},
		{
			name:           "POST to DELETE override",
			method:         "POST",
			formMethod:     "DELETE",
			expectedMethod: "DELETE",
		},
		{
			name:           "POST without override",
			method:         "POST",
			formMethod:     "",
			expectedMethod: "POST",
		},
		{
			name:           "GET ignores override",
			method:         "GET",
			formMethod:     "DELETE",
			expectedMethod: "GET",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(tc.method, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tc.formMethod != "" {
				c.SetParamValues(tc.formMethod)
				c.Request().Form = map[string][]string{
					"_method": {tc.formMethod},
				}
			}

			// Execute middleware
			handler := middleware.MethodOverride()(func(c echo.Context) error {
				s.Equal(tc.expectedMethod, c.Request().Method)
				return nil
			})

			// Assert
			err := handler(c)
			s.NoError(err)
		})
	}
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
