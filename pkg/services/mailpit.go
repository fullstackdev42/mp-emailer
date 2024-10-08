package services

import (
	"fmt"
	"net/smtp"
)

type MailpitEmailService struct {
	host string
	port string
}

func NewMailpitEmailService(host, port string) *MailpitEmailService {
	return &MailpitEmailService{
		host: host,
		port: port,
	}
}

func (s *MailpitEmailService) SendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	message := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	return smtp.SendMail(addr, nil, "no-reply@example.com", []string{to}, message)
}
