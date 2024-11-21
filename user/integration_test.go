package user_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/mp-emailer/config"
	"github.com/jonesrussell/mp-emailer/mocks"
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
	sessionManager user.SessionManager
}

func (s *IntegrationTestSuite) SetupTest() {
	s.echo = echo.New()

	// Initialize mocks and store
	mockStore := sessions.NewCookieStore([]byte("test-key"))
	mockLogger := mocks.NewMockLoggerInterface(s.T())
	mockErrorHandler := mocksShared.NewMockErrorHandlerInterface(s.T())
	mockFlashHandler := mocksShared.NewMockFlashHandlerInterface(s.T())
	mockTemplateRenderer := mocksShared.NewMockTemplateRendererInterface(s.T())
	mockSessionManager := mocksUser.NewMockSessionManager(s.T())
	mockUserService := mocksUser.NewMockServiceInterface(s.T())
	mockRepo := mocksUser.NewMockRepositoryInterface(s.T())
	mockConfig := &config.Config{}

	// Create test session
	mockSession := sessions.NewSession(mockStore, "test_session")
	mockSession.Values = make(map[interface{}]interface{})

	// Define test user with complete data
	testUser := &user.User{
		BaseModel: shared.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$somehashedpassword",
	}

	// Setup logger expectations with correct argument patterns
	mockLogger.On("Debug", "Processing login request").Return()
	mockLogger.On("Debug", "Attempting user authentication", "username", testUser.Username).Return()
	mockLogger.On("Debug", "User authenticated successfully", "username", testUser.Username, "userID", testUser.ID).Return()
	mockLogger.On("Debug", "Login process completed successfully", "username", testUser.Username).Return()

	// Setup session expectations with proper user data
	mockSessionManager.On("GetSession", mock.Anything).Return(mockSession, nil)
	mockSessionManager.On("SetSessionValues", mockSession, testUser).Run(func(args mock.Arguments) {
		session := args.Get(0).(*sessions.Session)
		user := args.Get(1).(*user.User)
		// Set the actual session values
		session.Values["user_id"] = user.ID.String()
		session.Values["username"] = user.Username
	}).Return()
	mockSessionManager.On("SaveSession", mock.Anything, mockSession).Return(nil)

	// Setup authentication expectations with complete user data
	mockUserService.On("AuthenticateUser",
		mock.Anything,
		testUser.Username,
		"securepassword123",
	).Return(true, testUser, nil)

	// Setup expectations for registration
	mockUserService.On("RegisterUser",
		mock.Anything,
		mock.MatchedBy(func(params *user.RegisterDTO) bool {
			return params.Username == testUser.Username &&
				params.Email == testUser.Email &&
				params.Password == "securepassword123" &&
				params.PasswordConfirm == "securepassword123"
		}),
	).Return(&user.DTO{
		ID:        testUser.ID,
		Username:  testUser.Username,
		Email:     testUser.Email,
		CreatedAt: testUser.CreatedAt,
		UpdatedAt: testUser.UpdatedAt,
	}, nil)

	// Setup flash handler expectations
	mockFlashHandler.On("SetFlashAndSaveSession",
		mock.Anything,
		"Registration successful! Please log in.",
	).Run(func(_ mock.Arguments) {
		mockSession.AddFlash("Registration successful! Please log in.")
	}).Return(nil)

	mockFlashHandler.On("SetFlashAndSaveSession",
		mock.Anything,
		"Successfully logged in!",
	).Run(func(_ mock.Arguments) {
		mockSession.AddFlash("Successfully logged in!")
	}).Return(nil)

	// Create handler with mocked dependencies
	s.handler = user.NewHandler(user.HandlerParams{
		BaseHandlerParams: shared.BaseHandlerParams{
			Store:            mockStore,
			Logger:           mockLogger,
			ErrorHandler:     mockErrorHandler,
			Config:           mockConfig,
			TemplateRenderer: mockTemplateRenderer,
		},
		Service:        mockUserService,
		FlashHandler:   mockFlashHandler,
		Repo:           mockRepo,
		SessionManager: mockSessionManager,
	})

	s.sessionManager = mockSessionManager
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestRegistrationToLoginFlow() {
	// Test data
	username := "testuser"
	email := "test@example.com"
	password := "securepassword123"

	// Step 1: Register new user
	registrationForm := url.Values{}
	registrationForm.Set("username", username)
	registrationForm.Set("email", email)
	registrationForm.Set("password", password)
	registrationForm.Set("password_confirm", password)

	req := httptest.NewRequest(http.MethodPost, "/user/register", strings.NewReader(registrationForm.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)

	// Execute registration
	err := s.handler.RegisterPOST(c)
	s.NoError(err)
	s.Equal(http.StatusSeeOther, rec.Code)
	s.Equal("/user/login", rec.Header().Get("Location"))

	// Get registration session and verify flash message
	regSession, err := s.sessionManager.GetSession(c)
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
	session, err := s.sessionManager.GetSession(c)
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
	loginFlash := session.Flashes() // This will clear the flashes
	s.Len(loginFlash, 1)
	s.Equal("Successfully logged in!", loginFlash[0])
}