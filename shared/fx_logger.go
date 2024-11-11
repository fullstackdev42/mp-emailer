package shared

import (
	"github.com/jonesrussell/loggo"
	"go.uber.org/fx/fxevent"
)

// CustomFxLogger wraps the application logger for fx logging
type CustomFxLogger struct {
	logger loggo.LoggerInterface
}

// NewCustomFxLogger creates a new CustomFxLogger instance
func NewCustomFxLogger(logger loggo.LoggerInterface) fxevent.Logger {
	return &CustomFxLogger{logger: logger}
}

func (l *CustomFxLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.logger.Info("Starting", "caller", e.FunctionName)
	case *fxevent.OnStopExecuting:
		l.logger.Info("Stopping", "caller", e.FunctionName)
	case *fxevent.Provided:
		l.logger.Info("Provided constructor",
			"constructor", e.ConstructorName,
			"types", e.OutputTypeNames,
			"module", e.ModuleName,
			"private", e.Private)
		if e.Err != nil {
			l.logger.Error("Constructor provision failed", e.Err)
		}
	}
}
