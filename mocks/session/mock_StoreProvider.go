// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	http "net/http"

	echo "github.com/labstack/echo/v4"

	mock "github.com/stretchr/testify/mock"

	session "github.com/jonesrussell/mp-emailer/session"
)

// MockStoreProvider is an autogenerated mock type for the StoreProvider type
type MockStoreProvider struct {
	mock.Mock
}

type MockStoreProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *MockStoreProvider) EXPECT() *MockStoreProvider_Expecter {
	return &MockStoreProvider_Expecter{mock: &_m.Mock}
}

// GetStore provides a mock function with given fields: r
func (_m *MockStoreProvider) GetStore(r *http.Request) session.Store {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for GetStore")
	}

	var r0 session.Store
	if rf, ok := ret.Get(0).(func(*http.Request) session.Store); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(session.Store)
		}
	}

	return r0
}

// MockStoreProvider_GetStore_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetStore'
type MockStoreProvider_GetStore_Call struct {
	*mock.Call
}

// GetStore is a helper method to define mock.On call
//   - r *http.Request
func (_e *MockStoreProvider_Expecter) GetStore(r interface{}) *MockStoreProvider_GetStore_Call {
	return &MockStoreProvider_GetStore_Call{Call: _e.mock.On("GetStore", r)}
}

func (_c *MockStoreProvider_GetStore_Call) Run(run func(r *http.Request)) *MockStoreProvider_GetStore_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*http.Request))
	})
	return _c
}

func (_c *MockStoreProvider_GetStore_Call) Return(_a0 session.Store) *MockStoreProvider_GetStore_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockStoreProvider_GetStore_Call) RunAndReturn(run func(*http.Request) session.Store) *MockStoreProvider_GetStore_Call {
	_c.Call.Return(run)
	return _c
}

// SetStore provides a mock function with given fields: c, store
func (_m *MockStoreProvider) SetStore(c echo.Context, store session.Store) {
	_m.Called(c, store)
}

// MockStoreProvider_SetStore_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetStore'
type MockStoreProvider_SetStore_Call struct {
	*mock.Call
}

// SetStore is a helper method to define mock.On call
//   - c echo.Context
//   - store session.Store
func (_e *MockStoreProvider_Expecter) SetStore(c interface{}, store interface{}) *MockStoreProvider_SetStore_Call {
	return &MockStoreProvider_SetStore_Call{Call: _e.mock.On("SetStore", c, store)}
}

func (_c *MockStoreProvider_SetStore_Call) Run(run func(c echo.Context, store session.Store)) *MockStoreProvider_SetStore_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(echo.Context), args[1].(session.Store))
	})
	return _c
}

func (_c *MockStoreProvider_SetStore_Call) Return() *MockStoreProvider_SetStore_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockStoreProvider_SetStore_Call) RunAndReturn(run func(echo.Context, session.Store)) *MockStoreProvider_SetStore_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockStoreProvider creates a new instance of MockStoreProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockStoreProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockStoreProvider {
	mock := &MockStoreProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}