package user

import (
	"github.com/jonesrussell/loggo"
)

type LoggingDecorator struct {
	service ServiceInterface
	logger  loggo.LoggerInterface
}

func NewLoggingDecorator(service ServiceInterface, logger loggo.LoggerInterface) ServiceInterface {
	return &LoggingDecorator{
		service: service,
		logger:  logger,
	}
}

// Implement LoggableService methods
func (d *LoggingDecorator) Info(message string, params ...interface{}) {
	d.logger.Info(message, params...)
}

func (d *LoggingDecorator) Warn(message string, params ...interface{}) {
	d.logger.Warn(message, params...)
}

func (d *LoggingDecorator) Error(message string, err error, params ...interface{}) {
	d.logger.Error(message, err, params...)
}

// Implement all ServiceInterface methods with logging
func (d *LoggingDecorator) GetUser(dto *GetDTO) (*DTO, error) {
	d.logger.Info("Getting user", "dto", dto)
	user, err := d.service.GetUser(dto)
	if err != nil {
		d.logger.Error("Failed to get user", err, "dto", dto)
	}
	return user, err
}

// Update LoginUser method to match ServiceInterface
func (d *LoggingDecorator) LoginUser(dto *LoginDTO) (string, error) {
	d.logger.Info("Logging in user", "dto", dto)
	token, err := d.service.LoginUser(dto)
	if err != nil {
		d.logger.Error("Failed to login user", err, "dto", dto)
	}
	return token, err
}

// Add RegisterUser method to match ServiceInterface
func (d *LoggingDecorator) RegisterUser(dto *RegisterDTO) (*DTO, error) {
	d.logger.Info("Registering new user", "dto", dto)
	user, err := d.service.RegisterUser(dto)
	if err != nil {
		d.logger.Error("Failed to register user", err, "dto", dto)
	}
	return user, err
}
