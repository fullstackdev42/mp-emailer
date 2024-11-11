package core

import (
	"database/sql"

	"gorm.io/gorm"
)

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
	DB() *gorm.DB
	GetSQLDB() (*sql.DB, error)
}

type Result interface {
	Scan(dest interface{}) Result
	Error() error
}

type AssociationInterface interface {
	Find(out interface{}) error
	Append(values ...interface{}) error
	Replace(values ...interface{}) error
	Delete(values ...interface{}) error
	Clear() error
	Count() int64
}

type Migrator interface {
	Up() error
	Close() error
	CreateTable(dst ...interface{}) error
	DropTable(dst ...interface{}) error
	HasTable(dst interface{}) bool
	AutoMigrate(dst ...interface{}) error
}
