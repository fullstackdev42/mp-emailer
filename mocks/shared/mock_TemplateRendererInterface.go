// Code generated by mockery v2.46.3. DO NOT EDIT.

package shared

import (
	io "io"

	echo "github.com/labstack/echo/v4"

	mock "github.com/stretchr/testify/mock"

	shared "github.com/fullstackdev42/mp-emailer/shared"
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

// Render provides a mock function with given fields: w, name, data, c
func (_m *MockTemplateRendererInterface) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	ret := _m.Called(w, name, data, c)

	if len(ret) == 0 {
		panic("no return value specified for Render")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Writer, string, interface{}, echo.Context) error); ok {
		r0 = rf(w, name, data, c)
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
//   - w io.Writer
//   - name string
//   - data interface{}
//   - c echo.Context
func (_e *MockTemplateRendererInterface_Expecter) Render(w interface{}, name interface{}, data interface{}, c interface{}) *MockTemplateRendererInterface_Render_Call {
	return &MockTemplateRendererInterface_Render_Call{Call: _e.mock.On("Render", w, name, data, c)}
}

func (_c *MockTemplateRendererInterface_Render_Call) Run(run func(w io.Writer, name string, data interface{}, c echo.Context)) *MockTemplateRendererInterface_Render_Call {
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

// RenderPage provides a mock function with given fields: c, templateName, pageData, errorHandler
func (_m *MockTemplateRendererInterface) RenderPage(c echo.Context, templateName string, pageData shared.Data, errorHandler shared.ErrorHandlerInterface) error {
	ret := _m.Called(c, templateName, pageData, errorHandler)

	if len(ret) == 0 {
		panic("no return value specified for RenderPage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, string, shared.Data, shared.ErrorHandlerInterface) error); ok {
		r0 = rf(c, templateName, pageData, errorHandler)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTemplateRendererInterface_RenderPage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RenderPage'
type MockTemplateRendererInterface_RenderPage_Call struct {
	*mock.Call
}

// RenderPage is a helper method to define mock.On call
//   - c echo.Context
//   - templateName string
//   - pageData shared.Data
//   - errorHandler shared.ErrorHandlerInterface
func (_e *MockTemplateRendererInterface_Expecter) RenderPage(c interface{}, templateName interface{}, pageData interface{}, errorHandler interface{}) *MockTemplateRendererInterface_RenderPage_Call {
	return &MockTemplateRendererInterface_RenderPage_Call{Call: _e.mock.On("RenderPage", c, templateName, pageData, errorHandler)}
}

func (_c *MockTemplateRendererInterface_RenderPage_Call) Run(run func(c echo.Context, templateName string, pageData shared.Data, errorHandler shared.ErrorHandlerInterface)) *MockTemplateRendererInterface_RenderPage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(echo.Context), args[1].(string), args[2].(shared.Data), args[3].(shared.ErrorHandlerInterface))
	})
	return _c
}

func (_c *MockTemplateRendererInterface_RenderPage_Call) Return(_a0 error) *MockTemplateRendererInterface_RenderPage_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTemplateRendererInterface_RenderPage_Call) RunAndReturn(run func(echo.Context, string, shared.Data, shared.ErrorHandlerInterface) error) *MockTemplateRendererInterface_RenderPage_Call {
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
