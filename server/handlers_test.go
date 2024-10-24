package server

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/fullstackdev42/mp-emailer/mocks"
	mocksEmail "github.com/fullstackdev42/mp-emailer/mocks/email"
	"github.com/fullstackdev42/mp-emailer/user"

	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRenderer is a mock of TemplateRenderer
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
	mockEmailService := mocksEmail.NewMockService(t)
	mockTemplateManager := new(MockRenderer)
	mockUserService := user.NewMockServiceInterface(t)
	mockCampaignService := campaign.NewMockServiceInterface(t)

	handler := NewHandler(mockLogger, mockEmailService, mockTemplateManager, mockUserService, mockCampaignService)

	assert.NotNil(t, handler)
	assert.IsType(t, &Handler{}, handler)
	assert.Equal(t, mockLogger, handler.Logger)
	assert.Equal(t, mockEmailService, handler.emailService)
	assert.Equal(t, mockTemplateManager, handler.templateManager)
	assert.Equal(t, mockUserService, handler.userService)
	assert.Equal(t, mockCampaignService, handler.campaignService)
	assert.NotNil(t, handler.errorHandler)
}

func TestHandler_HandleIndex(t *testing.T) {
	e := echo.New()
	mockLogger := mocks.NewMockLoggerInterface(t)
	mockCampaignService := campaign.NewMockServiceInterface(t)
	mockTemplateManager := new(MockRenderer)

	handler := &Handler{
		Logger:          mockLogger,
		campaignService: mockCampaignService,
		templateManager: mockTemplateManager,
		errorHandler:    &shared.ErrorHandler{Logger: mockLogger},
	}

	mockLogger.EXPECT().Debug(mock.Anything, mock.Anything).Times(2)
	mockCampaignService.EXPECT().GetAllCampaigns().Return([]campaign.Campaign{}, nil)
	mockTemplateManager.On("Render", mock.Anything, "home.gohtml", mock.Anything, mock.Anything).Return(nil)

	c, rec := SetupTestContext(e, "/")
	c.Set("isAuthenticated", true)

	err := handler.HandleIndex(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockCampaignService.AssertExpectations(t)
	mockTemplateManager.AssertExpectations(t)
}

func TestHandler_HandleIndex_Error(t *testing.T) {
	e := echo.New()
	mockLogger := mocks.NewMockLoggerInterface(t)
	mockCampaignService := campaign.NewMockServiceInterface(t)
	mockTemplateManager := new(MockRenderer)

	handler := &Handler{
		Logger:          mockLogger,
		campaignService: mockCampaignService,
		templateManager: mockTemplateManager,
		errorHandler:    &shared.ErrorHandler{Logger: mockLogger},
	}

	mockLogger.EXPECT().Debug(mock.Anything, mock.Anything).Times(2)
	mockLogger.EXPECT().Error(mock.Anything, mock.Anything).Times(1)
	mockCampaignService.EXPECT().GetAllCampaigns().Return(nil, errors.New("database error"))

	c, rec := SetupTestContext(e, "/")
	c.Set("isAuthenticated", true)

	err := handler.HandleIndex(c)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockCampaignService.AssertExpectations(t)
}
