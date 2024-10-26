// Code generated by mockery v2.46.3. DO NOT EDIT.

package email

import (
	smtp "net/smtp"

	mock "github.com/stretchr/testify/mock"
)

// MockSMTPClient is an autogenerated mock type for the SMTPClient type
type MockSMTPClient struct {
	mock.Mock
}

type MockSMTPClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSMTPClient) EXPECT() *MockSMTPClient_Expecter {
	return &MockSMTPClient_Expecter{mock: &_m.Mock}
}

// SendMail provides a mock function with given fields: addr, a, from, to, msg
func (_m *MockSMTPClient) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	ret := _m.Called(addr, a, from, to, msg)

	if len(ret) == 0 {
		panic("no return value specified for SendMail")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, smtp.Auth, string, []string, []byte) error); ok {
		r0 = rf(addr, a, from, to, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSMTPClient_SendMail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendMail'
type MockSMTPClient_SendMail_Call struct {
	*mock.Call
}

// SendMail is a helper method to define mock.On call
//   - addr string
//   - a smtp.Auth
//   - from string
//   - to []string
//   - msg []byte
func (_e *MockSMTPClient_Expecter) SendMail(addr interface{}, a interface{}, from interface{}, to interface{}, msg interface{}) *MockSMTPClient_SendMail_Call {
	return &MockSMTPClient_SendMail_Call{Call: _e.mock.On("SendMail", addr, a, from, to, msg)}
}

func (_c *MockSMTPClient_SendMail_Call) Run(run func(addr string, a smtp.Auth, from string, to []string, msg []byte)) *MockSMTPClient_SendMail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(smtp.Auth), args[2].(string), args[3].([]string), args[4].([]byte))
	})
	return _c
}

func (_c *MockSMTPClient_SendMail_Call) Return(_a0 error) *MockSMTPClient_SendMail_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSMTPClient_SendMail_Call) RunAndReturn(run func(string, smtp.Auth, string, []string, []byte) error) *MockSMTPClient_SendMail_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSMTPClient creates a new instance of MockSMTPClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSMTPClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSMTPClient {
	mock := &MockSMTPClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}