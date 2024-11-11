package server

import (
	"github.com/labstack/echo/v4"
)

// Route represents a server route.
type Route struct {
	Method     string
	Pattern    string
	Handler    echo.HandlerFunc
	Middleware []echo.MiddlewareFunc
}

// NewRoute creates a new Route.
func NewRoute(method, pattern string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) Route {
	return Route{
		Method:     method,
		Pattern:    pattern,
		Handler:    handler,
		Middleware: middleware,
	}
}

// ApplyRoutes applies a slice of Routes to an Echo instance.
func ApplyRoutes(e *echo.Echo, routes []Route) {
	for _, route := range routes {
		e.Add(route.Method, route.Pattern, route.Handler, route.Middleware...)
	}
}

// RegisterRoutes registers all server routes
func RegisterRoutes(handler HandlerInterface, e *echo.Echo) {
	// Create routes
	routes := []Route{
		NewRoute("GET", "/", handler.IndexGET),
		NewRoute("GET", "/health", handler.HealthCheck),
	}

	// Apply routes
	ApplyRoutes(e, routes)
}
