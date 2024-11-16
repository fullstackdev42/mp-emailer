package config

import (
	"fmt"
	"path/filepath"
	"time"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database" env:"sensitive"`
	Email    EmailConfig    `yaml:"email" env:"sensitive"`
	Auth     AuthConfig     `yaml:"auth" env:"sensitive"`
	Log      LogConfig      `yaml:"log"`
	Server   ServerConfig   `yaml:"server"`
}

type AppConfig struct {
	Debug bool        `env:"APP_DEBUG" envDefault:"false"`
	Env   Environment `env:"APP_ENV" envDefault:"development"`
	Host  string      `env:"APP_HOST" envDefault:"0.0.0.0"`
	Port  int         `env:"APP_PORT" envDefault:"8080"`
}

type DatabaseConfig struct {
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT" envDefault:"3306"`
	Name     string `env:"DB_NAME,required"`
}

type EmailConfig struct {
	Provider      EmailProvider `env:"EMAIL_PROVIDER" envDefault:"smtp"`
	MailgunKey    string        `env:"MAILGUN_KEY,required"`
	MailgunDomain string        `env:"MAILGUN_DOMAIN,required"`
	SMTP          SMTPConfig
}

type SMTPConfig struct {
	From     string `env:"SMTP_FROM,required"`
	Host     string `env:"SMTP_HOST,required"`
	Password string `env:"SMTP_PASSWORD,required"`
	Port     int    `env:"SMTP_PORT" envDefault:"587"`
	Username string `env:"SMTP_USERNAME,required"`
}

type AuthConfig struct {
	JWTExpiry     string `env:"JWT_EXPIRY" envDefault:"24h"`
	JWTSecret     string `env:"JWT_SECRET,required"`
	SessionName   string `env:"SESSION_NAME" envDefault:"session"`
	SessionSecret string `env:"SESSION_SECRET,required"`
}

type LogConfig struct {
	File  string
	Level string
}

type ServerConfig struct {
	MigrationsPath              string
	RepresentativeLookupBaseURL string
}

// DSN returns the database connection string
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)
}

// GetAbsolutePath returns the absolute path of a given path
func (c *Config) GetAbsolutePath(path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}
	return abs
}

// GetMigrationsPath returns the path to the migrations directory
func (c *Config) GetMigrationsPath() string {
	return c.Server.MigrationsPath
}

// GetLogFilePath returns the path to the log file
func (c *Config) GetLogFilePath() string {
	return c.Log.File
}

// GetJWTExpiryDuration returns the parsed JWT expiry duration
func (c *Config) GetJWTExpiryDuration() (time.Duration, error) {
	return time.ParseDuration(c.Auth.JWTExpiry)
}

func (c *Config) setupPaths() error {
	// Get absolute path for migrations
	migrationsPath, err := filepath.Abs(c.Server.MigrationsPath)
	if err != nil {
		return fmt.Errorf("invalid migrations path: %w", err)
	}
	c.Server.MigrationsPath = filepath.Clean(migrationsPath)

	// Get absolute path for log file
	logPath, err := filepath.Abs(c.Log.File)
	if err != nil {
		return fmt.Errorf("invalid log file path: %w", err)
	}
	c.Log.File = filepath.Clean(logPath)

	return nil
}
