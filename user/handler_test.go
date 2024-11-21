package user_test

import (
	"errors"
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
	mocksUser "github.com/jonesrussell/mp-emailer/mocks/user"
)

type HandlerTestSuite struct {
	testutil.BaseTestSuite
	handler        *user.Handler
	UserService    *mocksUser.MockServiceInterface
	UserRepo       *mocksUser.MockRepositoryInterface
	SessionManager *mocksSession.MockManager
}

func (s *HandlerTestSuite) SetupTest() {
	s.BaseTestSuite.SetupTest()

	s.UserService = mocksUser.NewMockServiceInterface(s.T())
	s.UserRepo = mocksUser.NewMockRepositoryInterface(s.T())
	s.SessionManager = mocksSession.NewMockManager(s.T())

	s.Config.Auth.SessionName = "test_session"

	flashHandler := shared.NewFlashHandler(shared.FlashHandlerParams{
		Store:        s.Store,
		Config:       s.Config,
		Logger:       s.Logger,
		ErrorHandler: s.ErrorHandler,
	})

	s.Echo.Renderer = s.TemplateRenderer

	params := user.HandlerParams{
		BaseHandlerParams: shared.BaseHandlerParams{
			Logger:           s.Logger,
			ErrorHandler:     s.ErrorHandler,
			TemplateRenderer: s.TemplateRenderer,
			Store:            s.Store,
			Config:           s.Config,
		},
		Service:        s.UserService,
		FlashHandler:   flashHandler,
		Repo:           s.UserRepo,
		SessionManager: s.SessionManager,
	}

	s.handler = user.NewHandler(params)
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) debugResponse(rec *httptest.ResponseRecorder) {
	s.T().Logf("Response Status: %d", rec.Code)
	s.T().Logf("Response Headers: %v", rec.Header())
	s.T().Logf("Response Body: %s", rec.Body.String())
}

func (s *HandlerTestSuite) TestLoginPOST() {
	tests := []struct {
		name    string
		payload string

		setupMocks     func() *sessions.Session
		expectedStatus int
		expectedPath   string
	}{
		{
			name:    "Successful login",
			payload: `{"username": "testuser", "password": "password123"}`,
			setupMocks: func() *sessions.Session {
				sess := sessions.NewSession(s.Store, "test_session")
				sess.Values = make(map[interface{}]interface{})

				testUser := &user.User{
					BaseModel: shared.BaseModel{
						ID: uuid.New(),
					},
					Username: "testuser",
				}

				s.Logger.On("Debug", "Processing login request").Return()
				s.Logger.On("Debug", "Attempting user authentication", "username", testUser.Username).Return()
				s.Logger.On("Debug", "User authenticated successfully", "username", testUser.Username, "userID", testUser.ID).Return()
				s.Logger.On("Debug", "Login process completed successfully", "username", testUser.Username).Return()

				s.SessionManager.On("GetSession", mock.Anything, "test_session").Return(sess, nil)
				s.SessionManager.On("SetSessionValues", sess, testUser).Return()
				s.SessionManager.On("SaveSession", mock.Anything, sess).Return(nil)

				s.UserService.On("AuthenticateUser",
					mock.Anything,
					testUser.Username,
					"password123",
				).Return(true, testUser, nil)

				return sess
			},
			expectedStatus: http.StatusSeeOther,
			expectedPath:   "/",
		},
		{
			name:    "Session store error",
			payload: `{"username": "testuser", "password": "password123"}`,
			setupMocks: func() *sessions.Session {
				sessionErr := errors.New("session store error")

				s.Logger.On("Debug", "Processing login request").Return()
				s.Logger.On("Error", "Failed to get session", sessionErr).Return()

				s.SessionManager.On("GetSession", mock.Anything, "test_session").
					Return(nil, sessionErr)

				s.ErrorHandler.On("HandleHTTPError",
					mock.AnythingOfType("*echo.context"),
					sessionErr,
					"Error getting session",
					http.StatusInternalServerError,
				).Return(echo.NewHTTPError(http.StatusInternalServerError))

				return nil
			},
			expectedStatus: http.StatusInternalServerError,
			expectedPath:   "",
		},
		{
			name:    "Invalid credentials",
			payload: `{"username": "wronguser", "password": "wrongpass"}`,
			setupMocks: func() *sessions.Session {
				sess := sessions.NewSession(s.Store, "test_session")
				sess.Values = make(map[interface{}]interface{})

				s.Logger.On("Debug", "Processing login request").Return()
				s.Logger.On("Debug", "Attempting user authentication", "username", "wronguser").Return()
				s.Logger.On("Debug", "Invalid login attempt", "username", "wronguser").Return()

				s.SessionManager.On("GetSession", mock.Anything, "test_session").Return(sess, nil)

				s.TemplateRenderer.On("Render",
					mock.AnythingOfType("*bytes.Buffer"),
					"login",
					mock.AnythingOfType("*shared.Data"),
					mock.AnythingOfType("*echo.context"),
				).Return(nil)

				s.UserService.On("AuthenticateUser",
					mock.Anything,
					"wronguser",
					"wrongpass",
				).Return(false, nil, nil)

				return sess
			},
			expectedStatus: http.StatusUnauthorized,
			expectedPath:   "",
		},
		{
			name:    "Session save error",
			payload: `{"username": "testuser", "password": "password123"}`,
			setupMocks: func() *sessions.Session {
				sess := sessions.NewSession(s.Store, "test_session")
				sess.Values = make(map[interface{}]interface{})

				testUser := &user.User{
					BaseModel: shared.BaseModel{
						ID: uuid.New(),
					},
					Username: "testuser",
				}

				saveErr := errors.New("failed to save session")

				s.Logger.On("Debug", "Processing login request").Return()
				s.Logger.On("Debug", "Attempting user authentication", "username", testUser.Username).Return()
				s.Logger.On("Debug", "User authenticated successfully", "username", testUser.Username, "userID", testUser.ID).Return()
				s.Logger.On("Error", "Failed to save session", saveErr).Return()

				s.SessionManager.On("GetSession", mock.Anything, "test_session").Return(sess, nil)
				s.SessionManager.On("SetSessionValues", sess, testUser).Return()
				s.SessionManager.On("SaveSession", mock.Anything, sess).Return(saveErr)

				s.UserService.On("AuthenticateUser",
					mock.Anything,
					testUser.Username,
					"password123",
				).Return(true, testUser, nil)

				s.ErrorHandler.On("HandleHTTPError",
					mock.AnythingOfType("*echo.context"),
					saveErr,
					"Error saving session",
					http.StatusInternalServerError,
				).Return(echo.NewHTTPError(http.StatusInternalServerError))

				return sess
			},
			expectedStatus: http.StatusInternalServerError,
			expectedPath:   "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest() // Reset mocks for each test case
			sess := tt.setupMocks()

			// Create request with JSON payload
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			// Create Echo context
			c := s.Echo.NewContext(req, rec)
			if sess != nil {
				c.Set("session", sess)
			}

			// Execute handler
			err := s.handler.LoginPOST(c)
			if err != nil {
				s.T().Logf("Handler returned error: %v", err)
				if he, ok := err.(*echo.HTTPError); ok {
					s.T().Logf("HTTP Error: code=%d, message=%v", he.Code, he.Message)
					rec.Code = he.Code // Set the status code from the error
				}
			}

			// Debug response
			s.debugResponse(rec)

			// Verify response
			s.Equal(tt.expectedStatus, rec.Code)
			if tt.expectedPath != "" {
				s.Equal(tt.expectedPath, rec.Header().Get("Location"))
			}

			// Verify all mocks
			s.SessionManager.AssertExpectations(s.T())
		})
	}
}
