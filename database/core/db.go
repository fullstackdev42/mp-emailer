package core

import (
	"database/sql"

	"gorm.io/gorm"
)

type DB struct {
	GormDB *gorm.DB
}

func (db *DB) DB() *gorm.DB {
	return db.GormDB
}

func (db *DB) Exists(model interface{}, query string, args ...interface{}) (bool, error) {
	var count int64
	result := db.GormDB.Model(model).Where(query, args...).Count(&count)
	return count > 0, result.Error
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
	stmt := db.GormDB.Raw(query, args...).Statement
	return NewResult(stmt)
}

func (db *DB) Delete(value interface{}) error {
	return db.GormDB.Delete(value).Error
}

func (db *DB) Unscoped() Interface {
	return &DB{GormDB: db.GormDB.Unscoped()}
}

func (db *DB) WithTrashed() Interface {
	return &DB{GormDB: db.GormDB.Unscoped()}
}

func (db *DB) Preload(query string, args ...interface{}) Interface {
	return &DB{GormDB: db.GormDB.Preload(query, args...)}
}

func (db *DB) Association(column string) AssociationInterface {
	return NewAssociation(db.GormDB.Association(column))
}

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

func (db *DB) AutoMigrate(dst ...interface{}) error {
	return db.GormDB.AutoMigrate(dst...)
}

type GormMigrator struct {
	migrator gorm.Migrator
}

// CreateTable implements Migrator.
func (m *GormMigrator) CreateTable(dst ...interface{}) error {
	return m.migrator.CreateTable(dst...)
}

// DropTable implements Migrator.
func (m *GormMigrator) DropTable(dst ...interface{}) error {
	return m.migrator.DropTable(dst...)
}

// HasTable implements Migrator.
func (m *GormMigrator) HasTable(dst interface{}) bool {
	return m.migrator.HasTable(dst)
}

// Up implements Migrator.
func (m *GormMigrator) Up() error {
	panic("unimplemented")
}

func (m GormMigrator) AutoMigrate(dst ...interface{}) error {
	return m.migrator.AutoMigrate(dst...)
}

func (m GormMigrator) Close() error {
	return nil
}

func (db *DB) Migrator() Migrator {
	return &GormMigrator{migrator: db.GormDB.Migrator()}
}

func (db *DB) GetSQLDB() (*sql.DB, error) {
	return db.GormDB.DB()
}

// Add Error method to DB struct
func (db *DB) Error() string {
	if db.GormDB.Error != nil {
		return db.GormDB.Error.Error()
	}
	return ""
}

// Add this method to your DB struct
func (db *DB) Update(value interface{}) error {
	return db.GormDB.Save(value).Error
}
