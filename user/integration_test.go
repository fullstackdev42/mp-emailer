package user_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/mp-emailer/config"
	"github.com/jonesrussell/mp-emailer/mocks"
	mocksSession "github.com/jonesrussell/mp-emailer/mocks/session"
	mocksShared "github.com/jonesrussell/mp-emailer/mocks/shared"
	mocksUser "github.com/jonesrussell/mp-emailer/mocks/user"
	"github.com/jonesrussell/mp-emailer/shared"
	"github.com/jonesrussell/mp-emailer/user"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	handler        *user.Handler
	echo           *echo.Echo
	sessionManager *mocksSession.MockManager
	userService    *mocksUser.MockServiceInterface
}

func (s *IntegrationTestSuite) SetupTest() {
	s.echo = echo.New()

	// Initialize mocks and store
	mockStore := sessions.NewCookieStore([]byte("test-key"))
	mockLogger := mocks.NewMockLoggerInterface(s.T())
	mockErrorHandler := mocksShared.NewMockErrorHandlerInterface(s.T())
	mockTemplateRenderer := mocksShared.NewMockTemplateRendererInterface(s.T())
	mockSessionManager := mocksSession.NewMockManager(s.T())
	mockUserService := mocksUser.NewMockServiceInterface(s.T())
	mockRepo := mocksUser.NewMockRepositoryInterface(s.T())
	mockConfig := &config.Config{}

	// Add logger expectations
	mockLogger.On("Debug", "Processing login request").Return()
	mockLogger.On("Debug", "Attempting user authentication", "username", mock.AnythingOfType("string")).Return()
	mockLogger.On("Debug", "User authenticated successfully",
		"username", mock.AnythingOfType("string"),
		"userID", mock.AnythingOfType("uuid.UUID")).Return()
	mockLogger.On("Debug", "Login process completed successfully",
		"username", mock.AnythingOfType("string")).Return()

	// Create test session
	mockSession := sessions.NewSession(mockStore, "test_session")
	mockSession.Values = make(map[interface{}]interface{})

	// Create handler with mocked dependencies
	s.handler = user.NewHandler(user.HandlerParams{
		BaseHandlerParams: shared.BaseHandlerParams{
			Logger:           mockLogger,
			ErrorHandler:     mockErrorHandler,
			Config:           mockConfig,
			TemplateRenderer: mockTemplateRenderer,
		},
		Service: mockUserService,
		Repo:    mockRepo,
	})

	s.sessionManager = mockSessionManager
	s.userService = mockUserService
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestRegistrationToLoginFlow() {
	// Test data
	username := "testuser"
	email := "test@example.com"
	password := "securepassword123"

	// Create test user
	testUser := &user.User{
		BaseModel: shared.BaseModel{
			ID: uuid.New(),
		},
		Username: username,
		Email:    email,
	}

	// Update mock expectation for registration to return *user.DTO instead of *user.User
	s.userService.On("RegisterUser", mock.Anything, mock.MatchedBy(func(dto interface{}) bool {
		if d, ok := dto.(*user.RegisterDTO); ok {
			return d.Username == username &&
				d.Email == email &&
				d.Password == password &&
				d.PasswordConfirm == password
		}
		return false
	})).Return(&user.DTO{
		ID:       testUser.ID,
		Username: testUser.Username,
		Email:    testUser.Email,
	}, nil)

	// Update mock expectation for authentication
	s.userService.On("AuthenticateUser", mock.Anything, username, password).Return(true, testUser, nil)

	// Session expectations remain the same
	session := sessions.NewSession(sessions.NewCookieStore([]byte("test-key")), "test_session")
	session.Values = make(map[interface{}]interface{})

	s.sessionManager.On("GetSession", mock.Anything, mock.Anything).Return(session, nil)
	s.sessionManager.On("SaveSession", mock.Anything, session).Return(nil)
	s.sessionManager.On("SetSessionValues", session, mock.AnythingOfType("*user.User")).Run(func(args mock.Arguments) {
		sess := args.Get(0).(*sessions.Session)
		user := args.Get(1).(*user.User)
		// Simulate setting session values
		sess.Values["user_id"] = user.ID
		sess.Values["username"] = user.Username
	}).Return()

	// Step 1: Register new user
	jsonData := strings.NewReader(`{"username":"` + username + `","email":"` + email + `","password":"` + password + `","password_confirm":"` + password + `"}`)

	req := httptest.NewRequest(http.MethodPost, "/user/register", jsonData)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)

	// Execute registration
	err := s.handler.RegisterPOST(c)
	s.NoError(err)
	s.Equal(http.StatusSeeOther, rec.Code)
	s.Equal("/user/login", rec.Header().Get("Location"))

	// Get registration session and verify flash message
	regSession, err := s.sessionManager.GetSession(c, "test_session")
	s.NoError(err)
	regFlash := regSession.Flashes() // This will clear the flashes
	s.Len(regFlash, 1)
	s.Equal("Registration successful! Please log in.", regFlash[0])

	// Step 2: Attempt login with new credentials
	loginForm := url.Values{}
	loginForm.Set("username", username)
	loginForm.Set("password", password)

	req = httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(loginForm.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	c = s.echo.NewContext(req, rec)

	// Create and set session using GetSession
	session, err = s.sessionManager.GetSession(c, "test_session")
	s.NoError(err)
	c.Set("session", session)

	// Execute login
	err = s.handler.LoginPOST(c)
	s.NoError(err)
	s.Equal(http.StatusSeeOther, rec.Code)
	s.Equal("/", rec.Header().Get("Location"))

	// Verify session contains user data
	s.NotNil(session.Values["user_id"])
	s.Equal(username, session.Values["username"])

	// Verify login flash message
	loginFlash := session.Flashes()
	s.Len(loginFlash, 1)
	s.Equal("Successfully logged in!", loginFlash[0])
}
