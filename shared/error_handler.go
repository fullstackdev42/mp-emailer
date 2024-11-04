package shared

import (
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

// ErrorHandlerInterface defines the methods that an error handler must implement
type ErrorHandlerInterface interface {
	HandleHTTPError(c echo.Context, err error, message string, statusCode int) error
}

// ErrorHandler is the concrete type for handling errors
type ErrorHandler struct{}

// NewErrorHandler creates a new instance of ErrorHandler
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// HandleHTTPError handles HTTP errors consistently
func (h *ErrorHandler) HandleHTTPError(c echo.Context, err error, message string, status int) error {
	data := Data{
		Title: "Error",
		Content: map[string]interface{}{
			"error":   err.Error(),
			"message": message,
		},
		StatusCode: status,
		PageName:   "error",
	}

	return c.Render(status, "error", data)
}

// LoggingErrorHandlerDecorator adds logging functionality to the ErrorHandlerInterface
type LoggingErrorHandlerDecorator struct {
	errorHandler ErrorHandlerInterface
	logger       loggo.LoggerInterface
}

// NewLoggingErrorHandlerDecorator creates a new instance of LoggingErrorHandlerDecorator
func NewLoggingErrorHandlerDecorator(errorHandler ErrorHandlerInterface, logger loggo.LoggerInterface) *LoggingErrorHandlerDecorator {
	return &LoggingErrorHandlerDecorator{
		errorHandler: errorHandler,
		logger:       logger,
	}
}

// HandleHTTPError logs the error and then handles it
func (d *LoggingErrorHandlerDecorator) HandleHTTPError(c echo.Context, err error, message string, statusCode int) error {
	d.logger.Error("Unhandled error", err, "url", c.Request().URL.String())
	return d.errorHandler.HandleHTTPError(c, err, message, statusCode)
}
