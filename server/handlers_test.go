package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/email"
	"github.com/fullstackdev42/mp-emailer/mocks"
	mocksEmail "github.com/fullstackdev42/mp-emailer/mocks/email"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRenderer is a mock of echo.Renderer
type MockRenderer struct {
	mock.Mock
}

func (m *MockRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	args := m.Called(w, name, data, c)
	return args.Error(0)
}

// SetupTestContext sets up the common test context
func SetupTestContext(e *echo.Echo, path string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func TestNewHandler(t *testing.T) {
	mockLogger := mocks.NewMockLoggerInterface(t)
	mockEmailService := &mocksEmail.MockService{}
	mockTemplateManager := &TemplateManager{}
	mockUserService := new(user.Service)
	mockCampaignService := new(campaign.Service)

	handler := NewHandler(
		mockLogger,
		mockEmailService,
		mockTemplateManager,
		mockUserService,
		mockCampaignService,
	)

	assert.NotNil(t, handler)
	assert.IsType(t, &Handler{}, handler)
	assert.Equal(t, mockLogger, handler.Logger)
	assert.Equal(t, mockEmailService, handler.emailService)
	assert.Equal(t, mockTemplateManager, handler.templateManager)
	assert.Equal(t, mockUserService, handler.userService)
	assert.Equal(t, mockCampaignService, handler.campaignService)
}

func TestHandler_HandleIndex(t *testing.T) {
	e := echo.New()
	mockRenderer := new(MockRenderer)
	e.Renderer = mockRenderer

	c, rec := SetupTestContext(e, "/")

	handler := &Handler{
		Logger:          mocks.NewMockLoggerInterface(t),
		Store:           sessions.NewCookieStore([]byte("test-secret")),
		emailService:    *new(email.Service),
		templateManager: &TemplateManager{},
		userService:     new(user.Service),
		campaignService: new(campaign.Service),
	}

	mockRenderer.On("Render", mock.Anything, "index.html", nil, mock.Anything).Return(nil)

	err := handler.HandleIndex(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockRenderer.AssertExpectations(t)
}
