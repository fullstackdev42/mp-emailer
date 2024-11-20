package database

import (
	"context"

	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/mp-emailer/config"
	dbconfig "github.com/jonesrussell/mp-emailer/database/config"
	"github.com/jonesrussell/mp-emailer/database/core"
	"gorm.io/gorm"
)

// Interface defines the contract for database operations
type Interface interface {
	Connect(ctx context.Context, cfg *config.Config) (*gorm.DB, error)
	Close() error
	GetDB() core.Interface
	WithContext(ctx context.Context) Interface
	Transaction(ctx context.Context, fn func(tx core.Interface) error) error
	Begin(ctx context.Context) (Interface, error)
	Commit() error
	Rollback() error
}

// New creates a new Database instance
func New(logger loggo.LoggerInterface) Interface {
	return &Database{
		logger: logger,
	}
}

// Database implements Interface
type Database struct {
	db     core.Interface
	logger loggo.LoggerInterface
	config *config.Config
}

func (d *Database) Connect(ctx context.Context, cfg *config.Config) (*gorm.DB, error) {
	connector := &dbconfig.DefaultConnector{}
	db, err := dbconfig.ConnectWithRetry(ctx, cfg, dbconfig.NewDefaultRetryConfig(), d.logger, connector)
	if err != nil {
		return nil, err
	}
	d.db = &core.DB{GormDB: db}
	d.config = cfg
	return db, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.db.GetSQLDB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *Database) GetDB() core.Interface {
	return d.db
}

func (d *Database) WithContext(ctx context.Context) Interface {
	newDB := &Database{
		db:     d.db.WithContext(ctx),
		logger: d.logger,
		config: d.config,
	}
	return newDB
}

func (d *Database) Transaction(ctx context.Context, fn func(tx core.Interface) error) error {
	return d.db.Transaction(ctx, fn)
}

func (d *Database) Begin(ctx context.Context) (Interface, error) {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &Database{
		db:     tx,
		logger: d.logger,
		config: d.config,
	}, nil
}

func (d *Database) Commit() error {
	return d.db.Commit()
}

func (d *Database) Rollback() error {
	return d.db.Rollback()
}
