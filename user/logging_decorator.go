package user

import (
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/jonesrussell/loggo"
)

// LoggingServiceDecorator adds logging functionality to the ServiceInterface
type LoggingServiceDecorator struct {
	shared.LoggingServiceDecorator
	service ServiceInterface
}

// NewLoggingServiceDecorator creates a new instance of LoggingServiceDecorator
func NewLoggingServiceDecorator(service ServiceInterface, logger loggo.LoggerInterface) *LoggingServiceDecorator {
	return &LoggingServiceDecorator{
		LoggingServiceDecorator: *shared.NewLoggingServiceDecorator(service, logger),
		service:                 service,
	}
}

// Implement User-specific methods with logging
func (d *LoggingServiceDecorator) RegisterUser(dto *CreateDTO) (*User, error) {
	d.Logger.Info("Registering new user", "username", dto.Username)
	user, err := d.service.RegisterUser(dto)
	if err != nil {
		d.Logger.Error("Failed to register user", err)
	}
	return user, err
}

// GetUser logs the user retrieval and returns the user and error
func (d *LoggingServiceDecorator) GetUser(dto *GetDTO) (*User, error) {
	d.Logger.Info("Fetching user", "id", dto.Username)
	// Convert the returned DTO to User type
	userDTO, err := d.service.GetUser(dto)
	if err != nil {
		d.Logger.Error("Failed to fetch user", err)
		return nil, err
	}
	// Convert DTO to User or ensure service returns User directly
	user := &User{
		// Map DTO fields to User fields
		Username: userDTO.Username,
		// ... other field mappings
	}
	return user, nil
}

// Continue with other methods (e.g., UpdateUser, DeleteUser, etc.)...
