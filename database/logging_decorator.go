package database

import (
	"database/sql"
	"time"

	"github.com/jonesrussell/loggo"
)

// Decorator interface defines methods that can be decorated
type Decorator interface {
	Interface
}

// LoggingDBDecorator adds logging functionality to database operations
type LoggingDBDecorator struct {
	db     Interface
	logger loggo.LoggerInterface
}

// NewLoggingDBDecorator creates a new logging decorator for database operations
func NewLoggingDBDecorator(db Interface, logger loggo.LoggerInterface) Decorator {
	return &LoggingDBDecorator{
		db:     db,
		logger: logger,
	}
}

// UserExists decorates the database UserExists method with logging
func (d *LoggingDBDecorator) UserExists(username, email string) (bool, error) {
	start := time.Now()
	d.logger.Debug("Checking if user exists",
		"username", username,
		"email", email,
	)

	exists, err := d.db.UserExists(username, email)

	d.logger.Debug("UserExists completed",
		"username", username,
		"email", email,
		"exists", exists,
		"duration", time.Since(start),
		"error", err,
	)

	return exists, err
}

// CreateUser decorates the database CreateUser method with logging
func (d *LoggingDBDecorator) CreateUser(id, username, email, passwordHash string) error {
	start := time.Now()
	d.logger.Debug("Creating new user",
		"id", id,
		"username", username,
		"email", email,
	)

	err := d.db.CreateUser(id, username, email, passwordHash)

	d.logger.Debug("CreateUser completed",
		"id", id,
		"username", username,
		"email", email,
		"duration", time.Since(start),
		"error", err,
	)

	return err
}

// LoginUser decorates the database LoginUser method with logging
func (d *LoggingDBDecorator) LoginUser(username, password string) (string, error) {
	start := time.Now()
	d.logger.Debug("Attempting user login",
		"username", username,
	)

	userID, err := d.db.LoginUser(username, password)

	d.logger.Debug("LoginUser completed",
		"username", username,
		"success", err == nil,
		"duration", time.Since(start),
		"error", err,
	)

	return userID, err
}

// Exec decorates the database Exec method with logging
func (d *LoggingDBDecorator) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	d.logger.Debug("Executing SQL query",
		"query", query,
		"args", args,
	)

	result, err := d.db.Exec(query, args...)

	d.logger.Debug("Exec completed",
		"duration", time.Since(start),
		"error", err,
	)

	return result, err
}

// Query decorates the database Query method with logging
func (d *LoggingDBDecorator) Query(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	d.logger.Debug("Executing SQL query",
		"query", query,
		"args", args,
	)

	rows, err := d.db.Query(query, args...)

	d.logger.Debug("Query completed",
		"duration", time.Since(start),
		"error", err,
	)

	return rows, err
}

// QueryRow decorates the database QueryRow method with logging
func (d *LoggingDBDecorator) QueryRow(query string, args ...interface{}) *sql.Row {
	start := time.Now()
	d.logger.Debug("Executing SQL query row",
		"query", query,
		"args", args,
	)

	row := d.db.QueryRow(query, args...)

	d.logger.Debug("QueryRow completed",
		"duration", time.Since(start),
	)

	return row
}
