package server

import (
	"github.com/labstack/echo/v4"
)

// Route represents a server route.
type Route struct {
	Method  string
	Pattern string
	Handler echo.HandlerFunc
}

// NewRoute creates a new Route.
func NewRoute(method, pattern string, handler echo.HandlerFunc) Route {
	return Route{method, pattern, handler}
}
