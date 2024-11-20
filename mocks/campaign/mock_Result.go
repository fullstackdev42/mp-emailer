// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	campaign "github.com/jonesrussell/mp-emailer/campaign"
	mock "github.com/stretchr/testify/mock"
)

// MockResult is an autogenerated mock type for the Result type
type MockResult struct {
	mock.Mock
}

type MockResult_Expecter struct {
	mock *mock.Mock
}

func (_m *MockResult) EXPECT() *MockResult_Expecter {
	return &MockResult_Expecter{mock: &_m.Mock}
}

// Error provides a mock function with given fields:
func (_m *MockResult) Error() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Error")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockResult_Error_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Error'
type MockResult_Error_Call struct {
	*mock.Call
}

// Error is a helper method to define mock.On call
func (_e *MockResult_Expecter) Error() *MockResult_Error_Call {
	return &MockResult_Error_Call{Call: _e.mock.On("Error")}
}

func (_c *MockResult_Error_Call) Run(run func()) *MockResult_Error_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockResult_Error_Call) Return(_a0 error) *MockResult_Error_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockResult_Error_Call) RunAndReturn(run func() error) *MockResult_Error_Call {
	_c.Call.Return(run)
	return _c
}

// Scan provides a mock function with given fields: dest
func (_m *MockResult) Scan(dest interface{}) campaign.Result {
	ret := _m.Called(dest)

	if len(ret) == 0 {
		panic("no return value specified for Scan")
	}

	var r0 campaign.Result
	if rf, ok := ret.Get(0).(func(interface{}) campaign.Result); ok {
		r0 = rf(dest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(campaign.Result)
		}
	}

	return r0
}

// MockResult_Scan_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Scan'
type MockResult_Scan_Call struct {
	*mock.Call
}

// Scan is a helper method to define mock.On call
//   - dest interface{}
func (_e *MockResult_Expecter) Scan(dest interface{}) *MockResult_Scan_Call {
	return &MockResult_Scan_Call{Call: _e.mock.On("Scan", dest)}
}

func (_c *MockResult_Scan_Call) Run(run func(dest interface{})) *MockResult_Scan_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *MockResult_Scan_Call) Return(_a0 campaign.Result) *MockResult_Scan_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockResult_Scan_Call) RunAndReturn(run func(interface{}) campaign.Result) *MockResult_Scan_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockResult creates a new instance of MockResult. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockResult(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockResult {
	mock := &MockResult{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
