// Code generated by mockery v2.46.3. DO NOT EDIT.

package campaign

import (
	campaign "github.com/fullstackdev42/mp-emailer/campaign"
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

// ComposeEmail provides a mock function with given fields: mp, _a1, userData
func (_m *MockServiceInterface) ComposeEmail(mp campaign.Representative, _a1 *campaign.Campaign, userData map[string]string) string {
	ret := _m.Called(mp, _a1, userData)

	if len(ret) == 0 {
		panic("no return value specified for ComposeEmail")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(campaign.Representative, *campaign.Campaign, map[string]string) string); ok {
		r0 = rf(mp, _a1, userData)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockServiceInterface_ComposeEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ComposeEmail'
type MockServiceInterface_ComposeEmail_Call struct {
	*mock.Call
}

// ComposeEmail is a helper method to define mock.On call
//   - mp campaign.Representative
//   - _a1 *campaign.Campaign
//   - userData map[string]string
func (_e *MockServiceInterface_Expecter) ComposeEmail(mp interface{}, _a1 interface{}, userData interface{}) *MockServiceInterface_ComposeEmail_Call {
	return &MockServiceInterface_ComposeEmail_Call{Call: _e.mock.On("ComposeEmail", mp, _a1, userData)}
}

func (_c *MockServiceInterface_ComposeEmail_Call) Run(run func(mp campaign.Representative, _a1 *campaign.Campaign, userData map[string]string)) *MockServiceInterface_ComposeEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(campaign.Representative), args[1].(*campaign.Campaign), args[2].(map[string]string))
	})
	return _c
}

func (_c *MockServiceInterface_ComposeEmail_Call) Return(_a0 string) *MockServiceInterface_ComposeEmail_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockServiceInterface_ComposeEmail_Call) RunAndReturn(run func(campaign.Representative, *campaign.Campaign, map[string]string) string) *MockServiceInterface_ComposeEmail_Call {
	_c.Call.Return(run)
	return _c
}

// CreateCampaign provides a mock function with given fields: _a0
func (_m *MockServiceInterface) CreateCampaign(_a0 *campaign.Campaign) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for CreateCampaign")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*campaign.Campaign) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockServiceInterface_CreateCampaign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateCampaign'
type MockServiceInterface_CreateCampaign_Call struct {
	*mock.Call
}

// CreateCampaign is a helper method to define mock.On call
//   - _a0 *campaign.Campaign
func (_e *MockServiceInterface_Expecter) CreateCampaign(_a0 interface{}) *MockServiceInterface_CreateCampaign_Call {
	return &MockServiceInterface_CreateCampaign_Call{Call: _e.mock.On("CreateCampaign", _a0)}
}

func (_c *MockServiceInterface_CreateCampaign_Call) Run(run func(_a0 *campaign.Campaign)) *MockServiceInterface_CreateCampaign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*campaign.Campaign))
	})
	return _c
}

func (_c *MockServiceInterface_CreateCampaign_Call) Return(_a0 error) *MockServiceInterface_CreateCampaign_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockServiceInterface_CreateCampaign_Call) RunAndReturn(run func(*campaign.Campaign) error) *MockServiceInterface_CreateCampaign_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteCampaign provides a mock function with given fields: id
func (_m *MockServiceInterface) DeleteCampaign(id int) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteCampaign")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(id)
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
//   - id int
func (_e *MockServiceInterface_Expecter) DeleteCampaign(id interface{}) *MockServiceInterface_DeleteCampaign_Call {
	return &MockServiceInterface_DeleteCampaign_Call{Call: _e.mock.On("DeleteCampaign", id)}
}

func (_c *MockServiceInterface_DeleteCampaign_Call) Run(run func(id int)) *MockServiceInterface_DeleteCampaign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockServiceInterface_DeleteCampaign_Call) Return(_a0 error) *MockServiceInterface_DeleteCampaign_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockServiceInterface_DeleteCampaign_Call) RunAndReturn(run func(int) error) *MockServiceInterface_DeleteCampaign_Call {
	_c.Call.Return(run)
	return _c
}

// FetchCampaign provides a mock function with given fields: id
func (_m *MockServiceInterface) FetchCampaign(id int) (*campaign.Campaign, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for FetchCampaign")
	}

	var r0 *campaign.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*campaign.Campaign, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *campaign.Campaign); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*campaign.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
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
//   - id int
func (_e *MockServiceInterface_Expecter) FetchCampaign(id interface{}) *MockServiceInterface_FetchCampaign_Call {
	return &MockServiceInterface_FetchCampaign_Call{Call: _e.mock.On("FetchCampaign", id)}
}

func (_c *MockServiceInterface_FetchCampaign_Call) Run(run func(id int)) *MockServiceInterface_FetchCampaign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockServiceInterface_FetchCampaign_Call) Return(_a0 *campaign.Campaign, _a1 error) *MockServiceInterface_FetchCampaign_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockServiceInterface_FetchCampaign_Call) RunAndReturn(run func(int) (*campaign.Campaign, error)) *MockServiceInterface_FetchCampaign_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllCampaigns provides a mock function with given fields:
func (_m *MockServiceInterface) GetAllCampaigns() ([]campaign.Campaign, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAllCampaigns")
	}

	var r0 []campaign.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]campaign.Campaign, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []campaign.Campaign); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]campaign.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockServiceInterface_GetAllCampaigns_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllCampaigns'
type MockServiceInterface_GetAllCampaigns_Call struct {
	*mock.Call
}

// GetAllCampaigns is a helper method to define mock.On call
func (_e *MockServiceInterface_Expecter) GetAllCampaigns() *MockServiceInterface_GetAllCampaigns_Call {
	return &MockServiceInterface_GetAllCampaigns_Call{Call: _e.mock.On("GetAllCampaigns")}
}

func (_c *MockServiceInterface_GetAllCampaigns_Call) Run(run func()) *MockServiceInterface_GetAllCampaigns_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockServiceInterface_GetAllCampaigns_Call) Return(_a0 []campaign.Campaign, _a1 error) *MockServiceInterface_GetAllCampaigns_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockServiceInterface_GetAllCampaigns_Call) RunAndReturn(run func() ([]campaign.Campaign, error)) *MockServiceInterface_GetAllCampaigns_Call {
	_c.Call.Return(run)
	return _c
}

// GetCampaignByID provides a mock function with given fields: id
func (_m *MockServiceInterface) GetCampaignByID(id int) (*campaign.Campaign, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetCampaignByID")
	}

	var r0 *campaign.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*campaign.Campaign, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *campaign.Campaign); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*campaign.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
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
//   - id int
func (_e *MockServiceInterface_Expecter) GetCampaignByID(id interface{}) *MockServiceInterface_GetCampaignByID_Call {
	return &MockServiceInterface_GetCampaignByID_Call{Call: _e.mock.On("GetCampaignByID", id)}
}

func (_c *MockServiceInterface_GetCampaignByID_Call) Run(run func(id int)) *MockServiceInterface_GetCampaignByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockServiceInterface_GetCampaignByID_Call) Return(_a0 *campaign.Campaign, _a1 error) *MockServiceInterface_GetCampaignByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockServiceInterface_GetCampaignByID_Call) RunAndReturn(run func(int) (*campaign.Campaign, error)) *MockServiceInterface_GetCampaignByID_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateCampaign provides a mock function with given fields: _a0
func (_m *MockServiceInterface) UpdateCampaign(_a0 *campaign.Campaign) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for UpdateCampaign")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*campaign.Campaign) error); ok {
		r0 = rf(_a0)
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
//   - _a0 *campaign.Campaign
func (_e *MockServiceInterface_Expecter) UpdateCampaign(_a0 interface{}) *MockServiceInterface_UpdateCampaign_Call {
	return &MockServiceInterface_UpdateCampaign_Call{Call: _e.mock.On("UpdateCampaign", _a0)}
}

func (_c *MockServiceInterface_UpdateCampaign_Call) Run(run func(_a0 *campaign.Campaign)) *MockServiceInterface_UpdateCampaign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*campaign.Campaign))
	})
	return _c
}

func (_c *MockServiceInterface_UpdateCampaign_Call) Return(_a0 error) *MockServiceInterface_UpdateCampaign_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockServiceInterface_UpdateCampaign_Call) RunAndReturn(run func(*campaign.Campaign) error) *MockServiceInterface_UpdateCampaign_Call {
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
