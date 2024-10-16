package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fullstackdev42/mp-emailer/mocks"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const testSessionName = "test_session"

func TestNewHandler(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := mocks.NewMockLoggerInterface(t)
	mockStore := new(MockSessionStore)

	handler := NewHandler(mockRepo, mockLogger, mockStore)

	assert.NotNil(t, handler)
	assert.IsType(t, &Handler{}, handler)
	assert.Equal(t, mockRepo, handler.repo)
	assert.Equal(t, mockLogger, handler.logger)
	assert.Equal(t, mockStore, handler.Store)
}

func TestHandler_RegisterGET(t *testing.T) {
	t.Run("Successful GET request", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/register", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo := new(MockRepository)
		mockLogger := mocks.NewMockLoggerInterface(t)
		mockStore := new(MockSessionStore)

		// Set up expectations for the mock session store
		mockSession := sessions.NewSession(mockStore, testSessionName)
		mockStore.On("Get", req, testSessionName).Return(mockSession, nil)
		mockStore.On("Save", req, rec, mockSession).Return(nil)

		// Create the handler with the mocked dependencies
		handler := NewHandler(mockRepo, mockLogger, mockStore)

		// Perform the request
		err := handler.RegisterGET(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Check that the response body contains expected content
		assert.Contains(t, rec.Body.String(), "Register")

		mockRepo.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockStore.AssertExpectations(t)
	})
}

func TestHandler_LoginGET(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := new(MockRepository)
	mockLogger := mocks.NewMockLoggerInterface(t)
	mockStore := new(MockSessionStore)

	// Set up expectations for the mock session store
	mockSession := sessions.NewSession(mockStore, testSessionName)
	mockStore.On("Get", req, testSessionName).Return(mockSession, nil)
	mockStore.On("Save", req, rec, mockSession).Return(nil)

	handler := NewHandler(mockRepo, mockLogger, mockStore)

	// Perform the request
	err := handler.LoginGET(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Check that the response body contains expected content
	assert.Contains(t, rec.Body.String(), "Login")

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}
