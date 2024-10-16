package mocks

import (
	"github.com/jonesrussell/loggo"
	"github.com/stretchr/testify/mock"
)

// MockLoggerInterface is an autogenerated mock type for the LoggerInterface type
type MockLoggerInterface struct {
	mock.Mock
}

type MockLoggerInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLoggerInterface) EXPECT() *MockLoggerInterface_Expecter {
	return &MockLoggerInterface_Expecter{mock: &_m.Mock}
}

// Debug provides a mock function with given fields: msg, args
func (_m *MockLoggerInterface) Debug(msg string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// MockLoggerInterface_Debug_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Debug'
type MockLoggerInterface_Debug_Call struct {
	*mock.Call
}

// Debug is a helper method to define mock.On call
//   - msg string
//   - args ...interface{}
func (_e *MockLoggerInterface_Expecter) Debug(msg interface{}, args ...interface{}) *MockLoggerInterface_Debug_Call {
	return &MockLoggerInterface_Debug_Call{Call: _e.mock.On("Debug",
		append([]interface{}{msg}, args...)...)}
}

func (_c *MockLoggerInterface_Debug_Call) Run(run func(msg string, args ...interface{})) *MockLoggerInterface_Debug_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockLoggerInterface_Debug_Call) Return() *MockLoggerInterface_Debug_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockLoggerInterface_Debug_Call) RunAndReturn(run func(string, ...interface{})) *MockLoggerInterface_Debug_Call {
	_c.Call.Return(run)
	return _c
}

// Error provides a mock function with given fields: msg, err, args
func (_m *MockLoggerInterface) Error(msg string, err error, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg, err)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// MockLoggerInterface_Error_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Error'
type MockLoggerInterface_Error_Call struct {
	*mock.Call
}

// Error is a helper method to define mock.On call
//   - msg string
//   - err error
//   - args ...interface{}
func (_e *MockLoggerInterface_Expecter) Error(msg interface{}, err interface{}, args ...interface{}) *MockLoggerInterface_Error_Call {
	return &MockLoggerInterface_Error_Call{Call: _e.mock.On("Error",
		append([]interface{}{msg, err}, args...)...)}
}

func (_c *MockLoggerInterface_Error_Call) Run(run func(msg string, err error, args ...interface{})) *MockLoggerInterface_Error_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), args[1].(error), variadicArgs...)
	})
	return _c
}

func (_c *MockLoggerInterface_Error_Call) Return() *MockLoggerInterface_Error_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockLoggerInterface_Error_Call) RunAndReturn(run func(string, error, ...interface{})) *MockLoggerInterface_Error_Call {
	_c.Call.Return(run)
	return _c
}

// Fatal provides a mock function with given fields: msg, err, args
func (_m *MockLoggerInterface) Fatal(msg string, err error, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg, err)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// MockLoggerInterface_Fatal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Fatal'
type MockLoggerInterface_Fatal_Call struct {
	*mock.Call
}

// Fatal is a helper method to define mock.On call
//   - msg string
//   - err error
//   - args ...interface{}
func (_e *MockLoggerInterface_Expecter) Fatal(msg interface{}, err interface{}, args ...interface{}) *MockLoggerInterface_Fatal_Call {
	return &MockLoggerInterface_Fatal_Call{Call: _e.mock.On("Fatal",
		append([]interface{}{msg, err}, args...)...)}
}

func (_c *MockLoggerInterface_Fatal_Call) Run(run func(msg string, err error, args ...interface{})) *MockLoggerInterface_Fatal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), args[1].(error), variadicArgs...)
	})
	return _c
}

func (_c *MockLoggerInterface_Fatal_Call) Return() *MockLoggerInterface_Fatal_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockLoggerInterface_Fatal_Call) RunAndReturn(run func(string, error, ...interface{})) *MockLoggerInterface_Fatal_Call {
	_c.Call.Return(run)
	return _c
}

// Info provides a mock function with given fields: msg, args
func (_m *MockLoggerInterface) Info(msg string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// MockLoggerInterface_Info_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Info'
type MockLoggerInterface_Info_Call struct {
	*mock.Call
}

// Info is a helper method to define mock.On call
//   - msg string
//   - args ...interface{}
func (_e *MockLoggerInterface_Expecter) Info(msg interface{}, args ...interface{}) *MockLoggerInterface_Info_Call {
	return &MockLoggerInterface_Info_Call{Call: _e.mock.On("Info",
		append([]interface{}{msg}, args...)...)}
}

func (_c *MockLoggerInterface_Info_Call) Run(run func(msg string, args ...interface{})) *MockLoggerInterface_Info_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockLoggerInterface_Info_Call) Return() *MockLoggerInterface_Info_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockLoggerInterface_Info_Call) RunAndReturn(run func(string, ...interface{})) *MockLoggerInterface_Info_Call {
	_c.Call.Return(run)
	return _c
}

// IsDebugEnabled provides a mock function with given fields:
func (_m *MockLoggerInterface) IsDebugEnabled() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsDebugEnabled")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockLoggerInterface_IsDebugEnabled_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsDebugEnabled'
type MockLoggerInterface_IsDebugEnabled_Call struct {
	*mock.Call
}

// IsDebugEnabled is a helper method to define mock.On call
func (_e *MockLoggerInterface_Expecter) IsDebugEnabled() *MockLoggerInterface_IsDebugEnabled_Call {
	return &MockLoggerInterface_IsDebugEnabled_Call{Call: _e.mock.On("IsDebugEnabled")}
}

func (_c *MockLoggerInterface_IsDebugEnabled_Call) Run(run func()) *MockLoggerInterface_IsDebugEnabled_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockLoggerInterface_IsDebugEnabled_Call) Return(_a0 bool) *MockLoggerInterface_IsDebugEnabled_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLoggerInterface_IsDebugEnabled_Call) RunAndReturn(run func() bool) *MockLoggerInterface_IsDebugEnabled_Call {
	_c.Call.Return(run)
	return _c
}

// Warn provides a mock function with given fields: msg, args
func (_m *MockLoggerInterface) Warn(msg string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// MockLoggerInterface_Warn_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Warn'
type MockLoggerInterface_Warn_Call struct {
	*mock.Call
}

// Warn is a helper method to define mock.On call
//   - msg string
//   - args ...interface{}
func (_e *MockLoggerInterface_Expecter) Warn(msg interface{}, args ...interface{}) *MockLoggerInterface_Warn_Call {
	return &MockLoggerInterface_Warn_Call{Call: _e.mock.On("Warn",
		append([]interface{}{msg}, args...)...)}
}

func (_c *MockLoggerInterface_Warn_Call) Run(run func(msg string, args ...interface{})) *MockLoggerInterface_Warn_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockLoggerInterface_Warn_Call) Return() *MockLoggerInterface_Warn_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockLoggerInterface_Warn_Call) RunAndReturn(run func(string, ...interface{})) *MockLoggerInterface_Warn_Call {
	_c.Call.Return(run)
	return _c
}

// WithOperation provides a mock function with given fields: operationID
func (_m *MockLoggerInterface) WithOperation(operationID string) loggo.LoggerInterface {
	ret := _m.Called(operationID)

	if len(ret) == 0 {
		panic("no return value specified for WithOperation")
	}

	var r0 loggo.LoggerInterface
	if rf, ok := ret.Get(0).(func(string) loggo.LoggerInterface); ok {
		r0 = rf(operationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(loggo.LoggerInterface)
		}
	}

	return r0
}

// MockLoggerInterface_WithOperation_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithOperation'
type MockLoggerInterface_WithOperation_Call struct {
	*mock.Call
}

// WithOperation is a helper method to define mock.On call
//   - operationID string
func (_e *MockLoggerInterface_Expecter) WithOperation(operationID interface{}) *MockLoggerInterface_WithOperation_Call {
	return &MockLoggerInterface_WithOperation_Call{Call: _e.mock.On("WithOperation", operationID)}
}

func (_c *MockLoggerInterface_WithOperation_Call) Run(run func(operationID string)) *MockLoggerInterface_WithOperation_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockLoggerInterface_WithOperation_Call) Return(_a0 loggo.LoggerInterface) *MockLoggerInterface_WithOperation_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLoggerInterface_WithOperation_Call) RunAndReturn(run func(string) loggo.LoggerInterface) *MockLoggerInterface_WithOperation_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockLoggerInterface creates a new instance of MockLoggerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockLoggerInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockLoggerInterface {
	mock := &MockLoggerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}