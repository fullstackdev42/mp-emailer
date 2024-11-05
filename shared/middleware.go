package shared

import (
	"github.com/labstack/echo/v4"
)

// MethodOverride middleware converts POST requests with _method parameter to the specified method
func MethodOverride() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == "POST" {
				method := c.FormValue("_method")
				if method != "" {
					c.Request().Method = method
				}
			}
			return next(c)
		}
	}
}
