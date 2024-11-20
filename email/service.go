package email

type ServiceInterface interface {
	SendPasswordReset(email string, token string) error
}
