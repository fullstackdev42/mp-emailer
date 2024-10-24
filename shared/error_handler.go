package shared

import (
	"net/http"

	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

type ErrorHandler struct {
	logger loggo.LoggerInterface
}

func NewErrorHandler(logger loggo.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

func (eh *ErrorHandler) HandleHTTPError(c echo.Context, err error, message string, statusCode int) error {
	eh.logger.Error("Unhandled error", err, "url", c.Request().URL.String())
	return c.Render(statusCode, "error.html", PageData{
		Title:   http.StatusText(statusCode),
		Content: map[string]string{"message": message},
	})
}
