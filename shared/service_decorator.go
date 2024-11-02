package shared

import (
	"github.com/jonesrussell/loggo"
)

// LoggableService defines the basic logging operations that a service should support
type LoggableService interface {
	Info(message string, params ...interface{})
	Warn(message string, params ...interface{})
	Error(message string, err error, params ...interface{})
}

// ServiceWithLogging combines LoggableService with common service operations
type ServiceWithLogging interface {
	LoggableService
}

// GenericLoggingDecorator adds logging functionality to any LoggableService
type GenericLoggingDecorator[T ServiceWithLogging] struct {
	Service T
	Logger  loggo.LoggerInterface
}

// NewGenericLoggingDecorator creates a new instance of GenericLoggingDecorator
func NewGenericLoggingDecorator[T ServiceWithLogging](service T, logger loggo.LoggerInterface) *GenericLoggingDecorator[T] {
	return &GenericLoggingDecorator[T]{
		Service: service,
		Logger:  logger,
	}
}

// Base logging methods
func (d *GenericLoggingDecorator[T]) Info(message string, params ...interface{}) {
	d.Logger.Info(message, params...)
	d.Service.Info(message, params...)
}

func (d *GenericLoggingDecorator[T]) Warn(message string, params ...interface{}) {
	d.Logger.Warn(message, params...)
	d.Service.Warn(message, params...)
}

func (d *GenericLoggingDecorator[T]) Error(message string, err error, params ...interface{}) {
	d.Logger.Error(message, err, params...)
	d.Service.Error(message, err, params...)
}
