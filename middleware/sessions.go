package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

// NewSessionsMiddleware creates a new session middleware
func NewSessionsMiddleware(store sessions.Store, logger loggo.LoggerInterface, sessionName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Debug("Session middleware processing request", "path", c.Request().URL.Path)

			session, err := store.Get(c.Request(), sessionName)
			if err != nil {
				logger.Error("Failed to get session", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Session error")
			}

			// Create a custom response writer to intercept the status code
			resWriter := c.Response().Writer
			c.Response().Writer = &responseWriter{
				ResponseWriter: resWriter,
				statusCode:     http.StatusOK,
			}

			// Call the next handler
			if err = next(c); err != nil {
				return err
			}

			// Save session after handler execution
			if err := store.Save(c.Request(), resWriter, session); err != nil {
				logger.Error("Failed to save session", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save session")
			}

			return nil
		}
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}
