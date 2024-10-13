package email

import "github.com/fullstackdev42/mp-emailer/config"

type Service interface {
	SendEmail(to, subject, body string) error
}

func NewEmailService(config *config.Config) Service {
	if config.AppEnv == "production" {
		return NewMailgunEmailService(config.MailgunDomain, config.MailgunAPIKey)
	}
	return NewMailpitEmailService(config.MailpitHost, config.MailpitPort)
}
