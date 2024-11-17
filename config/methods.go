package config

import (
	"github.com/fullstackdev42/mp-emailer/version"
	"github.com/jonesrussell/loggo"
)

func (c *Config) GetLogLevel() loggo.Level {
	switch c.Log.Level {
	case "debug":
		return loggo.LevelDebug
	case "info":
		return loggo.LevelInfo
	case "warn":
		return loggo.LevelWarn
	case "error":
		return loggo.LevelError
	default:
		return loggo.LevelInfo
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
