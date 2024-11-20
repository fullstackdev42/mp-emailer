package core

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

// Reader handles database read operations
type Reader interface {
	FindOne(ctx context.Context, model interface{}, query string, args ...interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) Result
	Exists(ctx context.Context, model interface{}, query string, args ...interface{}) (bool, error)
	Preload(query string, args ...interface{}) Interface
}

// Writer handles database write operations
type Writer interface {
	Create(ctx context.Context, value interface{}) error
	Update(ctx context.Context, value interface{}) error
	Delete(ctx context.Context, value interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) error
}

// QueryBuilder handles query construction
type QueryBuilder interface {
	Where(query interface{}, args ...interface{}) Interface
	Or(query interface{}, args ...interface{}) Interface
	Not(query interface{}, args ...interface{}) Interface
	Order(value interface{}) Interface
	Group(query string) Interface
	Having(query string, args ...interface{}) Interface
	Joins(query string, args ...interface{}) Interface
	Limit(limit int) Interface
	Offset(offset int) Interface
}

// TransactionManager handles database transactions
type TransactionManager interface {
	Transaction(ctx context.Context, fn func(tx Interface) error) error
	Begin(ctx context.Context) (Interface, error)
	Commit() error
	Rollback() error
}

// Interface combines all database operations
type Interface interface {
	Reader
	Writer
	QueryBuilder
	TransactionManager
	DB() *gorm.DB
	WithContext(ctx context.Context) Interface
	Unscoped() Interface
	WithTrashed() Interface
	Association(column string) AssociationInterface
	AutoMigrate(dst ...interface{}) error
	Migrator() Migrator
	GetSQLDB() (*sql.DB, error)
	Error() string
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
