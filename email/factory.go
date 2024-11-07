package email

import (
	"fmt"
	"net/smtp"

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
func NewEmailService(p Params) (Service, error) {
	switch p.Config.Provider {
	case ProviderSMTP:
		if p.Config.SMTPHost == "" || p.Config.SMTPPort == "" {
			return nil, fmt.Errorf("SMTP configuration is incomplete")
		}

		smtpClient := &SMTPClientImpl{
			auth: smtp.PlainAuth("", p.Config.SMTPUsername, p.Config.SMTPPassword, p.Config.SMTPHost),
		}

		return NewMailpitEmailService(
			p.Config.SMTPHost,
			p.Config.SMTPPort,
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
