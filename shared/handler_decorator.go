package shared

import (
	"github.com/jonesrussell/mp-emailer/logger"
	"github.com/labstack/echo/v4"
)

// HandlerLoggable defines the base logging interface for handlers
type HandlerLoggable interface {
	Info(message string, params ...interface{})
	Warn(message string, params ...interface{})
	Error(message string, err error, params ...interface{})
}

// LoggingHandlerDecorator adds logging functionality to any handler
type LoggingHandlerDecorator[T HandlerLoggable] struct {
	Handler T
	Logger  logger.Interface
}

// NewLoggingHandlerDecorator creates a new instance of LoggingHandlerDecorator
func NewLoggingHandlerDecorator[T HandlerLoggable](handler T, log logger.Interface) *LoggingHandlerDecorator[T] {
	return &LoggingHandlerDecorator[T]{
		Handler: handler,
		Logger:  log,
	}
}

// Info implements logging
func (d *LoggingHandlerDecorator[T]) Info(message string, params ...interface{}) {
	d.Logger.Info(message, params...)
	d.Handler.Info(message, params...)
}

// Warn implements logging
func (d *LoggingHandlerDecorator[T]) Warn(message string, params ...interface{}) {
	d.Logger.Warn(message, params...)
	d.Handler.Warn(message, params...)
}

// Error implements logging
func (d *LoggingHandlerDecorator[T]) Error(message string, err error, params ...interface{}) {
	d.Logger.Error(message, err, params...)
	d.Handler.Error(message, err, params...)
}

// IndexGET forwards the handler method while adding logging
func (d *LoggingHandlerDecorator[T]) IndexGET(c echo.Context) error {
	d.Logger.Info("Handling index request", "path", c.Path())
	// Type assertion to access the IndexGET method
	if handler, ok := interface{}(d.Handler).(interface{ IndexGET(echo.Context) error }); ok {
		return handler.IndexGET(c)
	}
	return echo.ErrMethodNotAllowed
}

func (d *LoggingHandlerDecorator[T]) HealthCheck(c echo.Context) error {
	d.Logger.Info("Health check requested")
	// Type assertion using interface{} conversion first, like IndexGET
	if handler, ok := interface{}(d.Handler).(interface{ HealthCheck(echo.Context) error }); ok {
		return handler.HealthCheck(c)
	}
	return echo.ErrMethodNotAllowed
}
