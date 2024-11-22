package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	mocksSession "github.com/jonesrussell/mp-emailer/mocks/session"
)

func (s *MiddlewareTestSuite) TestSessionMiddleware() {
	// Create mock session manager
	mockSessionManager := mocksSession.NewMockManager(s.T())

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)

	// Create test handler that checks if session manager is in context
	testHandler := func(c echo.Context) error {
		sessionManager := c.Get("session_manager")
		assert.NotNil(s.T(), sessionManager)
		assert.Equal(s.T(), mockSessionManager, sessionManager)
		return nil
	}

	// Create and execute middleware
	handler := s.manager.SessionMiddleware(mockSessionManager)(testHandler)
	err := handler(c)

	// Assert
	assert.NoError(s.T(), err)
}
