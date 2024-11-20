package email

import (
	"fmt"
	"net/smtp"

	"github.com/jonesrussell/mp-emailer/config"
	"github.com/mailgun/mailgun-go/v4"
)

// Provider represents supported email service providers
type Provider = config.EmailProvider

const (
	ProviderSMTP    Provider = config.EmailProviderSMTP
	ProviderMailgun Provider = config.EmailProviderMailgun
)

// Config holds the configuration needed for email services
type Config struct {
	Provider      Provider `env:"EMAIL_PROVIDER" envDefault:"smtp"`
	SMTPHost      string   `env:"SMTP_HOST"`
	SMTPPort      int      `env:"SMTP_PORT" envDefault:"587"`
	SMTPUsername  string   `env:"SMTP_USERNAME"`
	SMTPPassword  string   `env:"SMTP_PASSWORD"`
	SMTPFrom      string   `env:"SMTP_FROM"`
	MailgunDomain string   `env:"MAILGUN_DOMAIN"`
	MailgunAPIKey string   `env:"MAILGUN_API_KEY"`
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
func NewEmailService(p Params) (Service, error) {
	switch p.Config.Provider {
	case ProviderSMTP:
		if p.Config.SMTPHost == "" {
			return nil, fmt.Errorf("SMTP configuration is incomplete")
		}

		smtpClient := &SMTPClientImpl{
			auth: smtp.PlainAuth("", p.Config.SMTPUsername, p.Config.SMTPPassword, p.Config.SMTPHost),
		}

		return NewMailpitEmailService(
			p.Config.SMTPHost,
			fmt.Sprintf("%d", p.Config.SMTPPort),
			smtpClient,
			p.Config.SMTPFrom,
		), nil

	case ProviderMailgun:
		if p.Config.MailgunDomain == "" || p.Config.MailgunAPIKey == "" {
			return nil, fmt.Errorf("Mailgun configuration is incomplete")
		}

		mg := mailgun.NewMailgun(p.Config.MailgunDomain, p.Config.MailgunAPIKey)

		return NewMailgunEmailService(
			p.Config.MailgunDomain,
			p.Config.MailgunAPIKey,
			mg,
			p.Logger,
		), nil

	default:
		return nil, fmt.Errorf("unsupported email provider: %s", p.Config.Provider)
	}
}
