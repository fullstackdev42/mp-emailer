package core

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

var _ Interface = (*DB)(nil)

type DB struct {
	GormDB *gorm.DB
}

func (db *DB) DB() *gorm.DB {
	return db.GormDB
}

// Read operations with context
func (db *DB) Exists(ctx context.Context, model interface{}, query string, args ...interface{}) (bool, error) {
	var count int64
	result := db.GormDB.WithContext(ctx).Model(model).Where(query, args...).Count(&count)
	return count > 0, result.Error
}

func (db *DB) FindOne(ctx context.Context, model interface{}, query string, args ...interface{}) error {
	return db.GormDB.WithContext(ctx).Where(query, args...).First(model).Error
}

func (db *DB) Query(ctx context.Context, query string, args ...interface{}) Result {
	stmt := db.GormDB.WithContext(ctx).Raw(query, args...).Statement
	return NewResult(stmt)
}

// Write operations with context
func (db *DB) Create(ctx context.Context, value interface{}) error {
	return db.GormDB.WithContext(ctx).Create(value).Error
}

func (db *DB) Update(ctx context.Context, value interface{}) error {
	return db.GormDB.WithContext(ctx).Save(value).Error
}

func (db *DB) Delete(ctx context.Context, value interface{}) error {
	return db.GormDB.WithContext(ctx).Delete(value).Error
}

func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) error {
	return db.GormDB.WithContext(ctx).Exec(query, args...).Error
}

// Query builder methods
func (db *DB) Where(query interface{}, args ...interface{}) Interface {
	return &DB{GormDB: db.GormDB.Where(query, args...)}
}

func (db *DB) Or(query interface{}, args ...interface{}) Interface {
	return &DB{GormDB: db.GormDB.Or(query, args...)}
}

func (db *DB) Not(query interface{}, args ...interface{}) Interface {
	return &DB{GormDB: db.GormDB.Not(query, args...)}
}

func (db *DB) Order(value interface{}) Interface {
	return &DB{GormDB: db.GormDB.Order(value)}
}

func (db *DB) Group(query string) Interface {
	return &DB{GormDB: db.GormDB.Group(query)}
}

func (db *DB) Having(query string, args ...interface{}) Interface {
	return &DB{GormDB: db.GormDB.Having(query, args...)}
}

func (db *DB) Joins(query string, args ...interface{}) Interface {
	return &DB{GormDB: db.GormDB.Joins(query, args...)}
}

func (db *DB) Limit(limit int) Interface {
	return &DB{GormDB: db.GormDB.Limit(limit)}
}

func (db *DB) Offset(offset int) Interface {
	return &DB{GormDB: db.GormDB.Offset(offset)}
}

// Transaction methods
func (db *DB) Transaction(ctx context.Context, fn func(tx Interface) error) error {
	return db.GormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&DB{GormDB: tx})
	})
}

func (db *DB) Begin(ctx context.Context) (Interface, error) {
	tx := db.GormDB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &DB{GormDB: tx}, nil
}

func (db *DB) Commit() error {
	return db.GormDB.Commit().Error
}

func (db *DB) Rollback() error {
	return db.GormDB.Rollback().Error
}

// Context handling
func (db *DB) WithContext(ctx context.Context) Interface {
	return &DB{GormDB: db.GormDB.WithContext(ctx)}
}

// Utility methods
func (db *DB) GetSQLDB() (*sql.DB, error) {
	return db.GormDB.DB()
}

func (db *DB) Error() string {
	if db.GormDB.Error != nil {
		return db.GormDB.Error.Error()
	}
	return ""
}

// Migration methods
func (db *DB) Migrator() Migrator {
	return &GormMigrator{migrator: db.GormDB.Migrator()}
}

func (db *DB) AutoMigrate(dst ...interface{}) error {
	return db.GormDB.AutoMigrate(dst...)
}

// Add Association method
func (db *DB) Association(column string) AssociationInterface {
	return NewAssociation(db.GormDB.Association(column))
}

// Add Preload method
func (db *DB) Preload(query string, args ...interface{}) Interface {
	return &DB{GormDB: db.GormDB.Preload(query, args...)}
}

// Add Unscoped method
func (db *DB) Unscoped() Interface {
	return &DB{GormDB: db.GormDB.Unscoped()}
}

// Add WithTrashed method
func (db *DB) WithTrashed() Interface {
	return &DB{GormDB: db.GormDB.Unscoped()}
}

// Update GormMigrator implementation
type GormMigrator struct {
	migrator gorm.Migrator
}

func (m *GormMigrator) Up() error {
	return nil // Implement as needed
}

func (m *GormMigrator) Close() error {
	return nil // Implement as needed
}

func (m *GormMigrator) CreateTable(dst ...interface{}) error {
	return m.migrator.CreateTable(dst...)
}

func (m *GormMigrator) DropTable(dst ...interface{}) error {
	return m.migrator.DropTable(dst...)
}

func (m *GormMigrator) HasTable(dst interface{}) bool {
	return m.migrator.HasTable(dst)
}

func (m *GormMigrator) AutoMigrate(dst ...interface{}) error {
	return m.migrator.AutoMigrate(dst...)
}
