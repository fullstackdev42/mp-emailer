package database

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fullstackdev42/mp-emailer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// TestNewDB checks the functionality of the NewDB function when connecting to the database.
func TestNewDB(t *testing.T) {
	// Create a mock logger
	mockLogger := mocks.NewMockLoggerInterface(t)

	// Set up expectations for the mock logger
	mockLogger.EXPECT().Debug(mock.Anything).Return()
	mockLogger.EXPECT().Error("error connecting to database", mock.AnythingOfType("*mysql.MySQLError")).Return()

	// Define the Data Source Name (DSN) for testing
	testDSN := "user:password@tcp(localhost:3306)/testdb"

	// Call the NewDB function with the test DSN and mock logger
	testDB, err := NewDB(testDSN, mockLogger)

	// Assert that an error is returned
	assert.Error(t, err)

	// Assert that the testDB is nil when there is an error
	assert.Nil(t, testDB)

	// Assert that the error message contains the expected text
	assert.Contains(t, err.Error(), "error connecting to database")
}

func TestUserExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	testDB := &DB{SQL: db}

	// Set up expectations for the mock query and its result
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE username = \\? OR email = \\?").
		WithArgs("testuser", "test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Call the UserExists function and assert the result
	exists, err := testDB.UserExists("testuser", "test@example.com")
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	testDB := &DB{SQL: db}

	// Set up expectations for the mock exec and its result
	mock.ExpectExec("INSERT INTO users \\(username, email, password_hash\\) VALUES \\(\\?, \\?, \\?\\)").
		WithArgs("testuser", "test@example.com", "hashedpassword").
		WillReturnResult(sqlmock.NewResult(1, 1))
	// Call the CreateUser function and assert the result
	err = testDB.CreateUser("123e4567-e89b-12d3-a456-426614174000", "testuser", "test@example.com", "hashedpassword")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVerifyUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	testDB := &DB{SQL: db}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	// Test with correct password
	mock.ExpectQuery("SELECT id, password_hash FROM users WHERE username = \\?").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash"}).AddRow("user123", string(hashedPassword)))
	userID, err := testDB.VerifyUser("testuser", "correctpassword")
	assert.NoError(t, err)
	assert.Equal(t, "user123", userID)

	// Test with incorrect password
	mock.ExpectQuery("SELECT id, password_hash FROM users WHERE username = \\?").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash"}).AddRow("user123", string(hashedPassword)))
	_, err = testDB.VerifyUser("testuser", "wrongpassword")
	assert.Error(t, err)
	assert.Equal(t, "invalid username or password", err.Error())

	// Test with non-existent user
	mock.ExpectQuery("SELECT id, password_hash FROM users WHERE username = \\?").
		WithArgs("nonexistentuser").
		WillReturnError(sql.ErrNoRows)
	_, err = testDB.VerifyUser("nonexistentuser", "anypassword")
	assert.Error(t, err)
	assert.Equal(t, "invalid username or password", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}
