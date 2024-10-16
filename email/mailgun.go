package email

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type MailgunClient interface {
	NewMessage(from, subject, text string, to ...string) *mailgun.Message
	Send(ctx context.Context, message *mailgun.Message) (string, string, error)
}

type MailgunEmailService struct {
	domain string
	apiKey string
	client MailgunClient
}

func NewMailgunEmailService(domain, apiKey string) *MailgunEmailService {
	return &MailgunEmailService{
		domain: domain,
		apiKey: apiKey,
	}
}

func (s *MailgunEmailService) SendEmail(to, subject, body string) error {
	message := s.client.NewMessage("no-reply@"+s.domain, subject, body, to)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := s.client.Send(ctx, message)
	return err
}
