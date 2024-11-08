package database

import (
	"gorm.io/gorm"
)

// Interface defines the contract for database operations
type Interface interface {
	Exists(model interface{}, query string, args ...interface{}) (bool, error)
	Create(value interface{}) error
	FindOne(model interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) error
	Query(query string, args ...interface{}) Result
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
