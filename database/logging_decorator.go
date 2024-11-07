package database

import (
	"github.com/jonesrussell/loggo"
	"gorm.io/gorm"
)

// LoggingDBDecorator is a decorator for logging database operations
type LoggingDBDecorator struct {
	DB     Interface
	Logger loggo.LoggerInterface
}

// NewLoggingDBDecorator creates a new LoggingDBDecorator
func NewLoggingDBDecorator(db Interface, logger loggo.LoggerInterface) *LoggingDBDecorator {
	return &LoggingDBDecorator{
		DB:     db,
		Logger: logger,
	}
}

func (d *LoggingDBDecorator) Exists(model interface{}, query string, args ...interface{}) (bool, error) {
	d.Logger.Debug("Checking existence", "model", model, "query", query, "args", args)
	exists, err := d.DB.Exists(model, query, args...)
	if err != nil {
		d.Logger.Error("Error checking existence", err, "model", model, "query", query, "args", args)
	}
	return exists, err
}

func (d *LoggingDBDecorator) Create(model interface{}) error {
	d.Logger.Debug("Creating model", "model", model)
	err := d.DB.Create(model)
	if err != nil {
		d.Logger.Error("Error creating model", err, "model", model)
	}
	return err
}

func (d *LoggingDBDecorator) FindOne(model interface{}, query string, args ...interface{}) error {
	d.Logger.Debug("Finding one", "model", model, "query", query, "args", args)
	err := d.DB.FindOne(model, query, args...)
	if err != nil {
		d.Logger.Error("Error finding one", err, "model", model, "query", query, "args", args)
	}
	return err
}

func (d *LoggingDBDecorator) Exec(query string, args ...interface{}) error {
	d.Logger.Debug("Executing query", "query", query, "args", args)
	err := d.DB.Exec(query, args...)
	if err != nil {
		d.Logger.Error("Error executing query", err, "query", query, "args", args)
	}
	return err
}

func (d *LoggingDBDecorator) Query(query string, args ...interface{}) *gorm.DB {
	d.Logger.Debug("Querying", "query", query, "args", args)
	return d.DB.Query(query, args...)
}
