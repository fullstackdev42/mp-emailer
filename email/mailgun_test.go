package email

import (
	"testing"

	mocksEmail "github.com/jonesrussell/mp-emailer/mocks/email"
	mocksLogger "github.com/jonesrussell/mp-emailer/mocks/logger"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMailgunEmailService_SendEmail(t *testing.T) {
	mockMailgun := new(mocksEmail.MockMailgunClient)
	mockLogger := mocksLogger.NewMockInterface(t)

	// Set up logger expectations
	mockLogger.On("Debug", "Email sent successfully", "messageId", "").Return()

	service := &MailgunEmailService{
		domain: "example.com",
		apiKey: "key",
		client: mockMailgun,
		Logger: mockLogger,
	}

	message := &mailgun.Message{}

	mockMailgun.On("NewMessage", "no-reply@example.com", "Subject", "Body", "test@example.com").Return(message)
	mockMailgun.On("Send", mock.Anything, message).Return("", "", nil)

	err := service.SendEmail("test@example.com", "Subject", "Body", false)

	assert.NoError(t, err)
	mockMailgun.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
