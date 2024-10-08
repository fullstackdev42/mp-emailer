package services

type EmailService interface {
	SendEmail(to, subject, body string) error
}
