package server

import (
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

func LoggingMiddleware(logger loggo.LoggerInterface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			logger.Info("Request",
				"method", req.Method,
				"uri", req.RequestURI,
				"remote_addr", c.RealIP(),
			)

			err := next(c)

			logger.Info("Response",
				"status", res.Status,
				"size", res.Size,
			)

			return err
		}
	}
}
