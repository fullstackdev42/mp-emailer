package user_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fullstackdev42/mp-emailer/internal/testutil"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	mocksMiddleware "github.com/fullstackdev42/mp-emailer/mocks/middleware"
	mocksUser "github.com/fullstackdev42/mp-emailer/mocks/user"
)

type HandlerTestSuite struct {
	testutil.BaseTestSuite
	handler     *user.Handler
	UserService *mocksUser.MockServiceInterface
	UserRepo    *mocksUser.MockRepositoryInterface
}

func (s *HandlerTestSuite) SetupTest() {
	s.BaseTestSuite.SetupTest()

	// Initialize user-specific mocks
	s.UserService = mocksUser.NewMockServiceInterface(s.T())
	s.UserRepo = mocksUser.NewMockRepositoryInterface(s.T())

	// Ensure valid session name in config
	s.Config.SessionName = "test_session" // Valid cookie name without spaces

	// Create FlashHandler
	flashHandler := shared.NewFlashHandler(shared.FlashHandlerParams{
		Store:        s.Store,
		Config:       s.Config,
		Logger:       s.Logger,
		ErrorHandler: s.ErrorHandler,
	})

	// Register renderer with Echo
	s.Echo.Renderer = s.TemplateRenderer

	params := user.HandlerParams{
		BaseHandlerParams: shared.BaseHandlerParams{
			Logger:           s.Logger,
			ErrorHandler:     s.ErrorHandler,
			TemplateRenderer: s.TemplateRenderer,
			Store:            s.Store,
			Config:           s.Config,
		},
		Service:      s.UserService,
		FlashHandler: flashHandler,
		Repo:         s.UserRepo,
	}

	result, err := user.NewHandler(params)
	s.NoError(err)
	s.handler = result.Handler
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) TestLoginPOST() {
	tests := []struct {
		name           string
		payload        string
		setupMocks     func() *sessions.Session
		expectedStatus int
		expectedPath   string
	}{
		{
			name:    "Successful login",
			payload: `{"username": "testuser", "password": "password123"}`,
			setupMocks: func() *sessions.Session {
				testUser := &user.User{
					BaseModel: shared.BaseModel{
						ID: uuid.New(),
					},
					Username: "testuser",
				}

				// Create test session
				mockStore := mocksMiddleware.NewMockSessionStore(s.T())
				s.Store = mockStore
				sess := sessions.NewSession(mockStore, s.Config.SessionName)
				sess.Values = make(map[interface{}]interface{})

				// Mock session store with more specific matchers
				mockStore.On("Get",
					mock.MatchedBy(func(*http.Request) bool { return true }),
					s.Config.SessionName,
				).Return(sess, nil)

				// Mock session store Save method
				mockStore.On("Save",
					mock.MatchedBy(func(*http.Request) bool { return true }),
					mock.MatchedBy(func(http.ResponseWriter) bool { return true }),
					mock.MatchedBy(func(*sessions.Session) bool { return true }),
				).Return(nil).Times(2)

				// Mock authentication
				s.UserService.EXPECT().AuthenticateUser("testuser", "password123").
					Return(true, testUser, nil)

				return sess
			},
			expectedStatus: http.StatusSeeOther,
			expectedPath:   "/",
		},
		{
			name:    "Invalid credentials",
			payload: `{"username": "wronguser", "password": "wrongpass"}`,
			setupMocks: func() *sessions.Session {
				mockStore := mocksMiddleware.NewMockSessionStore(s.T())
				s.Store = mockStore
				sess := sessions.NewSession(mockStore, s.Config.SessionName)
				sess.Values = make(map[interface{}]interface{})

				mockStore.On("Get",
					mock.MatchedBy(func(*http.Request) bool { return true }),
					s.Config.SessionName,
				).Return(sess, nil)

				// Mock template renderer
				s.TemplateRenderer.On("Render",
					mock.AnythingOfType("*bytes.Buffer"),
					"login",
					mock.AnythingOfType("*shared.Data"),
					mock.AnythingOfType("*echo.context"),
				).Return(nil)

				// Mock failed authentication
				s.UserService.EXPECT().AuthenticateUser("wronguser", "wrongpass").
					Return(false, nil, nil)

				return sess
			},
			expectedStatus: http.StatusUnauthorized,
			expectedPath:   "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest() // Reset mocks for each test case
			sess := tt.setupMocks()

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := s.Echo.NewContext(req, rec)
			c.Set("session", sess)

			err := s.handler.LoginPOST(c)
			s.NoError(err)

			s.Equal(tt.expectedStatus, rec.Code)
			if tt.expectedPath != "" {
				s.Equal(tt.expectedPath, rec.Header().Get("Location"))
			}
		})
	}
}
