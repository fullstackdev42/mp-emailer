package logger

// Interface defines the standard logging operations
type Interface interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, err error, fields ...interface{})
	Fatal(msg string, err error, fields ...interface{})
	IsDebugEnabled() bool
	With(fields ...interface{}) Interface
	WithOperation(operation string) Interface
}
