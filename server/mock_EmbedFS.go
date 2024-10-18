package server

import (
	"io/fs"

	"github.com/stretchr/testify/mock"
)

// MockEmbedFS is a mock implementation of the embed.FS interface
type MockEmbedFS struct {
	mock.Mock
}

func NewMockEmbedFS() *MockEmbedFS {
	return &MockEmbedFS{}
}

func (m *MockEmbedFS) ReadFile(name string) ([]byte, error) {
	args := m.Called(name)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockEmbedFS) Open(name string) (fs.File, error) {
	args := m.Called(name)
	return args.Get(0).(fs.File), args.Error(1)
}

func (m *MockEmbedFS) ReadDir(name string) ([]fs.DirEntry, error) {
	args := m.Called(name)
	if args.Get(0) != nil {
		return args.Get(0).([]fs.DirEntry), args.Error(1)
	}
	return nil, args.Error(1)
}
