package email

import (
	"testing"

	"github.com/fullstackdev42/mp-emailer/mocks"
	mocksEmail "github.com/fullstackdev42/mp-emailer/mocks/email"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMailgunEmailService_SendEmail(t *testing.T) {
	MockMailgunClient := new(mocksEmail.MockMailgunClient)

	service := &MailgunEmailService{
		domain: "example.com",
		apiKey: "key",
		client: MockMailgunClient,
		logger: mocks.NewMockLoggerInterface(t),
	}

	message := &mailgun.Message{}

	// Update this line to pass the 'to' string directly
	MockMailgunClient.On("NewMessage", "no-reply@example.com", "Subject", "Body", "test@example.com").Return(message)
	MockMailgunClient.On("Send", mock.Anything, message).Return("", "", nil)

	err := service.SendEmail("test@example.com", "Subject", "Body")

	assert.NoError(t, err)
	MockMailgunClient.AssertExpectations(t)
}
