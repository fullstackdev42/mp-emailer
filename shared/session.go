package shared

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/mp-emailer/config"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// FlashHandlerInterface defines the methods that a flash handler must implement
type FlashHandlerInterface interface {
	SetFlashAndSaveSession(c echo.Context, message string) error
	ClearSession(c echo.Context) error
}

// FlashHandler handles flash messages and session operations
type FlashHandler struct {
	Store        sessions.Store
	SessionName  string
	Logger       loggo.LoggerInterface
	ErrorHandler ErrorHandlerInterface
}

// FlashHandlerParams for dependency injection
type FlashHandlerParams struct {
	fx.In
	Store        sessions.Store
	Config       *config.Config // To get SessionName
	Logger       loggo.LoggerInterface
	ErrorHandler ErrorHandlerInterface
}

// NewFlashHandler creates a new FlashHandler with dependency injection
func NewFlashHandler(params FlashHandlerParams) FlashHandlerInterface {
	return &FlashHandler{
		Store:        params.Store,
		SessionName:  params.Config.Auth.SessionName,
		Logger:       params.Logger,
		ErrorHandler: params.ErrorHandler,
	}
}

func (f *FlashHandler) SetFlashAndSaveSession(c echo.Context, message string) error {
	sess, err := f.Store.Get(c.Request(), f.SessionName)
	if err != nil {
		return f.ErrorHandler.HandleHTTPError(c, err, "Error getting session", http.StatusInternalServerError)
	}
	sess.AddFlash(message, "messages")
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		f.Logger.Error("Failed to save session", err)
		return f.ErrorHandler.HandleHTTPError(c, err, "Error saving session", http.StatusInternalServerError)
	}
	return nil
}

func (f *FlashHandler) ClearSession(c echo.Context) error {
	sess, err := f.Store.Get(c.Request(), f.SessionName)
	if err != nil {
		return f.ErrorHandler.HandleHTTPError(c, err, "Error getting session", http.StatusInternalServerError)
	}
	sess.Values = make(map[interface{}]interface{})
	sess.Options.MaxAge = -1
	return sess.Save(c.Request(), c.Response())
}

// Ensure FlashHandler implements FlashHandlerInterface
var _ FlashHandlerInterface = (*FlashHandler)(nil)
