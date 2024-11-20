package email

import (
	"fmt"
)

type Service interface {
	SendEmail(to string, subject string, body string, isHTML bool) error
	SendPasswordReset(to string, resetToken string) error
}

// Ensure both services implement the interface
var (
	_ Service = (*MailpitEmailService)(nil)
	_ Service = (*MailgunEmailService)(nil)
)

// Add any missing methods to ensure both services fully implement the interface

func (s *MailpitEmailService) SendPasswordReset(to string, resetToken string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf(`Hello,

A password reset has been requested for your account. 
To reset your password, please use the following token: %s

If you did not request this reset, please ignore this email.

Best regards,
Your Application Team`, resetToken)

	return s.SendEmail(to, subject, body, false)
}

func (s *MailgunEmailService) SendPasswordReset(to string, resetToken string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf(`Hello,

A password reset has been requested for your account. 
To reset your password, please use the following token: %s

If you did not request this reset, please ignore this email.

Best regards,
Your Application Team`, resetToken)

	return s.SendEmail(to, subject, body, false)
}
