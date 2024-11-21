package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
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

	// Create mock dependencies
	s.config = &config.Config{
		App: config.AppConfig{
			Env: "test",
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
	// Test that all middleware are registered without error
	s.manager.Register(s.echo)

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Add a handler that we can use to verify middleware execution
	handlerCalled := false
	handler := func(_ echo.Context) error {
		handlerCalled = true
		return nil
	}

	// Create a test route
	s.echo.GET("/", handler)

	// Execute the request through the Echo instance
	s.echo.ServeHTTP(rec, req)

	// Verify that the handler was called (indicating middleware chain is working)
	s.True(handlerCalled, "Expected handler to be called through middleware chain")

	// Verify response headers that indicate middleware was active
	// Logger middleware adds Server header
	s.Contains(rec.Header().Get("Server"), "Echo")
}
