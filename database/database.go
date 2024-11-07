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
	Query(query string, args ...interface{}) *gorm.DB
}

type DB struct {
	GormDB *gorm.DB
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
	result := db.GormDB.Create(value)
	return result.Error
}

func (db *DB) FindOne(model interface{}, query string, args ...interface{}) error {
	return db.GormDB.Where(query, args...).First(model).Error
}

func (db *DB) Exec(query string, args ...interface{}) error {
	return db.GormDB.Exec(query, args...).Error
}

func (db *DB) Query(query string, args ...interface{}) *gorm.DB {
	return db.GormDB.Raw(query, args...)
}
