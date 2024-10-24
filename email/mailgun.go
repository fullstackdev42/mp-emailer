package email

import (
	"context"
	"fmt"
	"time"

	"github.com/jonesrussell/loggo"
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
	logger loggo.LoggerInterface
}

func NewMailgunEmailService(domain, apiKey string, client MailgunClient, logger loggo.LoggerInterface) *MailgunEmailService {
	return &MailgunEmailService{
		domain: domain,
		apiKey: apiKey,
		client: client,
		logger: logger,
	}
}

func (s *MailgunEmailService) SendEmail(to, subject, body string) error {
	message := s.client.NewMessage(
		fmt.Sprintf("no-reply@%s", s.domain),
		subject,
		body,
		to,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := s.client.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
