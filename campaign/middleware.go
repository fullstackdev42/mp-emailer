package campaign

import (
	"net/http"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

// Custom middleware for protected routes
func AuthMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session := getSession(c, cfg)
			if session == nil || session.Values["user_id"] == nil {
				// Store the original requested URL in the session
				if session != nil {
					session.Values["redirect_after_login"] = c.Request().URL.String()
					_ = session.Save(c.Request(), c.Response().Writer)
				}
				return c.Redirect(http.StatusSeeOther, "/user/login")
			}
			return next(c)
		}
	}
}

func getSession(c echo.Context, cfg *config.Config) *sessions.Session {
	store := c.Get("store").(sessions.Store)
	session, _ := store.Get(c.Request(), cfg.SessionName)
	return session
}
