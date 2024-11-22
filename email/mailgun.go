package email

import (
	"context"
	"fmt"
	"time"

	"github.com/jonesrussell/mp-emailer/logger"
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
	Logger logger.Interface
}

func NewMailgunEmailService(domain, apiKey string, client MailgunClient, log logger.Interface) *MailgunEmailService {
	return &MailgunEmailService{
		domain: domain,
		apiKey: apiKey,
		client: client,
		Logger: log,
	}
}

func (s *MailgunEmailService) SendEmail(to, subject, body string, isHTML bool) error {
	message := s.client.NewMessage(
		fmt.Sprintf("no-reply@%s", s.domain),
		subject,
		body,
		to,
	)

	if isHTML {
		s.Logger.Debug("HTML Body content", "body", body)
		message.SetHTML(body)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, id, err := s.client.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.Logger.Debug("Email sent successfully", "messageId", id)
	return nil
}
