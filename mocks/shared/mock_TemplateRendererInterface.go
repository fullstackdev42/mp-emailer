// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	io "io"

	echo "github.com/labstack/echo/v4"

	mock "github.com/stretchr/testify/mock"
)

// MockTemplateRendererInterface is an autogenerated mock type for the TemplateRendererInterface type
type MockTemplateRendererInterface struct {
	mock.Mock
}

type MockTemplateRendererInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTemplateRendererInterface) EXPECT() *MockTemplateRendererInterface_Expecter {
	return &MockTemplateRendererInterface_Expecter{mock: &_m.Mock}
}

// Render provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *MockTemplateRendererInterface) Render(_a0 io.Writer, _a1 string, _a2 interface{}, _a3 echo.Context) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	if len(ret) == 0 {
		panic("no return value specified for Render")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Writer, string, interface{}, echo.Context) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTemplateRendererInterface_Render_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Render'
type MockTemplateRendererInterface_Render_Call struct {
	*mock.Call
}

// Render is a helper method to define mock.On call
//   - _a0 io.Writer
//   - _a1 string
//   - _a2 interface{}
//   - _a3 echo.Context
func (_e *MockTemplateRendererInterface_Expecter) Render(_a0 interface{}, _a1 interface{}, _a2 interface{}, _a3 interface{}) *MockTemplateRendererInterface_Render_Call {
	return &MockTemplateRendererInterface_Render_Call{Call: _e.mock.On("Render", _a0, _a1, _a2, _a3)}
}

func (_c *MockTemplateRendererInterface_Render_Call) Run(run func(_a0 io.Writer, _a1 string, _a2 interface{}, _a3 echo.Context)) *MockTemplateRendererInterface_Render_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(io.Writer), args[1].(string), args[2].(interface{}), args[3].(echo.Context))
	})
	return _c
}

func (_c *MockTemplateRendererInterface_Render_Call) Return(_a0 error) *MockTemplateRendererInterface_Render_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTemplateRendererInterface_Render_Call) RunAndReturn(run func(io.Writer, string, interface{}, echo.Context) error) *MockTemplateRendererInterface_Render_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTemplateRendererInterface creates a new instance of MockTemplateRendererInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTemplateRendererInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTemplateRendererInterface {
	mock := &MockTemplateRendererInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
