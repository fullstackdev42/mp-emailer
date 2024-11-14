package user_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fullstackdev42/mp-emailer/config"
	mocksShared "github.com/fullstackdev42/mp-emailer/mocks/shared"
	mocksUser "github.com/fullstackdev42/mp-emailer/mocks/user"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupTestHandler(t *testing.T) *user.Handler {
	mockService := mocksUser.NewMockServiceInterface(t)
	mockRepo := mocksUser.NewMockRepositoryInterface(t)
	mockErrorHandler := &shared.ErrorHandler{}
	mockTemplateRenderer := &mocksShared.MockTemplateRendererInterface{}
	store := sessions.NewCookieStore([]byte("test-secret"))
	flashHandler := shared.NewFlashHandler(shared.FlashHandlerParams{
		Store:        store,
		Config:       &config.Config{SessionName: "test_session"},
		Logger:       nil,
		ErrorHandler: mockErrorHandler,
	})

	return &user.Handler{
		Service:         mockService,
		Repo:            mockRepo,
		ErrorHandler:    mockErrorHandler,
		Store:           store,
		SessionName:     "test_session",
		Config:          &config.Config{},
		TemplateManager: mockTemplateRenderer,
		FlashHandler:    flashHandler,
	}
}

func TestLoginPOST(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name           string
		payload        string
		setupMocks     func(h *user.Handler)
		expectedStatus int
		expectedPath   string
	}{
		{
			name:    "Successful login",
			payload: `{"username": "testuser", "password": "password123"}`,
			setupMocks: func(h *user.Handler) {
				mockService := h.Service.(*mocksUser.MockServiceInterface)
				mockRepo := h.Repo.(*mocksUser.MockRepositoryInterface)

				// Create a real password hash for testing
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				testUser := &user.User{
					BaseModel: shared.BaseModel{
						ID: uuid.New(),
					},
					Username:     "testuser",
					PasswordHash: string(hashedPassword),
				}

				// Mock all Info calls in sequence
				mockService.On("Info", "Starting login attempt", "username", "").Return()
				mockService.On("Info", "User found, attempting password verification", "username", "testuser").Return()
				mockService.On("Info", "Password verified successfully", "username", "testuser").Return()
				mockService.On("Info", "Login successful", "username", "testuser").Return()

				mockRepo.On("FindByUsername", "testuser").Return(testUser, nil)
			},
			expectedStatus: http.StatusSeeOther,
			expectedPath:   "/",
		},
		{
			name:    "Invalid credentials",
			payload: `{"username": "wronguser", "password": "wrongpass"}`,
			setupMocks: func(h *user.Handler) {
				mockService := h.Service.(*mocksUser.MockServiceInterface)
				mockRepo := h.Repo.(*mocksUser.MockRepositoryInterface)

				// Mock all Info calls in sequence
				mockService.On("Info", "Starting login attempt", "username", "").Return()
				mockService.On("Info", "Login failed - user not found", "username", "wronguser").Return()

				mockRepo.On("FindByUsername", "wronguser").Return(nil, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedPath:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := setupTestHandler(t)
			tt.setupMocks(h)

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Create and set up the session
			sess, _ := h.Store.Get(req, h.SessionName)
			c.Set("session", sess)

			// Call the handler
			err := h.LoginPOST(c)

			if tt.expectedPath != "" {
				// For successful login
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, tt.expectedPath, rec.Header().Get("Location"))

				// Verify session values
				sess, _ = h.Store.Get(req, h.SessionName)
				assert.True(t, sess.Values["authenticated"].(bool))
				assert.Equal(t, "testuser", sess.Values["username"])
			} else {
				// For failed login
				assert.Error(t, err)
				if he, ok := err.(*echo.HTTPError); ok {
					assert.Equal(t, tt.expectedStatus, he.Code)
				}
			}
		})
	}
}
