package user

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/mocks"
	"github.com/gorilla/sessions"
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
	assert.Equal(t, mockLogger, handler.Logger)
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
	// Create a mock session store
	mockStore := sessions.NewCookieStore([]byte("test-secret"))

	// Initialize the handler with the mock store
	h := &Handler{
		Store:       mockStore,
		SessionName: "test-session",
		Logger:      mocks.NewMockLoggerInterface(t),
	}

	// Create a new echo context for testing
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the LogoutGET function
	err := h.LogoutGET(c)

	// Assert that there's no error
	assert.NoError(t, err)

	// Assert that the response is a redirect (303 See Other)
	assert.Equal(t, http.StatusSeeOther, rec.Code)

	// Assert that the Location header is set to "/"
	assert.Equal(t, "/", rec.Header().Get("Location"))
}
