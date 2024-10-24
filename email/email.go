package email

import (
	"fmt"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/jonesrussell/loggo"
)

type Service interface {
	SendEmail(to, subject, body string) error
}

func New(config *config.Config, logger loggo.LoggerInterface) (Service, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if config.AppEnv == "production" {
		if config.MailgunDomain == "" || config.MailgunAPIKey == "" {
			return nil, fmt.Errorf("Mailgun configuration is incomplete")
		}
		return NewMailgunEmailService(config.MailgunDomain, config.MailgunAPIKey, logger)
	}

	if config.MailpitHost == "" || config.MailpitPort == "" {
		return nil, fmt.Errorf("Mailpit configuration is incomplete")
	}
	return NewMailpitEmailService(config.MailpitHost, config.MailpitPort, logger)
}
