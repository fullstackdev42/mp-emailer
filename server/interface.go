package server

import "github.com/labstack/echo/v4"

// HandlerInterface defines the contract for server handlers
type HandlerInterface interface {
	HandleIndex(c echo.Context) error
}
