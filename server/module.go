package server

import (
	"log"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// ProvideModule provides the server module dependencies
func ProvideModule() fx.Option {
	return fx.Options(
		fx.Provide(
			NewHandler,
			NewTemplateManager,
			// Add other constructors here
		),
	)
}

// InvokeModule sets up the server routes
func InvokeModule(e *echo.Echo, handler *Handler) {
	// Register routes
	if success, err := handler.ProvideRoutes(e); err != nil {
		// Handle the error, e.g., log it or panic
		log.Fatalf("failed to provide routes: %v", err)
	} else if len(success) == 0 {
		log.Println("routes were not successfully provided")
	}
}
