package user

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// MockRenderer is a mock of echo.Renderer
type MockRenderer struct {
	LastRenderedTemplate string
}

func (m *MockRenderer) Render(_ io.Writer, name string, _ interface{}, _ echo.Context) error {
	m.LastRenderedTemplate = name
	return nil
}

// SetupTestContext sets up the common test context
func SetupTestContext(e *echo.Echo, path string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func TestNewHandler(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := mocks.NewMockLoggerInterface(t)
	mockStore := new(MockSessionStore)
	handler := NewHandler(mockRepo, mockLogger, mockStore, config.NewConfig())
	assert.NotNil(t, handler)
	assert.IsType(t, &Handler{}, handler)
	assert.Equal(t, mockRepo, handler.repo)
	assert.Equal(t, mockLogger, handler.logger)
	assert.Equal(t, mockStore, handler.Store)
}

func TestHandler_RegisterGET(t *testing.T) {
	e := echo.New()
	mockRenderer := &MockRenderer{}
	e.Renderer = mockRenderer
	c, rec := SetupTestContext(e, "/register")
	handler := &Handler{}
	err := handler.RegisterGET(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "register.html", mockRenderer.LastRenderedTemplate)
}

func TestHandler_LoginGET(t *testing.T) {
	e := echo.New()
	mockRenderer := &MockRenderer{}
	e.Renderer = mockRenderer
	c, rec := SetupTestContext(e, "/login")
	handler := &Handler{}
	err := handler.LoginGET(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "login.html", mockRenderer.LastRenderedTemplate)
}

func TestHandler_LogoutGET(t *testing.T) {
	e := echo.New()
	c, rec := SetupTestContext(e, "/logout")
	handler := &Handler{}
	err := handler.LogoutGET(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, "/", rec.Header().Get("Location"))
}
