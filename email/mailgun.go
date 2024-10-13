package email

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type MailgunEmailService struct {
	domain string
	apiKey string
}

func NewMailgunEmailService(domain, apiKey string) *MailgunEmailService {
	return &MailgunEmailService{
		domain: domain,
		apiKey: apiKey,
	}
}

func (s *MailgunEmailService) SendEmail(to, subject, body string) error {
	mg := mailgun.NewMailgun(s.domain, s.apiKey)
	message := mg.NewMessage("no-reply@"+s.domain, subject, body, to)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mg.Send(ctx, message)
	return err
}
