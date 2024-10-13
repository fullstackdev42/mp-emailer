package server

import (
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/user"
	"github.com/gorilla/sessions"
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(config *config.Config, logger *loggo.Logger, tmplManager *TemplateManager) *echo.Echo {
	e := echo.New()
	e.Static("/static", "web/public")
	e.Renderer = echo.Renderer(tmplManager)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	store := sessions.NewCookieStore([]byte(config.SessionSecret))
	e.Use(session.Middleware(store))
	e.Use(user.SetAuthStatusMiddleware(store, logger))

	return e
}
