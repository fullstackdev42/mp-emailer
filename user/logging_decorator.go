package user

import (
	"github.com/jonesrussell/loggo"
)

// LoggingServiceDecorator adds logging functionality to the User ServiceInterface
type LoggingServiceDecorator struct {
	service ServiceInterface
	logger  loggo.LoggerInterface
}

// NewLoggingServiceDecorator creates a new instance of LoggingServiceDecorator
func NewLoggingServiceDecorator(service ServiceInterface, logger loggo.LoggerInterface) *LoggingServiceDecorator {
	return &LoggingServiceDecorator{
		service: service,
		logger:  logger,
	}
}

// Implement User-specific methods with logging
func (d *LoggingServiceDecorator) RegisterUser(dto *RegisterDTO) (*DTO, error) {
	d.logger.Info("Registering new user", "username", dto.Username)
	user, err := d.service.RegisterUser(dto)
	if err != nil {
		d.logger.Error("Failed to register user", err)
	}
	return user, err
}

// LoginUser logs the login attempt and the result
func (d *LoggingServiceDecorator) LoginUser(dto *LoginDTO) (string, error) {
	d.logger.Info("Logging in user", "username", dto.Username)
	token, err := d.service.LoginUser(dto)
	if err != nil {
		d.logger.Error("Failed to login user", err)
	}
	return token, err
}

// GetUser fetches a user by their ID
func (d *LoggingServiceDecorator) GetUser(dto *GetDTO) (*DTO, error) {
	d.logger.Info("Fetching user", "id", dto.ID)
	user, err := d.service.GetUser(dto)
	if err != nil {
		d.logger.Error("Failed to fetch user", err)
	}
	return user, err
}
