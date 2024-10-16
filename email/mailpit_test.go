package email

import (
	"fmt"
	"net/smtp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSMTPClient is a mock implementation of the SMTP client
type MockSMTPClient struct {
	mock.Mock
}

func (m *MockSMTPClient) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	args := m.Called(addr, a, from, to, msg)
	return args.Error(0)
}

func TestMailpitEmailService_SendEmail(t *testing.T) {
	mockSMTP := new(MockSMTPClient)
	service := &MailpitEmailService{
		host:       "localhost",
		port:       "1025",
		smtpClient: mockSMTP,
	}

	addr := fmt.Sprintf("%s:%s", service.host, service.port)
	message := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", "test@example.com", "Subject", "Body"))

	mockSMTP.On("SendMail", addr, nil, "no-reply@example.com", []string{"test@example.com"}, message).Return(nil)

	err := service.SendEmail("test@example.com", "Subject", "Body")

	assert.NoError(t, err)
	mockSMTP.AssertExpectations(t)
}
