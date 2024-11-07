package email

import (
	"fmt"
	"testing"

	mocksEmail "github.com/fullstackdev42/mp-emailer/mocks/email"
	"github.com/stretchr/testify/assert"
)

func TestMailpitEmailService_SendEmail(t *testing.T) {
	mockSMTP := new(mocksEmail.MockSMTPClient)
	service := &MailpitEmailService{
		host:       "localhost",
		port:       "1025",
		smtpClient: mockSMTP,
	}

	addr := fmt.Sprintf("%s:%s", service.host, service.port)
	message := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", "test@example.com", "Subject", "Body"))

	mockSMTP.On("SendMail", addr, nil, "no-reply@example.com", []string{"test@example.com"}, message).Return(nil)

	err := service.SendEmail("test@example.com", "Subject", "Body", false)

	assert.NoError(t, err)
	mockSMTP.AssertExpectations(t)
}
