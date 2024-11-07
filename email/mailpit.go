package email

import (
	"fmt"

	"github.com/jonesrussell/loggo"
)

type MailpitEmailService struct {
	host       string
	port       string
	smtpClient SMTPClient
	logger     loggo.LoggerInterface
}

func NewMailpitEmailService(host, port string, smtpClient SMTPClient, logger loggo.LoggerInterface) *MailpitEmailService {
	return &MailpitEmailService{
		host:       host,
		port:       port,
		smtpClient: smtpClient,
		logger:     logger,
	}
}

func (s *MailpitEmailService) SendEmail(to, subject, body string, isHTML bool) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	contentType := "text/plain"
	if isHTML {
		contentType = "text/html"
		s.logger.Debug("HTML Body content", "body", body)
	}

	// Add proper email formatting with MIME version
	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: %s; charset=UTF-8\r\n"+
		"\r\n"+
		"%s",
		to, subject, contentType, body))

	s.logger.Debug("Full message", "message", string(message))
	return s.smtpClient.SendMail(addr, nil, "no-reply@example.com", []string{to}, message)
}
