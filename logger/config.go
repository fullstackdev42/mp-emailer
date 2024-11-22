package logger

import "fmt"

type Config struct {
	Level       string `yaml:"level" env:"LOG_LEVEL" envDefault:"info"`
	Format      string `yaml:"format" env:"LOG_FORMAT" envDefault:"json"`
	OutputPath  string `yaml:"file" env:"LOG_FILE" envDefault:"storage/logs/app.log"`
	Development bool   `yaml:"development" env:"LOG_DEVELOPMENT" envDefault:"false"`
}

func (c *Config) Validate() error {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLevels[c.Level] {
		return fmt.Errorf("invalid log level: %s", c.Level)
	}

	validFormats := map[string]bool{
		"json": true,
		"text": true,
	}

	if !validFormats[c.Format] {
		return fmt.Errorf("invalid log format: %s", c.Format)
	}

	return nil
}
