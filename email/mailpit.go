package email

import (
	"fmt"
)

type MailpitEmailService struct {
	host       string
	port       string
	smtpClient SMTPClient
	from       string
}

func NewMailpitEmailService(host, port string, smtpClient SMTPClient, from string) *MailpitEmailService {
	return &MailpitEmailService{
		host:       host,
		port:       port,
		smtpClient: smtpClient,
		from:       from,
	}
}

func (s *MailpitEmailService) SendEmail(to, subject, body string, isHTML bool) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	contentType := "text/plain"
	if isHTML {
		contentType = "text/html"
	}

	message := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: %s; charset=UTF-8\r\n"+
		"\r\n"+
		"%s",
		s.from, to, subject, contentType, body))

	return s.smtpClient.SendMail(addr, nil, s.from, []string{to}, message)
}
