package email

type Service interface {
	SendEmail(to, subject, body string) error
}
