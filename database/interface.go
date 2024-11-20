package database

import (
	"context"

	"gorm.io/gorm"
)

type Database interface {
	// Core operations
	Create(ctx context.Context, value interface{}) error
	FindOne(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Update(ctx context.Context, value interface{}) error
	Delete(ctx context.Context, value interface{}) error

	// Transaction support
	Transaction(ctx context.Context, fn func(tx Database) error) error

	// Utility methods
	Close() error
	DB() *gorm.DB
}
