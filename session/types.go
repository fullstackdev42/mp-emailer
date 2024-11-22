package session

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

// Manager defines the session manager interface
type Manager interface {
	GetSession(c echo.Context, name string) (*sessions.Session, error)
	SaveSession(c echo.Context, sess *sessions.Session) error
	GetFlashes(sess *sessions.Session) []interface{}
	AddFlash(sess *sessions.Session, value interface{})
	SetSessionValues(sess *sessions.Session, user interface{})
	GetSessionValue(sess *sessions.Session, key string) (interface{}, error)
	ValidateSession(c echo.Context) error
	ClearSession(c echo.Context, name string) error
	IsAuthenticated(c echo.Context) bool
}

// Store extends the basic sessions.Store interface with additional security features
type Store interface {
	sessions.Store

	// Additional security methods
	RegenerateID(r *http.Request, w http.ResponseWriter) (string, error)
	SetSecure(secure bool)
	SetSameSite(mode http.SameSite)
	SetOptions(options *sessions.Options)
	Cleanup(threshold time.Time) error
}

// Options configures session behavior
type Options struct {
	MaxAge          int
	CleanupInterval time.Duration
	SecurityKey     []byte
	CookieName      string
	Domain          string
	Secure          bool
	HTTPOnly        bool
	SameSite        http.SameSite
	Path            string
	MaxLength       int
	KeyPrefix       string
}

// Data represents the standard session data structure
type Data struct {
	UserID          interface{}
	Username        string
	LastAccessed    time.Time
	CreatedAt       time.Time
	IsAuthenticated bool
	CustomData      map[string]interface{}
}

// UserData represents the minimal user data needed for sessions
type UserData interface {
	GetID() interface{}
	GetUsername() string
	GetCustomData() map[string]interface{}
}

// Common session-related errors
var (
	ErrSessionNotFound    = echo.NewHTTPError(http.StatusUnauthorized, "session not found")
	ErrSessionExpired     = echo.NewHTTPError(http.StatusUnauthorized, "session expired")
	ErrInvalidSession     = echo.NewHTTPError(http.StatusBadRequest, "invalid session")
	ErrSessionStoreFailed = echo.NewHTTPError(http.StatusInternalServerError, "session store failed")
	ErrInvalidKeySize     = echo.NewHTTPError(http.StatusInternalServerError, "invalid security key size: must be 16, 24, or 32 bytes")
)

// Interface abstracts the gorilla session for testing
type Interface interface {
	Get(key interface{}) interface{}
	Set(key interface{}, val interface{})
	Delete(key interface{})
	IsNew() bool
	Save(r *http.Request, w http.ResponseWriter) error
	AddFlash(value interface{}, vars ...string)
	Flashes(vars ...string) []interface{}
	Options() *sessions.Options
	Values() map[interface{}]interface{}
	GetValues() map[interface{}]interface{}
	GetID() string
}
