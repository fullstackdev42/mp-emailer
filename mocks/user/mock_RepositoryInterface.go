// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	user "github.com/jonesrussell/mp-emailer/user"
	mock "github.com/stretchr/testify/mock"
)

// MockRepositoryInterface is an autogenerated mock type for the RepositoryInterface type
type MockRepositoryInterface struct {
	mock.Mock
}

type MockRepositoryInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRepositoryInterface) EXPECT() *MockRepositoryInterface_Expecter {
	return &MockRepositoryInterface_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: _a0
func (_m *MockRepositoryInterface) Create(_a0 *user.User) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*user.User) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRepositoryInterface_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockRepositoryInterface_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - _a0 *user.User
func (_e *MockRepositoryInterface_Expecter) Create(_a0 interface{}) *MockRepositoryInterface_Create_Call {
	return &MockRepositoryInterface_Create_Call{Call: _e.mock.On("Create", _a0)}
}

func (_c *MockRepositoryInterface_Create_Call) Run(run func(_a0 *user.User)) *MockRepositoryInterface_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*user.User))
	})
	return _c
}

func (_c *MockRepositoryInterface_Create_Call) Return(_a0 error) *MockRepositoryInterface_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRepositoryInterface_Create_Call) RunAndReturn(run func(*user.User) error) *MockRepositoryInterface_Create_Call {
	_c.Call.Return(run)
	return _c
}

// FindByEmail provides a mock function with given fields: email
func (_m *MockRepositoryInterface) FindByEmail(email string) (*user.User, error) {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for FindByEmail")
	}

	var r0 *user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*user.User, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) *user.User); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepositoryInterface_FindByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByEmail'
type MockRepositoryInterface_FindByEmail_Call struct {
	*mock.Call
}

// FindByEmail is a helper method to define mock.On call
//   - email string
func (_e *MockRepositoryInterface_Expecter) FindByEmail(email interface{}) *MockRepositoryInterface_FindByEmail_Call {
	return &MockRepositoryInterface_FindByEmail_Call{Call: _e.mock.On("FindByEmail", email)}
}

func (_c *MockRepositoryInterface_FindByEmail_Call) Run(run func(email string)) *MockRepositoryInterface_FindByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockRepositoryInterface_FindByEmail_Call) Return(_a0 *user.User, _a1 error) *MockRepositoryInterface_FindByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepositoryInterface_FindByEmail_Call) RunAndReturn(run func(string) (*user.User, error)) *MockRepositoryInterface_FindByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// FindByUsername provides a mock function with given fields: username
func (_m *MockRepositoryInterface) FindByUsername(username string) (*user.User, error) {
	ret := _m.Called(username)

	if len(ret) == 0 {
		panic("no return value specified for FindByUsername")
	}

	var r0 *user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*user.User, error)); ok {
		return rf(username)
	}
	if rf, ok := ret.Get(0).(func(string) *user.User); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepositoryInterface_FindByUsername_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByUsername'
type MockRepositoryInterface_FindByUsername_Call struct {
	*mock.Call
}

// FindByUsername is a helper method to define mock.On call
//   - username string
func (_e *MockRepositoryInterface_Expecter) FindByUsername(username interface{}) *MockRepositoryInterface_FindByUsername_Call {
	return &MockRepositoryInterface_FindByUsername_Call{Call: _e.mock.On("FindByUsername", username)}
}

func (_c *MockRepositoryInterface_FindByUsername_Call) Run(run func(username string)) *MockRepositoryInterface_FindByUsername_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockRepositoryInterface_FindByUsername_Call) Return(_a0 *user.User, _a1 error) *MockRepositoryInterface_FindByUsername_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepositoryInterface_FindByUsername_Call) RunAndReturn(run func(string) (*user.User, error)) *MockRepositoryInterface_FindByUsername_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRepositoryInterface creates a new instance of MockRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRepositoryInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRepositoryInterface {
	mock := &MockRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
