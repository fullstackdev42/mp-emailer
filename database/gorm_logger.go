package database

import (
	"context"
	"errors"
	"time"

	"github.com/jonesrussell/loggo"
	gormlogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	logger loggo.LoggerInterface
}

func NewGormLogger(logger loggo.LoggerInterface) gormlogger.Interface {
	return &GormLogger{
		logger: logger,
	}
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Info(msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Warn(msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Error(msg, errors.New(msg), data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		l.logger.Error("Database query failed",
			err,
			"sql", sql,
			"rows", rows,
			"elapsed", elapsed)
		return
	}

	l.logger.Debug("Database query",
		"sql", sql,
		"rows", rows,
		"elapsed", elapsed)
}
