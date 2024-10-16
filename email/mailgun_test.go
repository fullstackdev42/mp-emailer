package email

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMailgunClient is a mock implementation of the Mailgun client
type MockMailgunClient struct {
	mock.Mock
}

func (m *MockMailgunClient) NewMessage(from, subject, body string, to ...string) *mailgun.Message {
	args := m.Called(from, subject, body, to)
	return args.Get(0).(*mailgun.Message)
}

func (m *MockMailgunClient) Send(ctx context.Context, message *mailgun.Message) (string, string, error) {
	args := m.Called(ctx, message)
	return args.String(0), args.String(1), args.Error(2)
}

func TestMailgunEmailService_SendEmail(t *testing.T) {
	mockClient := new(MockMailgunClient)
	service := &MailgunEmailService{
		domain: "example.com",
		apiKey: "key",
		client: mockClient,
	}

	message := &mailgun.Message{} // Create an empty message object

	mockClient.On("NewMessage", "no-reply@example.com", "Subject", "Body", []string{"test@example.com"}).Return(message)
	mockClient.On("Send", mock.Anything, message).Return("", "", nil)

	err := service.SendEmail("test@example.com", "Subject", "Body")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}
