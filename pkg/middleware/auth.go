package middleware

import (
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(store sessions.Store, logger loggo.LoggerInterface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isAuthenticated := checkAuthentication(c, store, logger)
			c.Set("isAuthenticated", isAuthenticated)

			if !isAuthenticated {
				return echo.ErrUnauthorized
			}

			return next(c)
		}
	}
}

func checkAuthentication(c echo.Context, store sessions.Store, logger loggo.LoggerInterface) bool {
	session, err := store.Get(c.Request(), "mpe")
	if err != nil {
		logger.Error("Error getting session", err)
		return false
	}

	username, ok := session.Values["username"]
	if !ok || username == "" {
		logger.Info("No username found in session")
		return false
	}

	return true
}
