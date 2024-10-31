package server

import (
	"github.com/jonesrussell/loggo"
	"github.com/labstack/echo/v4"
)

// LoggingHandlerDecorator wraps a handler with logging
type LoggingHandlerDecorator struct {
	handler HandlerInterface
	logger  loggo.LoggerInterface
}

func NewLoggingHandlerDecorator(handler HandlerInterface, logger loggo.LoggerInterface) HandlerInterface {
	return &LoggingHandlerDecorator{
		handler: handler,
		logger:  logger,
	}
}

// Implement HandlerInterface methods
func (d *LoggingHandlerDecorator) HandleIndex(c echo.Context) error {
	d.logger.Info("Handling index request")
	defer d.logger.Info("Completed index request")
	return d.handler.HandleIndex(c)
}
