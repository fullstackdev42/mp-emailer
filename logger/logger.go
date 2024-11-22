package logger

import (
	"sync"
)

var (
	instance Interface
	once     sync.Once
)

// GetLogger returns the singleton logger instance
func GetLogger() Interface {
	if instance == nil {
		panic("logger not initialized - call Initialize() first")
	}
	return instance
}

// Initialize sets up the logger with the given configuration
func Initialize(cfg *Config) error {
	var err error
	once.Do(func() {
		if err = cfg.Validate(); err != nil {
			return
		}

		var logger *Logger
		logger, err = NewZapLogger(cfg)
		if err != nil {
			return
		}

		instance = logger

		// Log initialization success
		instance.Info("logger initialized successfully",
			"level", cfg.Level,
			"format", cfg.Format,
			"output", cfg.OutputPath,
		)
	})
	return err
}

// Cleanup performs any necessary cleanup of the logger
func Cleanup() error {
	if l, ok := instance.(*Logger); ok {
		return l.Sync()
	}
	return nil
}
