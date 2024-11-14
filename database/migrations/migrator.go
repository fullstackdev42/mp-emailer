package migrations

import (
	"gorm.io/gorm"
)

// GormMigrator wraps GORM's migrator to implement our Migrator interface
type GormMigrator struct {
	migrator gorm.Migrator
}

func (m *GormMigrator) Up() error {
	// This is a no-op since GORM handles migrations differently
	return nil
}

func (m *GormMigrator) Close() error {
	// No explicit close needed for GORM migrator
	return nil
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
