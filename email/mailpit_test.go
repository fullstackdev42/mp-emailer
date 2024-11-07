package email

import (
	"testing"

	"github.com/fullstackdev42/mp-emailer/mocks"
	mocksEmail "github.com/fullstackdev42/mp-emailer/mocks/email"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMailpitEmailService_SendEmail(t *testing.T) {
	// Create mocks
	mockLogger := new(mocks.MockLoggerInterface)
	mockSMTP := new(mocksEmail.MockSMTPClient)

	// Set up logger expectations
	mockLogger.On("Debug", "HTML Body content", "body", "Test Body").Return()
	mockLogger.On("Debug", "Full message", "message", mock.AnythingOfType("string")).Return()

	// Set up SMTP client expectations
	mockSMTP.On("SendMail",
		"localhost:1025",
		mock.Anything,
		"test@example.com",
		[]string{"recipient@example.com"},
		mock.AnythingOfType("[]uint8"),
	).Return(nil)

	// Create service with mocks
	service := NewMailpitEmailService(
		"localhost",
		"1025",
		mockSMTP,
		mockLogger,
		"test@example.com",
	)

	// Test sending HTML email
	err := service.SendEmail("recipient@example.com", "Test Subject", "Test Body", true)

	assert.NoError(t, err)
	mockLogger.AssertExpectations(t)
	mockSMTP.AssertExpectations(t)
}
