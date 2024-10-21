package shared

import (
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

type ErrorHandler struct {
	Logger loggo.LoggerInterface
}

func NewErrorHandler(logger loggo.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		Logger: logger,
	}
}

func (eh *ErrorHandler) HandleError(c echo.Context, err error, statusCode int, message string) error {
	eh.Logger.Error("Error occurred", err, "message", message, "statusCode", statusCode)

	pageData := PageData{
		Title: "Error",
		Error: message,
		Content: map[string]string{
			"Details": err.Error(),
		},
	}

	return c.Render(statusCode, "error.html", pageData)
}
