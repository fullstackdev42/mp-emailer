package user

import (
	"context"

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
func (d *LoggingDecorator) GetUser(ctx context.Context, dto *GetDTO) (*DTO, error) {
	d.logger.Info("Getting user", "dto", dto)
	user, err := d.service.GetUser(ctx, dto)
	if err != nil {
		d.logger.Error("Failed to get user", err, "dto", dto)
	}
	return user, err
}

// LoginUser logs in a user
func (d *LoggingDecorator) LoginUser(ctx context.Context, dto *LoginDTO) (string, error) {
	d.logger.Info("Logging in user", "dto", dto)
	token, err := d.service.LoginUser(ctx, dto)
	if err != nil {
		d.logger.Error("Failed to login user", err, "dto", dto)
	}
	return token, err
}

// RegisterUser registers a new user
func (d *LoggingDecorator) RegisterUser(ctx context.Context, dto *RegisterDTO) (*DTO, error) {
	d.logger.Info("Registering new user", "dto", dto)
	user, err := d.service.RegisterUser(ctx, dto)
	if err != nil {
		d.logger.Error("Failed to register user", err, "dto", dto)
	}
	return user, err
}

// AuthenticateUser authenticates a user with logging
func (d *LoggingDecorator) AuthenticateUser(ctx context.Context, username, password string) (*User, error) {
	d.logger.Info("Authenticating user", "username", username)
	user, err := d.service.AuthenticateUser(ctx, username, password)
	if err != nil {
		d.logger.Error("Failed to authenticate user", err, "username", username)
	}
	return user, err
}

// RequestPasswordReset decorates the password reset request with logging
func (d *LoggingDecorator) RequestPasswordReset(ctx context.Context, dto *PasswordResetDTO) error {
	d.logger.Info("Requesting password reset", "email", dto.Email)
	err := d.service.RequestPasswordReset(ctx, dto)
	if err != nil {
		d.logger.Error("Failed to request password reset", err, "email", dto.Email)
	}
	return err
}

// ResetPassword decorates the password reset completion with logging
func (d *LoggingDecorator) ResetPassword(ctx context.Context, dto *ResetPasswordDTO) error {
	d.logger.Info("Resetting password", "token", dto.Token)
	err := d.service.ResetPassword(ctx, dto)
	if err != nil {
		d.logger.Error("Failed to reset password", err, "token", dto.Token)
	}
	return err
}
