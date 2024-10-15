package user

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) VerifyUser(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockService) RegisterUser(username, password, email string) error {
	args := m.Called(username, password, email)
	return args.Error(0)
}

type mockSessionStore struct {
	sessions map[string]*sessions.Session
}

func newMockSessionStore() *mockSessionStore {
	return &mockSessionStore{
		sessions: make(map[string]*sessions.Session),
	}
}

func (m *mockSessionStore) Get(_ *http.Request, name string) (*sessions.Session, error) {
	session, ok := m.sessions[name]
	if !ok {
		session = sessions.NewSession(m, name)
		m.sessions[name] = session
	}
	return session, nil
}

func (m *mockSessionStore) New(_ *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(m, name)
	m.sessions[name] = session
	return session, nil
}

func (m *mockSessionStore) Save(_ *http.Request, _ http.ResponseWriter, s *sessions.Session) error {
	m.sessions[s.Name()] = s
	return nil
}

func TestHandler_HandleLogin(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockService, *loggo.MockLogger)
		method         string
		username       string
		password       string
		wantStatusCode int
		wantBody       string
		wantRedirect   string
		checkSession   func(*testing.T, *http.Response, *mockSessionStore)
	}{
		{
			name: "Successful login",
			setupMock: func(ms *MockService, ml *loggo.MockLogger) {
				ms.On("VerifyUser", "validuser", "validpass").Return("123", nil)
				ml.On("Debug", mock.Anything, mock.Anything).Return().Twice()
			},
			method:         http.MethodPost,
			username:       "validuser",
			password:       "validpass",
			wantStatusCode: http.StatusSeeOther,
			wantRedirect:   "/campaigns",
			checkSession: func(t *testing.T, _ *http.Response, mss *mockSessionStore) {
				sess, ok := mss.sessions["mpe"]
				assert.True(t, ok, "Session should be created")
				assert.Equal(t, "123", sess.Values["userID"])
				assert.Equal(t, "validuser", sess.Values["username"])
				assert.Equal(t, "123", sess.Values["userID"])
				assert.Equal(t, "validuser", sess.Values["username"])
			},
		},
		{
			name: "Invalid credentials",
			setupMock: func(ms *MockService, ml *loggo.MockLogger) {
				ms.On("VerifyUser", "testuser", "wrongpassword").Return("", fmt.Errorf("invalid username or password"))
				ml.On("Debug", mock.Anything, mock.Anything).Return().Twice()
				ml.On("Warn", mock.Anything, mock.Anything).Return()
			},
			method:         http.MethodPost,
			username:       "testuser",
			password:       "wrongpassword",
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       "Invalid username or password",
		},
		{
			name: "Empty username",
			setupMock: func(_ *MockService, ml *loggo.MockLogger) {
				ml.On("Debug", mock.Anything, mock.Anything).Return().Twice()
				ml.On("Warn", mock.Anything, mock.Anything).Return()
			},
			method:         http.MethodPost,
			username:       "",
			password:       "somepassword",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "Username and password are required",
		},
		{
			name: "Empty password",
			setupMock: func(_ *MockService, ml *loggo.MockLogger) {
				ml.On("Debug", mock.Anything, mock.Anything).Return().Twice()
				ml.On("Warn", mock.Anything, mock.Anything).Return()
			},
			method:         http.MethodPost,
			username:       "someuser",
			password:       "",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "Username and password are required",
		},
		// Session expiration
		{
			name: "Session Expired",
			setupMock: func(ms *MockService, ml *loggo.MockLogger) {
				ms.On("VerifyUser", "validuser", "validpass").Return("123", nil)
				ml.On("Debug", mock.Anything, mock.Anything).Return().Twice()
			},
			method:         http.MethodPost,
			username:       "validuser",
			password:       "validpass",
			wantStatusCode: http.StatusSeeOther,
			wantRedirect:   "/campaigns",
			checkSession: func(t *testing.T, _ *http.Response, mss *mockSessionStore) {
				_, ok := mss.sessions["mpe"]
				assert.True(t, ok, "Session should be created")
				// Simulate session expiration
				delete(mss.sessions, "mpe")
				_, ok = mss.sessions["mpe"]
				assert.False(t, ok, "Session should be expired")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			mockLogger := new(loggo.MockLogger)
			mockSessionStore := newMockSessionStore()

			if tt.setupMock != nil {
				tt.setupMock(mockService, mockLogger)
			}

			handler := NewHandler(mockService, mockLogger)

			e := echo.New()
			req := httptest.NewRequest(tt.method, "/login", strings.NewReader(url.Values{
				"username": {tt.username},
				"password": {tt.password},
			}.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set up the session store
			c.Set("_session_store", mockSessionStore)

			err := handler.HandleLogin(c)

			t.Logf("Response status: %d", rec.Code)
			t.Logf("Response headers: %+v", rec.Header())
			t.Logf("Response body: %s", rec.Body.String())

			if tt.wantStatusCode == http.StatusSeeOther {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatusCode, rec.Code)
				assert.Equal(t, tt.wantRedirect, rec.Header().Get("Location"))
				if tt.checkSession != nil {
					tt.checkSession(t, rec.Result(), mockSessionStore)
				}
			} else {
				if assert.Error(t, err) {
					httpError, ok := err.(*echo.HTTPError)
					if assert.True(t, ok) {
						assert.Equal(t, tt.wantStatusCode, httpError.Code)
						assert.Equal(t, tt.wantBody, httpError.Message)
					}
				}
			}

			mockService.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}
