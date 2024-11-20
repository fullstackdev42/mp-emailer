package email

import (
	"testing"

	mocksEmail "github.com/jonesrussell/mp-emailer/mocks/email"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMailpitEmailService_SendEmail(t *testing.T) {
	// Create mock SMTP client
	mockSMTP := new(mocksEmail.MockSMTPClient)

	// Set up SMTP client expectations
	mockSMTP.On("SendMail",
		"localhost:1025",
		mock.Anything,
		"test@example.com",
		[]string{"recipient@example.com"},
		mock.AnythingOfType("[]uint8"),
	).Return(nil)

	// Create service with mock SMTP client
	service := NewMailpitEmailService(
		"localhost",
		"1025",
		mockSMTP,
		"test@example.com",
	)

	// Test sending HTML email
	err := service.SendEmail("recipient@example.com", "Test Subject", "Test Body", true)

	assert.NoError(t, err)
	mockSMTP.AssertExpectations(t)
}
