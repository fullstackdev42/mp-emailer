package session_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/mp-emailer/mocks"
	"github.com/jonesrussell/mp-emailer/session"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStore implements sessions.Store for testing
type MockStore struct {
	mock.Mock
}

func (m *MockStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	args := m.Called(r, name)
	if sess, ok := args.Get(0).(*sessions.Session); ok {
		return sess, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	args := m.Called(r, w, s)
	return args.Error(0)
}

func (m *MockStore) New(r *http.Request, name string) (*sessions.Session, error) {
	args := m.Called(r, name)
	if sess, ok := args.Get(0).(*sessions.Session); ok {
		return sess, args.Error(1)
	}
	return sessions.NewSession(m, name), nil
}

func TestNewCleaner(t *testing.T) {
	// Arrange
	store := &MockStore{}
	logger := mocks.NewMockLoggerInterface(t)
	interval := 15 * time.Minute
	maxAge := 3600

	// Act
	cleaner := session.NewCleaner(store, interval, maxAge, logger)

	// Assert
	assert.NotNil(t, cleaner)
}

func TestCleanup(t *testing.T) {
	// Arrange
	store := &MockStore{}
	logger := mocks.NewMockLoggerInterface(t)
	interval := 15 * time.Millisecond
	maxAge := 3600

	cleaner := session.NewCleaner(store, interval, maxAge, logger)

	// Update expectation to match actual call with all parameters
	logger.EXPECT().Info(
		"Session cleanup check",
		"maxAge", maxAge,
		"threshold", mock.AnythingOfType("time.Time"),
		"cleanupInterval", interval,
	).Return()

	// Act
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	cleaner.StartCleanup(ctx)

	// Wait for at least one cleanup cycle
	time.Sleep(20 * time.Millisecond)

	// Assert
	mock.AssertExpectationsForObjects(t, logger)
}

func TestMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		setupSession   func(store *MockStore) *sessions.Session
		expectExpired  bool
		expectSaveErr  error
		expectedStatus int
		expectSave     bool
	}{
		{
			name: "Valid session",
			setupSession: func(store *MockStore) *sessions.Session {
				sess := sessions.NewSession(store, "test")
				sess.Values["last_accessed"] = time.Now()
				sess.Options = &sessions.Options{MaxAge: 3600}
				return sess
			},
			expectExpired:  false,
			expectSaveErr:  nil,
			expectedStatus: http.StatusOK,
			expectSave:     false,
		},
		{
			name: "Expired session",
			setupSession: func(store *MockStore) *sessions.Session {
				sess := sessions.NewSession(store, "test")
				sess.Values["last_accessed"] = time.Now().Add(-2 * time.Hour)
				sess.Options = &sessions.Options{MaxAge: 3600}
				return sess
			},
			expectExpired:  true,
			expectSaveErr:  nil,
			expectedStatus: http.StatusOK,
			expectSave:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			store := &MockStore{}
			logger := mocks.NewMockLoggerInterface(t)

			cleaner := session.NewCleaner(store, 15*time.Minute, 3600, logger)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			sess := tt.setupSession(store)

			// Set up mock expectations
			store.On("Get", mock.AnythingOfType("*http.Request"), "session").Return(sess, nil)

			if tt.expectSave {
				store.On("Save",
					mock.AnythingOfType("*http.Request"),
					mock.AnythingOfType("*httptest.ResponseRecorder"),
					mock.AnythingOfType("*sessions.Session"),
				).Run(func(args mock.Arguments) {
					s := args.Get(2).(*sessions.Session)
					assert.NotNil(t, s)
					assert.NotNil(t, s.Values["last_accessed"])
				}).Return(tt.expectSaveErr).Once()
			}

			if tt.expectExpired {
				logger.EXPECT().Error("Error deleting expired session", tt.expectSaveErr).Maybe()
			}

			next := echo.HandlerFunc(func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			})

			handler := cleaner.Middleware()(next)

			// Act
			err := handler(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// Wait a short time for any async operations
			time.Sleep(10 * time.Millisecond)

			store.AssertExpectations(t)
			logger.AssertExpectations(t)
		})
	}
}

func TestStartCleanup(t *testing.T) {
	// Arrange
	store := &MockStore{}
	logger := mocks.NewMockLoggerInterface(t)
	interval := 100 * time.Millisecond
	maxAge := 3600

	cleaner := session.NewCleaner(store, interval, maxAge, logger)

	// Update logger expectation to match actual call format
	logger.EXPECT().Info(
		"Session cleanup check",
		"maxAge", maxAge,
		"threshold", mock.AnythingOfType("time.Time"),
		"cleanupInterval", interval,
	).Return()

	// Act
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	cleaner.StartCleanup(ctx)
	time.Sleep(300 * time.Millisecond) // Wait for multiple cleanup cycles

	// Assert
	logger.AssertExpectations(t)
}
