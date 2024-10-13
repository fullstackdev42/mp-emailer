package server

import (
	"github.com/fullstackdev42/mp-emailer/pkg/config"
	"github.com/fullstackdev42/mp-emailer/pkg/database"
	"github.com/fullstackdev42/mp-emailer/pkg/templates"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(config *config.Config, logger *loggo.Logger, db *database.DB, tmplManager *templates.TemplateManager) *echo.Echo {
	e := echo.New()
	e.Static("/static", "web/public")
	e.Renderer = echo.Renderer(tmplManager)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	store := sessions.NewCookieStore([]byte(config.SessionSecret))
	e.Use(session.Middleware(store))
	e.Use(user.SetAuthStatusMiddleware(store, logger))
	e.Use(dbMiddleware(db))

	return e
}

func dbMiddleware(db *database.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	}
}
