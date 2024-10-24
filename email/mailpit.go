package email

import (
	"fmt"
	"net/smtp"

	"github.com/jonesrussell/loggo"
)

type MailpitEmailService struct {
	host       string
	port       string
	smtpClient SMTPClient
	logger     *loggo.Logger
}

type SMTPClient interface {
	SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

type smtpClientWrapper struct{}

func (s smtpClientWrapper) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return smtp.SendMail(addr, a, from, to, msg)
}

func NewMailpitEmailService(host, port string, logger loggo.LoggerInterface) (Service, error) {
	if host == "" || port == "" {
		return nil, fmt.Errorf("invalid Mailpit configuration: host and port must not be empty")
	}

	return &MailpitEmailService{
		host:       host,
		port:       port,
		smtpClient: smtpClientWrapper{},
		logger:     logger.(*loggo.Logger),
	}, nil
}

func (s *MailpitEmailService) SendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	message := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))
	return s.smtpClient.SendMail(addr, nil, "no-reply@example.com", []string{to}, message)
}
