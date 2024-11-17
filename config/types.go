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
