package session

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// StoreProvider defines the interface for providing and managing session stores
type StoreProvider interface {
	// GetStore returns the session store for the given request
	GetStore(r *http.Request) Store

	// SetStore sets the session store for the given context
	SetStore(c echo.Context, store Store)
}
