package database

import (
	"time"

	"github.com/jonesrussell/loggo"
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

func (d *LoggingDBDecorator) Query(query string, args ...interface{}) Result {
	d.Logger.Debug("Querying", "query", query, "args", args)
	return d.DB.Query(query, args...)
}

// Add Association method implementation
func (d *LoggingDBDecorator) Association(column string) AssociationInterface {
	d.Logger.Debug("Getting association", "column", column)
	return d.DB.Association(column)
}

// Add this method to the LoggingDBDecorator struct
func (d *LoggingDBDecorator) AutoMigrate(dst ...interface{}) error {
	d.Logger.Info("Executing AutoMigrate")
	start := time.Now()

	// Call the underlying db's AutoMigrate method
	err := d.DB.AutoMigrate(dst...)

	d.Logger.Info("Completed AutoMigrate",
		"duration", time.Since(start),
		"error", err)

	return err
}

// Add this method to the LoggingDBDecorator struct
func (d *LoggingDBDecorator) Delete(model interface{}) error {
	d.Logger.Debug("Deleting model", "model", model)
	err := d.DB.Delete(model)
	if err != nil {
		d.Logger.Error("Error deleting model", err, "model", model)
	}
	return err
}

// Add this method to the LoggingDBDecorator struct
func (d *LoggingDBDecorator) Group(name string) Interface {
	d.Logger.Debug("Creating database group", "name", name)
	return d.DB.Group(name)
}

// Add this method to the LoggingDBDecorator struct
func (d *LoggingDBDecorator) Having(query string, args ...interface{}) Interface {
	d.Logger.Debug("Applying Having clause", "query", query, "args", args)
	return d.DB.Having(query, args...)
}

// Add this method to the LoggingDBDecorator struct
func (d *LoggingDBDecorator) Joins(query string, args ...interface{}) Interface {
	d.Logger.Debug("Applying Joins clause", "query", query, "args", args)
	return d.DB.Joins(query, args...)
}

// Add this method to the LoggingDBDecorator struct
func (d *LoggingDBDecorator) Limit(limit int) Interface {
	d.Logger.Debug("Applying Limit clause", "limit", limit)
	return d.DB.Limit(limit)
}

// Add this method to the LoggingDBDecorator struct
func (d *LoggingDBDecorator) Offset(offset int) Interface {
	d.Logger.Debug("Applying Offset clause", "offset", offset)
	return d.DB.Offset(offset)
}

// Add this method to the LoggingDBDecorator struct
func (d *LoggingDBDecorator) Migrator() Migrator {
	d.Logger.Debug("Getting database migrator")
	return d.DB.Migrator()
}

func (d *LoggingDBDecorator) Not(query interface{}, args ...interface{}) Interface {
	d.Logger.Debug("Executing Not query", "query", query, "args", args)
	return d.DB.Not(query, args...)
}

func (d *LoggingDBDecorator) Or(query interface{}, args ...interface{}) Interface {
	d.Logger.Debug("Executing Or query", "query", query, "args", args)
	return d.DB.Or(query, args...)
}

func (d *LoggingDBDecorator) Order(value interface{}) Interface {
	d.Logger.Debug("Applying Order clause", "value", value)
	return d.DB.Order(value)
}

func (d *LoggingDBDecorator) Preload(query string, args ...interface{}) Interface {
	d.Logger.Debug("Applying Preload", "query", query, "args", args)
	return d.DB.Preload(query, args...)
}

func (d *LoggingDBDecorator) Unscoped() Interface {
	d.Logger.Debug("Applying Unscoped")
	return d.DB.Unscoped()
}

func (d *LoggingDBDecorator) Where(query interface{}, args ...interface{}) Interface {
	d.Logger.Debug("Applying Where clause", "query", query, "args", args)
	return d.DB.Where(query, args...)
}

func (d *LoggingDBDecorator) WithTrashed() Interface {
	d.Logger.Debug("Applying WithTrashed")
	return d.DB.WithTrashed()
}
