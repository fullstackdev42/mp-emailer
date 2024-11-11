package database

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/jonesrussell/loggo"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Interface defines the contract for database operations
type Interface interface {
	Exists(model interface{}, query string, args ...interface{}) (bool, error)
	Create(value interface{}) error
	FindOne(model interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) error
	Query(query string, args ...interface{}) Result
	Delete(value interface{}) error
	Unscoped() Interface
	WithTrashed() Interface
	Preload(query string, args ...interface{}) Interface
	Association(column string) AssociationInterface
	Where(query interface{}, args ...interface{}) Interface
	Or(query interface{}, args ...interface{}) Interface
	Not(query interface{}, args ...interface{}) Interface
	Order(value interface{}) Interface
	Group(query string) Interface
	Having(query string, args ...interface{}) Interface
	Joins(query string, args ...interface{}) Interface
	Limit(limit int) Interface
	Offset(offset int) Interface
	AutoMigrate(dst ...interface{}) error
	Migrator() Migrator
}

// Result represents the interface for database query results
type Result interface {
	Scan(dest interface{}) Result
	Error() error
}

type DB struct {
	GormDB *gorm.DB
}

// gormResult wraps gorm.DB to implement the Result interface
type gormResult struct {
	db *gorm.DB
}

func (r *gormResult) Scan(dest interface{}) Result {
	r.db = r.db.Scan(dest)
	return r
}

func (r *gormResult) Error() error {
	return r.db.Error
}

var _ Interface = (*DB)(nil)

// NewDB creates a new database instance with the provided gorm.DB
func NewDB(gormDB *gorm.DB) Interface {
	if gormDB == nil {
		panic("gormDB cannot be nil")
	}
	return &DB{
		GormDB: gormDB,
	}
}

// Implement the generic interface
func (db *DB) Exists(model interface{}, query string, args ...interface{}) (bool, error) {
	var count int64
	err := db.GormDB.Model(model).Where(query, args...).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (db *DB) Create(value interface{}) error {
	return db.GormDB.Create(value).Error
}

func (db *DB) FindOne(model interface{}, query string, args ...interface{}) error {
	return db.GormDB.Where(query, args...).First(model).Error
}

func (db *DB) Exec(query string, args ...interface{}) error {
	return db.GormDB.Exec(query, args...).Error
}

func (db *DB) Query(query string, args ...interface{}) Result {
	return &gormResult{db: db.GormDB.Raw(query, args...)}
}

func (db *DB) Delete(value interface{}) error {
	return db.GormDB.Delete(value).Error
}

func (db *DB) Unscoped() Interface {
	return &DB{GormDB: db.GormDB.Unscoped()}
}

func (db *DB) WithTrashed() Interface {
	return db
}

func (db *DB) Preload(query string, args ...interface{}) Interface {
	return &DB{GormDB: db.GormDB.Preload(query, args...)}
}

func (db *DB) Association(column string) AssociationInterface {
	return &gormAssociation{association: db.GormDB.Association(column)}
}

type AssociationInterface interface {
	Find(out interface{}) error
	Append(values ...interface{}) error
	Replace(values ...interface{}) error
	Delete(values ...interface{}) error
	Clear() error
	Count() int64
}

type gormAssociation struct {
	association *gorm.Association
}

func (a *gormAssociation) Find(out interface{}) error {
	return a.association.Find(out)
}

func (a *gormAssociation) Append(values ...interface{}) error {
	return a.association.Append(values...)
}

func (a *gormAssociation) Replace(values ...interface{}) error {
	return a.association.Replace(values...)
}

func (a *gormAssociation) Delete(values ...interface{}) error {
	return a.association.Delete(values...)
}

func (a *gormAssociation) Clear() error {
	return a.association.Clear()
}

func (a *gormAssociation) Count() int64 {
	return a.association.Count()
}

func (db *DB) AutoMigrate(dst ...interface{}) error {
	return db.GormDB.AutoMigrate(dst...)
}

func (db *DB) Migrator() Migrator {
	gormMigrator := db.GormDB.Migrator()
	return &customMigrator{migrator: gormMigrator}
}

func (db *DB) Group(query string) Interface {
	return &DB{GormDB: db.GormDB.Group(query)}
}

func (db *DB) Having(query string, args ...interface{}) Interface {
	return &DB{GormDB: db.GormDB.Having(query, args...)}
}

// Add custom migrator type
type customMigrator struct {
	migrator gorm.Migrator
}

func (m *customMigrator) AutoMigrate(dst ...interface{}) error {
	return m.migrator.AutoMigrate(dst...)
}

func (m *customMigrator) Close() error {
	// Implement any cleanup if needed
	return nil
}

// Add Joins method implementation
func (db *DB) Joins(query string, args ...interface{}) Interface {
	return &DB{GormDB: db.GormDB.Joins(query, args...)}
}

func (db *DB) Limit(limit int) Interface {
	return &DB{GormDB: db.GormDB.Limit(limit)}
}

func (db *DB) Offset(offset int) Interface {
	return &DB{GormDB: db.GormDB.Offset(offset)}
}

func (db *DB) Not(query interface{}, args ...interface{}) Interface {
	return &DB{GormDB: db.GormDB.Not(query, args...)}
}

func (db *DB) Order(value interface{}) Interface {
	return &DB{GormDB: db.GormDB.Order(value)}
}

func (db *DB) Or(query interface{}, args ...interface{}) Interface {
	return &DB{GormDB: db.GormDB.Or(query, args...)}
}

// Add missing methods to customMigrator
func (m *customMigrator) CreateTable(dst ...interface{}) error {
	return m.migrator.CreateTable(dst...)
}

func (m *customMigrator) DropTable(dst ...interface{}) error {
	return m.migrator.DropTable(dst...)
}

func (m *customMigrator) HasTable(dst interface{}) bool {
	return m.migrator.HasTable(dst)
}

// Add Where method implementation
func (db *DB) Where(query interface{}, args ...interface{}) Interface {
	return &DB{GormDB: db.GormDB.Where(query, args...)}
}

// Add Up method to customMigrator
func (m *customMigrator) Up() error {
	// Since GORM's migrator doesn't have an Up method directly,
	// we'll treat AutoMigrate as our "up" operation
	return m.AutoMigrate()
}

type RetryConfig struct {
	InitialInterval      time.Duration
	MaxInterval          time.Duration
	MaxElapsedTime       time.Duration
	MultiplicationFactor float64
}

func NewDefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		InitialInterval:      100 * time.Millisecond,
		MaxInterval:          10 * time.Second,
		MaxElapsedTime:       1 * time.Minute,
		MultiplicationFactor: 2.0,
	}
}

func ConnectWithRetry(cfg *config.Config, retryConfig *RetryConfig, logger loggo.LoggerInterface) (*gorm.DB, error) {
	var db *gorm.DB

	operation := func() error {
		var err error
		db, err = Connect(cfg)
		if err != nil {
			logger.Error("Failed to connect to database", err)
			return err
		}
		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = retryConfig.InitialInterval
	expBackoff.MaxInterval = retryConfig.MaxInterval
	expBackoff.MaxElapsedTime = retryConfig.MaxElapsedTime
	expBackoff.Multiplier = retryConfig.MultiplicationFactor

	err := backoff.Retry(operation, expBackoff)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after retries: %w", err)
	}

	logger.Info("Successfully connected to database after retry")
	return db, nil
}

// Update the existing fx provider
func ProvideDatabase(lc fx.Lifecycle, cfg *config.Config, logger loggo.LoggerInterface) (*gorm.DB, error) {
	retryConfig := NewDefaultRetryConfig()
	db, err := ConnectWithRetry(cfg, retryConfig, logger)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			return sqlDB.Close()
		},
	})

	return db, nil
}

func Connect(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
