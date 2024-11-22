package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func NewZapLogger(cfg *Config) (*Logger, error) {
	zapCfg := zap.NewProductionConfig()

	// Configure based on settings
	zapCfg.OutputPaths = []string{cfg.OutputPath}
	zapCfg.Level = zap.NewAtomicLevelAt(getZapLevel(cfg.Level))

	if cfg.Format == "text" {
		zapCfg.Encoding = "console"
	}

	if cfg.Development {
		zapCfg = zap.NewDevelopmentConfig()
	}

	logger, err := zapCfg.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return &Logger{logger}, nil
}

// Interface implementation
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.Logger.Debug(msg, toZapFields(fields...)...)
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	l.Logger.Info(msg, toZapFields(fields...)...)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.Logger.Warn(msg, toZapFields(fields...)...)
}

func (l *Logger) Error(msg string, err error, fields ...interface{}) {
	zapFields := toZapFields(fields...)
	if err != nil {
		zapFields = append(zapFields, zap.Error(err))
	}
	l.Logger.Error(msg, zapFields...)
}

func (l *Logger) Fatal(msg string, err error, fields ...interface{}) {
	zapFields := toZapFields(fields...)
	if err != nil {
		zapFields = append(zapFields, zap.Error(err))
	}
	l.Logger.Fatal(msg, zapFields...)
}

func (l *Logger) With(fields ...interface{}) Interface {
	return &Logger{l.Logger.With(toZapFields(fields...)...)}
}

func (l *Logger) IsDebugEnabled() bool {
	return l.Core().Enabled(zapcore.DebugLevel)
}

func (l *Logger) WithOperation(operation string) Interface {
	return &Logger{l.Logger.With(zap.String("operation", operation))}
}

// Helper functions
func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func toZapFields(fields ...interface{}) []zap.Field {
	if len(fields)%2 != 0 {
		return []zap.Field{zap.Error(fmt.Errorf("odd number of fields provided"))}
	}

	zapFields := make([]zap.Field, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		zapFields = append(zapFields, zap.Any(key, fields[i+1]))
	}
	return zapFields
}
