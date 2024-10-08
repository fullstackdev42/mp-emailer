package middleware

import (
	"github.com/labstack/echo/v4"
)

func IsAuthenticatedMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		isAuthenticated := checkAuthentication(c)
		c.Set("isAuthenticated", isAuthenticated)
		return next(c)
	}
}

func checkAuthentication(_ echo.Context) bool {
	// Your authentication check logic
	return true // Example: assuming all users are authenticated
}
