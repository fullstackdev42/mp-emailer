// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	campaign "github.com/jonesrussell/mp-emailer/campaign"

	mock "github.com/stretchr/testify/mock"
)

// MockServiceInterface is an autogenerated mock type for the ServiceInterface type
type MockServiceInterface struct {
	mock.Mock
}

type MockServiceInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockServiceInterface) EXPECT() *MockServiceInterface_Expecter {
	return &MockServiceInterface_Expecter{mock: &_m.Mock}
}

// ComposeEmail provides a mock function with given fields: ctx, params
func (_m *MockServiceInterface) ComposeEmail(ctx context.Context, params campaign.ComposeEmailParams) (string, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for ComposeEmail")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, campaign.ComposeEmailParams) (string, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, campaign.ComposeEmailParams) string); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, campaign.ComposeEmailParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockServiceInterface_ComposeEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ComposeEmail'
type MockServiceInterface_ComposeEmail_Call struct {
	*mock.Call
}

// ComposeEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - params campaign.ComposeEmailParams
func (_e *MockServiceInterface_Expecter) ComposeEmail(ctx interface{}, params interface{}) *MockServiceInterface_ComposeEmail_Call {
	return &MockServiceInterface_ComposeEmail_Call{Call: _e.mock.On("ComposeEmail", ctx, params)}
}

func (_c *MockServiceInterface_ComposeEmail_Call) Run(run func(ctx context.Context, params campaign.ComposeEmailParams)) *MockServiceInterface_ComposeEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(campaign.ComposeEmailParams))
	})
	return _c
}

func (_c *MockServiceInterface_ComposeEmail_Call) Return(_a0 string, _a1 error) *MockServiceInterface_ComposeEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockServiceInterface_ComposeEmail_Call) RunAndReturn(run func(context.Context, campaign.ComposeEmailParams) (string, error)) *MockServiceInterface_ComposeEmail_Call {
	_c.Call.Return(run)
	return _c
}

// CreateCampaign provides a mock function with given fields: ctx, dto
func (_m *MockServiceInterface) CreateCampaign(ctx context.Context, dto *campaign.CreateCampaignDTO) (*campaign.Campaign, error) {
	ret := _m.Called(ctx, dto)

	if len(ret) == 0 {
		panic("no return value specified for CreateCampaign")
	}

	var r0 *campaign.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *campaign.CreateCampaignDTO) (*campaign.Campaign, error)); ok {
		return rf(ctx, dto)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *campaign.CreateCampaignDTO) *campaign.Campaign); ok {
		r0 = rf(ctx, dto)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*campaign.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *campaign.CreateCampaignDTO) error); ok {
		r1 = rf(ctx, dto)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockServiceInterface_CreateCampaign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateCampaign'
type MockServiceInterface_CreateCampaign_Call struct {
	*mock.Call
}

// CreateCampaign is a helper method to define mock.On call
//   - ctx context.Context
//   - dto *campaign.CreateCampaignDTO
func (_e *MockServiceInterface_Expecter) CreateCampaign(ctx interface{}, dto interface{}) *MockServiceInterface_CreateCampaign_Call {
	return &MockServiceInterface_CreateCampaign_Call{Call: _e.mock.On("CreateCampaign", ctx, dto)}
}

func (_c *MockServiceInterface_CreateCampaign_Call) Run(run func(ctx context.Context, dto *campaign.CreateCampaignDTO)) *MockServiceInterface_CreateCampaign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*campaign.CreateCampaignDTO))
	})
	return _c
}

func (_c *MockServiceInterface_CreateCampaign_Call) Return(_a0 *campaign.Campaign, _a1 error) *MockServiceInterface_CreateCampaign_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockServiceInterface_CreateCampaign_Call) RunAndReturn(run func(context.Context, *campaign.CreateCampaignDTO) (*campaign.Campaign, error)) *MockServiceInterface_CreateCampaign_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteCampaign provides a mock function with given fields: ctx, params
func (_m *MockServiceInterface) DeleteCampaign(ctx context.Context, params campaign.DeleteCampaignDTO) error {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for DeleteCampaign")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, campaign.DeleteCampaignDTO) error); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockServiceInterface_DeleteCampaign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteCampaign'
type MockServiceInterface_DeleteCampaign_Call struct {
	*mock.Call
}

// DeleteCampaign is a helper method to define mock.On call
//   - ctx context.Context
//   - params campaign.DeleteCampaignDTO
func (_e *MockServiceInterface_Expecter) DeleteCampaign(ctx interface{}, params interface{}) *MockServiceInterface_DeleteCampaign_Call {
	return &MockServiceInterface_DeleteCampaign_Call{Call: _e.mock.On("DeleteCampaign", ctx, params)}
}

func (_c *MockServiceInterface_DeleteCampaign_Call) Run(run func(ctx context.Context, params campaign.DeleteCampaignDTO)) *MockServiceInterface_DeleteCampaign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(campaign.DeleteCampaignDTO))
	})
	return _c
}

func (_c *MockServiceInterface_DeleteCampaign_Call) Return(_a0 error) *MockServiceInterface_DeleteCampaign_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockServiceInterface_DeleteCampaign_Call) RunAndReturn(run func(context.Context, campaign.DeleteCampaignDTO) error) *MockServiceInterface_DeleteCampaign_Call {
	_c.Call.Return(run)
	return _c
}

// FetchCampaign provides a mock function with given fields: ctx, params
func (_m *MockServiceInterface) FetchCampaign(ctx context.Context, params campaign.GetCampaignParams) (*campaign.Campaign, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for FetchCampaign")
	}

	var r0 *campaign.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, campaign.GetCampaignParams) (*campaign.Campaign, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, campaign.GetCampaignParams) *campaign.Campaign); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*campaign.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, campaign.GetCampaignParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockServiceInterface_FetchCampaign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchCampaign'
type MockServiceInterface_FetchCampaign_Call struct {
	*mock.Call
}

// FetchCampaign is a helper method to define mock.On call
//   - ctx context.Context
//   - params campaign.GetCampaignParams
func (_e *MockServiceInterface_Expecter) FetchCampaign(ctx interface{}, params interface{}) *MockServiceInterface_FetchCampaign_Call {
	return &MockServiceInterface_FetchCampaign_Call{Call: _e.mock.On("FetchCampaign", ctx, params)}
}

func (_c *MockServiceInterface_FetchCampaign_Call) Run(run func(ctx context.Context, params campaign.GetCampaignParams)) *MockServiceInterface_FetchCampaign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(campaign.GetCampaignParams))
	})
	return _c
}

func (_c *MockServiceInterface_FetchCampaign_Call) Return(_a0 *campaign.Campaign, _a1 error) *MockServiceInterface_FetchCampaign_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockServiceInterface_FetchCampaign_Call) RunAndReturn(run func(context.Context, campaign.GetCampaignParams) (*campaign.Campaign, error)) *MockServiceInterface_FetchCampaign_Call {
	_c.Call.Return(run)
	return _c
}

// GetCampaignByID provides a mock function with given fields: ctx, params
func (_m *MockServiceInterface) GetCampaignByID(ctx context.Context, params campaign.GetCampaignParams) (*campaign.Campaign, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for GetCampaignByID")
	}

	var r0 *campaign.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, campaign.GetCampaignParams) (*campaign.Campaign, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, campaign.GetCampaignParams) *campaign.Campaign); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*campaign.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, campaign.GetCampaignParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockServiceInterface_GetCampaignByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCampaignByID'
type MockServiceInterface_GetCampaignByID_Call struct {
	*mock.Call
}

// GetCampaignByID is a helper method to define mock.On call
//   - ctx context.Context
//   - params campaign.GetCampaignParams
func (_e *MockServiceInterface_Expecter) GetCampaignByID(ctx interface{}, params interface{}) *MockServiceInterface_GetCampaignByID_Call {
	return &MockServiceInterface_GetCampaignByID_Call{Call: _e.mock.On("GetCampaignByID", ctx, params)}
}

func (_c *MockServiceInterface_GetCampaignByID_Call) Run(run func(ctx context.Context, params campaign.GetCampaignParams)) *MockServiceInterface_GetCampaignByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(campaign.GetCampaignParams))
	})
	return _c
}

func (_c *MockServiceInterface_GetCampaignByID_Call) Return(_a0 *campaign.Campaign, _a1 error) *MockServiceInterface_GetCampaignByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockServiceInterface_GetCampaignByID_Call) RunAndReturn(run func(context.Context, campaign.GetCampaignParams) (*campaign.Campaign, error)) *MockServiceInterface_GetCampaignByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetCampaigns provides a mock function with given fields: ctx
func (_m *MockServiceInterface) GetCampaigns(ctx context.Context) ([]campaign.Campaign, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetCampaigns")
	}

	var r0 []campaign.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]campaign.Campaign, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []campaign.Campaign); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]campaign.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockServiceInterface_GetCampaigns_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCampaigns'
type MockServiceInterface_GetCampaigns_Call struct {
	*mock.Call
}

// GetCampaigns is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockServiceInterface_Expecter) GetCampaigns(ctx interface{}) *MockServiceInterface_GetCampaigns_Call {
	return &MockServiceInterface_GetCampaigns_Call{Call: _e.mock.On("GetCampaigns", ctx)}
}

func (_c *MockServiceInterface_GetCampaigns_Call) Run(run func(ctx context.Context)) *MockServiceInterface_GetCampaigns_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockServiceInterface_GetCampaigns_Call) Return(_a0 []campaign.Campaign, _a1 error) *MockServiceInterface_GetCampaigns_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockServiceInterface_GetCampaigns_Call) RunAndReturn(run func(context.Context) ([]campaign.Campaign, error)) *MockServiceInterface_GetCampaigns_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateCampaign provides a mock function with given fields: ctx, dto
func (_m *MockServiceInterface) UpdateCampaign(ctx context.Context, dto *campaign.UpdateCampaignDTO) error {
	ret := _m.Called(ctx, dto)

	if len(ret) == 0 {
		panic("no return value specified for UpdateCampaign")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *campaign.UpdateCampaignDTO) error); ok {
		r0 = rf(ctx, dto)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockServiceInterface_UpdateCampaign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateCampaign'
type MockServiceInterface_UpdateCampaign_Call struct {
	*mock.Call
}

// UpdateCampaign is a helper method to define mock.On call
//   - ctx context.Context
//   - dto *campaign.UpdateCampaignDTO
func (_e *MockServiceInterface_Expecter) UpdateCampaign(ctx interface{}, dto interface{}) *MockServiceInterface_UpdateCampaign_Call {
	return &MockServiceInterface_UpdateCampaign_Call{Call: _e.mock.On("UpdateCampaign", ctx, dto)}
}

func (_c *MockServiceInterface_UpdateCampaign_Call) Run(run func(ctx context.Context, dto *campaign.UpdateCampaignDTO)) *MockServiceInterface_UpdateCampaign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*campaign.UpdateCampaignDTO))
	})
	return _c
}

func (_c *MockServiceInterface_UpdateCampaign_Call) Return(_a0 error) *MockServiceInterface_UpdateCampaign_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockServiceInterface_UpdateCampaign_Call) RunAndReturn(run func(context.Context, *campaign.UpdateCampaignDTO) error) *MockServiceInterface_UpdateCampaign_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockServiceInterface creates a new instance of MockServiceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockServiceInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockServiceInterface {
	mock := &MockServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
