package email

import (
	"testing"

	"github.com/fullstackdev42/mp-emailer/mocks"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMailgunEmailService_SendEmail(t *testing.T) {
	mockClient := new(MockMailgunClient)
	service := &MailgunEmailService{
		domain: "example.com",
		apiKey: "key",
		client: mockClient,
		logger: mocks.NewMockLoggerInterface(t),
	}

	message := &mailgun.Message{}

	// Update this line to pass the 'to' string directly
	mockClient.On("NewMessage", "no-reply@example.com", "Subject", "Body", "test@example.com").Return(message)
	mockClient.On("Send", mock.Anything, message).Return("", "", nil)

	err := service.SendEmail("test@example.com", "Subject", "Body")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}
