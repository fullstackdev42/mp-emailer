package database

import (
	"context"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormDB struct {
	db *gorm.DB
}

func New(dsn string) (Database, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &GormDB{db: db}, nil
}

func (g *GormDB) Create(ctx context.Context, value interface{}) error {
	return g.db.WithContext(ctx).Create(value).Error
}

func (g *GormDB) FindOne(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return g.db.WithContext(ctx).Where(query, args...).First(dest).Error
}

func (g *GormDB) Update(ctx context.Context, value interface{}) error {
	return g.db.WithContext(ctx).Save(value).Error
}

func (g *GormDB) Delete(ctx context.Context, value interface{}) error {
	return g.db.WithContext(ctx).Delete(value).Error
}

func (g *GormDB) Transaction(ctx context.Context, fn func(tx Database) error) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&GormDB{db: tx})
	})
}

func (g *GormDB) Close() error {
	sqlDB, err := g.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (g *GormDB) DB() *gorm.DB {
	return g.db
}
