package user

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/mocks"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
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
	mockRepo := NewMockRepositoryInterface(t)
	mockService := new(Service)
	mockLogger := mocks.NewMockLoggerInterface(t)
	mockStore := sessions.NewCookieStore([]byte("test-secret"))
	mockConfig := config.NewConfig()
	handler := NewHandler(mockRepo, mockService, mockLogger, mockStore, mockConfig)

	assert.NotNil(t, handler)
	assert.IsType(t, &Handler{}, handler)
	assert.Equal(t, mockRepo, handler.repo)
	assert.Equal(t, mockLogger, handler.Logger)
	assert.Equal(t, mockStore, handler.Store)
	assert.Equal(t, mockConfig, handler.Config)
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
	assert.Equal(t, "register.gohtml", mockRenderer.LastRenderedTemplate)
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
	assert.Equal(t, "login.gohtml", mockRenderer.LastRenderedTemplate)
}

func TestHandler_LogoutGET(t *testing.T) {
	mockStore := sessions.NewCookieStore([]byte("test-secret"))
	h := &Handler{
		Store:       mockStore,
		SessionName: "test-session",
		Logger:      mocks.NewMockLoggerInterface(t),
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.LogoutGET(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, "/", rec.Header().Get("Location"))
}

func TestHandler_RegisterPOST(t *testing.T) {
	mockRepo := NewMockRepositoryInterface(t)
	mockService := NewMockServiceInterface(t)
	mockLogger := mocks.NewMockLoggerInterface(t)
	mockStore := sessions.NewCookieStore([]byte("test-secret"))
	mockConfig := &config.Config{SessionName: "test-session"}
	h := NewHandler(mockRepo, mockService, mockLogger, mockStore, mockConfig)

	mockRepo.EXPECT().UserExists(mock.Anything, mock.Anything).Return(false, nil)
	mockRepo.EXPECT().CreateUser(mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockRepo.EXPECT().GetUserByUsername(mock.Anything).Return(&User{ID: "1", Username: "testuser"}, nil)
	mockService.EXPECT().RegisterUser(mock.AnythingOfType("user.RegisterUserParams")).Return(nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader("username=testuser&email=test@example.com&password=testpass"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.RegisterPOST(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, "/", rec.Header().Get("Location"))
	mockRepo.AssertExpectations(t)
	mockService.AssertExpectations(t)
}
