package user_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/mp-emailer/internal/testutil"
	"github.com/jonesrussell/mp-emailer/shared"
	"github.com/jonesrussell/mp-emailer/user"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	mocksSession "github.com/jonesrussell/mp-emailer/mocks/session"
	mocksShared "github.com/jonesrussell/mp-emailer/mocks/shared"
	mocksUser "github.com/jonesrussell/mp-emailer/mocks/user"
)

type HandlerTestSuite struct {
	testutil.BaseTestSuite
	handler          *user.Handler
	UserService      *mocksUser.MockServiceInterface
	UserRepo         *mocksUser.MockRepositoryInterface
	SessionManager   *mocksSession.MockManager
	StoreProvider    *mocksSession.MockStoreProvider
	TemplateRenderer *mocksShared.MockTemplateRendererInterface
}

func (s *HandlerTestSuite) SetupTest() {
	s.BaseTestSuite.SetupTest()

	s.UserService = mocksUser.NewMockServiceInterface(s.T())
	s.UserRepo = mocksUser.NewMockRepositoryInterface(s.T())
	s.SessionManager = mocksSession.NewMockManager(s.T())
	s.StoreProvider = mocksSession.NewMockStoreProvider(s.T())
	s.ErrorHandler = mocksShared.NewMockErrorHandlerInterface(s.T())
	s.TemplateRenderer = mocksShared.NewMockTemplateRendererInterface(s.T())

	s.Config.Auth.SessionName = "test_session"
	s.Echo.Renderer = s.TemplateRenderer

	params := user.HandlerParams{
		BaseHandlerParams: shared.BaseHandlerParams{
			Logger:           s.Logger,
			ErrorHandler:     s.ErrorHandler,
			TemplateRenderer: s.TemplateRenderer,
			Config:           s.Config,
			StoreProvider:    s.StoreProvider,
		},
		Service: s.UserService,
		Repo:    s.UserRepo,
	}

	s.handler = user.NewHandler(params)
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) TestLoginPOST_Success() {
	// Create test login params
	loginParams := &user.LoginDTO{
		Username: "testuser",
		Password: "password123",
	}

	// Create expected user response
	expectedUser := &user.User{
		Username: loginParams.Username,
		// Add other required User fields here
	}

	// Setup mock session
	mockSession := mocksSession.NewMockInterface(s.T())
	mockValues := make(map[interface{}]interface{})
	mockSession.On("Values").Return(mockValues)
	mockSession.On("IsNew").Return(false)
	mockSession.On("Options").Return(&sessions.Options{})
	mockSession.On("GetID").Return("test-session-id")
	mockSession.On("Save", mock.AnythingOfType("*http.Request"), mock.AnythingOfType("http.ResponseWriter")).Return(nil)

	// Setup session manager expectations
	s.handler.SessionManager = s.SessionManager // Important: Set the session manager
	s.SessionManager.On("GetSession", mock.AnythingOfType("*echo.Context"), s.Config.Auth.SessionName).
		Return(mockSession, nil)
	s.SessionManager.On("AddFlash", mockSession, "Successfully logged in!").Return(nil)
	s.SessionManager.On("SetAuthenticated", mockSession, true).Return(nil)
	s.SessionManager.On("SetSessionValues", mockSession, expectedUser).Return(nil)
	s.SessionManager.On("SaveSession", mock.AnythingOfType("*echo.Context"), mockSession).Return(nil)

	// Setup service expectations
	s.UserService.On("AuthenticateUser",
		mock.Anything,
		loginParams.Username,
		loginParams.Password,
	).Return(expectedUser, nil)

	// Setup logger expectations
	s.Logger.On("Debug", "LoginPOST: Starting").Return()
	s.Logger.On("Debug", "LoginPOST: Authentication successful", "username", loginParams.Username).Return()
	// Add error logger expectation in case session fails
	s.Logger.On("Error", "Error getting session", mock.AnythingOfType("*echo.HTTPError")).Return()

	// Create form data
	form := url.Values{}
	form.Add("username", loginParams.Username)
	form.Add("password", loginParams.Password)

	// Create request with form data
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	// Create context
	rec := httptest.NewRecorder()
	c := s.Echo.NewContext(req, rec)

	// Execute request
	err := s.handler.LoginPOST(c)

	// Assertions
	s.NoError(err)
	s.Equal(http.StatusSeeOther, rec.Code) // Expecting redirect after successful login

	// Verify all expectations were met
	s.UserService.AssertExpectations(s.T())
	s.SessionManager.AssertExpectations(s.T())
	mockSession.AssertExpectations(s.T())
	s.Logger.AssertExpectations(s.T())
}
