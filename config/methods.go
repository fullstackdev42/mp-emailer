package config

import (
	"github.com/jonesrussell/mp-emailer/version"
	"go.uber.org/zap/zapcore"
)

func (c *Config) GetLogLevel() zapcore.Level {
	switch c.Log.Level {
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

func (c *Config) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":      "ok",
		"environment": c.App.Env,
		"version": map[string]string{
			"version":   c.Version.Version,
			"buildDate": c.Version.BuildDate,
			"commit":    c.Version.Commit,
		},
		"features": map[string]bool{
			"mailgun": c.FeatureFlags.EnableMailgun,
			"smtp":    c.FeatureFlags.EnableSMTP,
			"metrics": c.FeatureFlags.EnableMetrics,
			"beta":    c.FeatureFlags.BetaFeatures,
		},
	}
}

func (c *Config) InitializeVersion() version.Info {
	return version.Info{
		Version:   c.Version.Version,
		BuildDate: c.Version.BuildDate,
		Commit:    c.Version.Commit,
	}
}
