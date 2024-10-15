package mocks

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/mock"
)

// SessionStore is an interface that matches the methods of sessions.Store
type SessionStore interface {
	Get(r *http.Request, name string) (*sessions.Session, error)
	New(r *http.Request, name string) (*sessions.Session, error)
	Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error
}

// MockSessionStore is a mock implementation of the SessionStore interface
type MockSessionStore struct {
	mock.Mock
}

// Get mocks the Get method
func (m *MockSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	args := m.Called(r, name)
	return args.Get(0).(*sessions.Session), args.Error(1)
}

// New mocks the New method
func (m *MockSessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	args := m.Called(r, name)
	return args.Get(0).(*sessions.Session), args.Error(1)
}

// Save mocks the Save method
func (m *MockSessionStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	args := m.Called(r, w, session)
	return args.Error(0)
}
