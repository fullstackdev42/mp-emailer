package user_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
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
	mockStore        *mocksSession.MockStore
	mockSession      *mocksSession.MockInterface
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

func (s *HandlerTestSuite) setupBaseMocks() {
	// Setup session mocks
	s.mockStore = mocksSession.NewMockStore(s.T())
	s.mockSession = mocksSession.NewMockInterface(s.T())

	// Setup mock session behavior
	mockValues := make(map[interface{}]interface{})
	s.mockSession.On("Values").Return(mockValues)
	s.mockSession.On("IsNew").Return(false)
	s.mockSession.On("Options").Return(&sessions.Options{})
	s.mockSession.On("GetID").Return("test-session-id")

	// Setup store provider expectations
	s.StoreProvider.On("GetStore", mock.AnythingOfType("*http.Request")).Return(s.mockStore)

	// Setup store expectations
	s.mockStore.On("New",
		mock.AnythingOfType("*http.Request"),
		s.Config.Auth.SessionName,
	).Return(s.mockSession, nil)

	// Setup session manager expectations
	s.SessionManager.On("GetSession",
		mock.MatchedBy(func(_ echo.Context) bool { return true }),
		s.Config.Auth.SessionName,
	).Return(s.mockSession, nil)

	// Set the session manager on the handler
	s.handler.SessionManager = s.SessionManager
}

func (s *HandlerTestSuite) createTestRequest(method, path, payload string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, path, strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return s.Echo.NewContext(req, rec), rec
}

func (s *HandlerTestSuite) TestLoginPOST_Success() {
	// Setup base mocks
	s.setupBaseMocks()

	mockUser := &user.User{
		BaseModel: shared.BaseModel{ID: uuid.New()},
		Username:  "testuser",
	}

	// Setup additional expectations specific to this test
	s.SessionManager.On("AddFlash", s.mockSession, "Successfully logged in!").Return()
	s.SessionManager.On("SetAuthenticated", mock.AnythingOfType("*echo.Context"), true).Return(nil)
	s.SessionManager.On("SetSessionValues", s.mockSession, mockUser).Return()

	// Setup service expectations
	s.UserService.On("AuthenticateUser", mock.Anything, "testuser", "testpass").Return(mockUser, nil)

	// Setup error handler expectations
	s.ErrorHandler.On("HandleHTTPError",
		mock.AnythingOfType("*echo.context"),
		mock.AnythingOfType("*echo.HTTPError"),
		"Error getting session",
		http.StatusInternalServerError,
	).Return(nil)

	// Setup logger expectations
	s.Logger.On("Error", "Error getting session", mock.AnythingOfType("*echo.HTTPError")).Return()
	s.Logger.On("Debug", "Adding flash message", "message", "Successfully logged in!").Return()

	// Execute test
	payload := `{"username":"testuser","password":"testpass"}`
	c, rec := s.createTestRequest(http.MethodPost, "/login", payload)

	err := s.handler.LoginPOST(c)
	s.NoError(err)
	s.Equal(http.StatusSeeOther, rec.Code)
	s.Equal("/", rec.Header().Get("Location"))

	// Verify all mocks
	s.UserService.AssertExpectations(s.T())
	s.SessionManager.AssertExpectations(s.T())
	s.Logger.AssertExpectations(s.T())
	s.mockStore.AssertExpectations(s.T())
	s.mockSession.AssertExpectations(s.T())
	s.ErrorHandler.AssertExpectations(s.T())
}
