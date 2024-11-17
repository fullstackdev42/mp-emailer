package config

// Log is used for logging configuration without sensitive fields
type Log struct {
	*Config
	JWTSecret     string `json:"-"`
	MailgunAPIKey string `json:"-"`
	SessionSecret string `json:"-"`
}

// EmailProvider represents the type of email service to use
type EmailProvider string

const (
	EmailProviderSMTP    EmailProvider = "smtp"
	EmailProviderMailgun EmailProvider = "mailgun"
)

type VersionConfig struct {
	Version   string `yaml:"version" env:"APP_VERSION"`
	BuildDate string `yaml:"build_date" env:"BUILD_DATE"`
	Commit    string `yaml:"commit" env:"GIT_COMMIT"`
}
