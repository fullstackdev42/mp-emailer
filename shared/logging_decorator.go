package shared

import (
	"github.com/jonesrussell/loggo"
)

// ServiceInterface is a generic interface for services to be decorated
type ServiceInterface interface {
	Info(message string, params ...interface{})
	Warn(message string, params ...interface{})
	Error(message string, err error, params ...interface{})
}

// LoggingServiceDecorator adds logging functionality to any ServiceInterface
type LoggingServiceDecorator struct {
	Service ServiceInterface
	Logger  loggo.LoggerInterface
}

// NewLoggingServiceDecorator creates a new instance of LoggingServiceDecorator
func NewLoggingServiceDecorator(service ServiceInterface, logger loggo.LoggerInterface) *LoggingServiceDecorator {
	return &LoggingServiceDecorator{
		Service: service,
		Logger:  logger,
	}
}
