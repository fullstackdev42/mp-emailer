package email

import (
	"testing"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMailgunEmailService is a mock implementation of the Mailgun email service
type MockMailgunEmailService struct {
	mock.Mock
}

func (m *MockMailgunEmailService) SendEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

// MockMailpitEmailService is a mock implementation of the Mailpit email service
type MockMailpitEmailService struct {
	mock.Mock
}

func (m *MockMailpitEmailService) SendEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.Config
		expected Service
	}{
		{
			name: "Production environment",
			config: &config.Config{
				AppEnv:        "production",
				MailgunDomain: "example.com",
				MailgunAPIKey: "key",
			},
			expected: &MailgunEmailService{},
		},
		{
			name: "Non-production environment",
			config: &config.Config{
				AppEnv:      "development",
				MailpitHost: "localhost",
				MailpitPort: "1025",
			},
			expected: &MailpitEmailService{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := New(tt.config)
			assert.IsType(t, tt.expected, service)
		})
	}
}

func TestSendEmail_Mailgun(t *testing.T) {
	mockService := new(MockMailgunEmailService)
	mockService.On("SendEmail", "test@example.com", "Subject", "Body").Return(nil)

	err := mockService.SendEmail("test@example.com", "Subject", "Body")

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestSendEmail_Mailpit(t *testing.T) {
	mockService := new(MockMailpitEmailService)
	mockService.On("SendEmail", "test@example.com", "Subject", "Body").Return(nil)

	err := mockService.SendEmail("test@example.com", "Subject", "Body")

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}
