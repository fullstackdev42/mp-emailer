package email

import (
	"fmt"
	"net/smtp"

	"github.com/jonesrussell/loggo"
	"github.com/mailgun/mailgun-go/v4"
)

// Provider represents supported email service providers
type Provider string

const (
	ProviderSMTP    Provider = "smtp"
	ProviderMailgun Provider = "mailgun"
)

// Config holds the configuration needed for email services
type Config struct {
	Provider      Provider
	SMTPHost      string
	SMTPPort      string
	SMTPUsername  string
	SMTPPassword  string
	SMTPFrom      string
	MailgunDomain string
	MailgunAPIKey string
}

// SMTPClient interface
type SMTPClient interface {
	SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

// SMTPClientImpl implements SMTPClient
type SMTPClientImpl struct {
	auth smtp.Auth
}

// SendMail implements the SMTPClient interface
func (s *SMTPClientImpl) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return smtp.SendMail(addr, a, from, to, msg)
}

// NewEmailService creates an email service based on the provided configuration
func NewEmailService(config Config) (Service, error) {
	switch config.Provider {
	case ProviderSMTP:
		if config.SMTPHost == "" || config.SMTPPort == "" {
			return nil, fmt.Errorf("SMTP configuration is incomplete")
		}

		smtpClient := &SMTPClientImpl{
			auth: smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost),
		}

		return NewMailpitEmailService(
			config.SMTPHost,
			config.SMTPPort,
			smtpClient,
		), nil

	case ProviderMailgun:
		if config.MailgunDomain == "" || config.MailgunAPIKey == "" {
			return nil, fmt.Errorf("Mailgun configuration is incomplete")
		}

		mg := mailgun.NewMailgun(config.MailgunDomain, config.MailgunAPIKey)
		logger, err := loggo.NewLogger("mailgun", loggo.LevelInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to create logger: %w", err)
		}

		return NewMailgunEmailService(
			config.MailgunDomain,
			config.MailgunAPIKey,
			mg,
			logger,
		), nil

	default:
		return nil, fmt.Errorf("unsupported email provider: %s", config.Provider)
	}
}
