package email

import (
	"fmt"
	"net/smtp"
)

type MailpitEmailService struct {
	host       string
	port       string
	smtpClient SMTPClient
}

type SMTPClient interface {
	SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

func (s *MailpitEmailService) SendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	message := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))
	return s.smtpClient.SendMail(addr, nil, "no-reply@example.com", []string{to}, message)
}

func NewMailpitEmailService(host, port string, smtpClient SMTPClient) *MailpitEmailService {
	return &MailpitEmailService{
		host:       host,
		port:       port,
		smtpClient: smtpClient,
	}
}
