package shared

import (
	"github.com/jonesrussell/mp-emailer/logger"
	"go.uber.org/fx/fxevent"
)

// CustomFxLogger wraps the application logger for fx logging
type CustomFxLogger struct {
	logger logger.Interface
}

// NewCustomFxLogger creates a new CustomFxLogger instance
func NewCustomFxLogger(logger logger.Interface) fxevent.Logger {
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
