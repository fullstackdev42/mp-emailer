package decorators

import (
	"database/sql"
	"time"

	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/mp-emailer/database/core"
	"gorm.io/gorm"
)

// LoggingDecorator is a decorator for logging database operations
type LoggingDecorator struct {
	Database core.Interface
	Logger   loggo.LoggerInterface
}

// NewLoggingDecorator creates a new LoggingDecorator
func NewLoggingDecorator(db core.Interface, logger loggo.LoggerInterface) core.Interface {
	return &LoggingDecorator{
		Database: db,
		Logger:   logger,
	}
}

func (d *LoggingDecorator) Exists(model interface{}, query string, args ...interface{}) (bool, error) {
	d.Logger.Debug("Checking existence", "model", model, "query", query, "args", args)
	exists, err := d.Database.Exists(model, query, args...)
	if err != nil {
		d.Logger.Error("Error checking existence", err, "model", model, "query", query, "args", args)
	}
	return exists, err
}

func (d *LoggingDecorator) Create(model interface{}) error {
	d.Logger.Debug("Creating model", "model", model)
	err := d.Database.Create(model)
	if err != nil {
		d.Logger.Error("Error creating model", err, "model", model)
	}
	return err
}

func (d *LoggingDecorator) FindOne(model interface{}, query string, args ...interface{}) error {
	d.Logger.Debug("Finding one", "model", model, "query", query, "args", args)
	err := d.Database.FindOne(model, query, args...)
	if err != nil {
		d.Logger.Error("Error finding one", err, "model", model, "query", query, "args", args)
	}
	return err
}

func (d *LoggingDecorator) Exec(query string, args ...interface{}) error {
	d.Logger.Debug("Executing query", "query", query, "args", args)
	err := d.Database.Exec(query, args...)
	if err != nil {
		d.Logger.Error("Error executing query", err, "query", query, "args", args)
	}
	return err
}

func (d *LoggingDecorator) Query(query string, args ...interface{}) core.Result {
	d.Logger.Debug("Querying", "query", query, "args", args)
	return d.Database.Query(query, args...)
}

// Add Association method implementation
func (d *LoggingDecorator) Association(column string) core.AssociationInterface {
	d.Logger.Debug("Getting association", "column", column)
	return d.Database.Association(column)
}

// Add this method to the LoggingDecorator struct
func (d *LoggingDecorator) AutoMigrate(dst ...interface{}) error {
	d.Logger.Info("Executing AutoMigrate")
	start := time.Now()

	// Call the underlying db's AutoMigrate method
	err := d.Database.AutoMigrate(dst...)

	d.Logger.Info("Completed AutoMigrate",
		"duration", time.Since(start),
		"error", err)

	return err
}

// Add this method to the LoggingDecorator struct
func (d *LoggingDecorator) Delete(model interface{}) error {
	d.Logger.Debug("Deleting model", "model", model)
	err := d.Database.Delete(model)
	if err != nil {
		d.Logger.Error("Error deleting model", err, "model", model)
	}
	return err
}

// Add this method to the LoggingDecorator struct
func (d *LoggingDecorator) Group(name string) core.Interface {
	d.Logger.Debug("Creating database group", "name", name)
	return d.Database.Group(name)
}

// Add this method to the LoggingDecorator struct
func (d *LoggingDecorator) Having(query string, args ...interface{}) core.Interface {
	d.Logger.Debug("Applying Having clause", "query", query, "args", args)
	return d.Database.Having(query, args...)
}

// Add this method to the LoggingDecorator struct
func (d *LoggingDecorator) Joins(query string, args ...interface{}) core.Interface {
	d.Logger.Debug("Applying Joins clause", "query", query, "args", args)
	return d.Database.Joins(query, args...)
}

// Add this method to the LoggingDecorator struct
func (d *LoggingDecorator) Limit(limit int) core.Interface {
	d.Logger.Debug("Applying Limit clause", "limit", limit)
	return d.Database.Limit(limit)
}

// Add this method to the LoggingDecorator struct
func (d *LoggingDecorator) Offset(offset int) core.Interface {
	d.Logger.Debug("Applying Offset clause", "offset", offset)
	return d.Database.Offset(offset)
}

// Add this method to the LoggingDecorator struct
func (d *LoggingDecorator) Not(query interface{}, args ...interface{}) core.Interface {
	d.Logger.Debug("Executing Not query", "query", query, "args", args)
	return d.Database.Not(query, args...)
}

func (d *LoggingDecorator) Or(query interface{}, args ...interface{}) core.Interface {
	d.Logger.Debug("Executing Or query", "query", query, "args", args)
	return d.Database.Or(query, args...)
}

func (d *LoggingDecorator) Order(value interface{}) core.Interface {
	d.Logger.Debug("Applying Order clause", "value", value)
	return d.Database.Order(value)
}

func (d *LoggingDecorator) Preload(query string, args ...interface{}) core.Interface {
	d.Logger.Debug("Applying Preload", "query", query, "args", args)
	return d.Database.Preload(query, args...)
}

func (d *LoggingDecorator) Unscoped() core.Interface {
	d.Logger.Debug("Applying Unscoped")
	return d.Database.Unscoped()
}

func (d *LoggingDecorator) Where(query interface{}, args ...interface{}) core.Interface {
	d.Logger.Debug("Applying Where clause", "query", query, "args", args)
	return d.Database.Where(query, args...)
}

func (d *LoggingDecorator) WithTrashed() core.Interface {
	d.Logger.Debug("Applying WithTrashed")
	return d.Database.WithTrashed()
}

// DB returns the underlying database interface
func (d *LoggingDecorator) DB() *gorm.DB {
	return d.Database.DB()
}

func (d *LoggingDecorator) GetSQLDB() (*sql.DB, error) {
	return d.Database.DB().DB()
}

// Add this method to the LoggingDecorator struct
func (d *LoggingDecorator) Error() string {
	err := d.Database.Error()
	if err != "" {
		d.Logger.Debug("Database error occurred", "error", err)
	}
	return err
}

// Add this method to the LoggingDecorator struct
func (d *LoggingDecorator) Migrator() core.Migrator {
	d.Logger.Debug("Getting database migrator")
	return d.Database.Migrator()
}

// Add this method to the LoggingDecorator struct
func (d *LoggingDecorator) Update(value interface{}) error {
	d.Logger.Debug("Updating model", "model", value)
	err := d.Database.Update(value)
	if err != nil {
		d.Logger.Error("Error updating model", err, "model", value)
	}
	return err
}
