package shared

import (
	"net/http"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// FlashHandler handles flash messages and session operations
type FlashHandler struct {
	Store        sessions.Store
	SessionName  string
	Logger       loggo.LoggerInterface
	ErrorHandler *ErrorHandler
}

// FlashHandlerParams for dependency injection
type FlashHandlerParams struct {
	fx.In
	Store        sessions.Store
	Config       *config.Config // To get SessionName
	Logger       loggo.LoggerInterface
	ErrorHandler *ErrorHandler
}

// NewFlashHandler creates a new FlashHandler with dependency injection
func NewFlashHandler(params FlashHandlerParams) *FlashHandler {
	return &FlashHandler{
		Store:        params.Store,
		SessionName:  params.Config.SessionName,
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
