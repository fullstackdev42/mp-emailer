package user

import (
	"github.com/jonesrussell/loggo"
)

// LoggingDecorator is a decorator for the ServiceInterface
type LoggingDecorator struct {
	service ServiceInterface
	logger  loggo.LoggerInterface
}

// NewLoggingDecorator creates a new instance of LoggingDecorator
func NewLoggingDecorator(service ServiceInterface, logger loggo.LoggerInterface) ServiceInterface {
	return &LoggingDecorator{
		service: service,
		logger:  logger,
	}
}

// Info logs an info message with the given parameters
func (d *LoggingDecorator) Info(message string, params ...interface{}) {
	d.logger.Info(message, params...)
}

// Warn logs a warning message with the given parameters
func (d *LoggingDecorator) Warn(message string, params ...interface{}) {
	d.logger.Warn(message, params...)
}

// Error logs an error message with the given parameters
func (d *LoggingDecorator) Error(message string, err error, params ...interface{}) {
	d.logger.Error(message, err, params...)
}

// GetUser gets a user
func (d *LoggingDecorator) GetUser(dto *GetDTO) (*DTO, error) {
	d.logger.Info("Getting user", "dto", dto)
	user, err := d.service.GetUser(dto)
	if err != nil {
		d.logger.Error("Failed to get user", err, "dto", dto)
	}
	return user, err
}

// LoginUser logs in a user
func (d *LoggingDecorator) LoginUser(dto *LoginDTO) (string, error) {
	d.logger.Info("Logging in user", "dto", dto)
	token, err := d.service.LoginUser(dto)
	if err != nil {
		d.logger.Error("Failed to login user", err, "dto", dto)
	}
	return token, err
}

// RegisterUser registers a new user
func (d *LoggingDecorator) RegisterUser(dto *RegisterDTO) (*DTO, error) {
	d.logger.Info("Registering new user", "dto", dto)
	user, err := d.service.RegisterUser(dto)
	if err != nil {
		d.logger.Error("Failed to register user", err, "dto", dto)
	}
	return user, err
}
