// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	campaign "github.com/fullstackdev42/mp-emailer/campaign"
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

// Create provides a mock function with given fields: dto
func (_m *MockRepositoryInterface) Create(dto *campaign.CreateCampaignDTO) (*campaign.Campaign, error) {
	ret := _m.Called(dto)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *campaign.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(*campaign.CreateCampaignDTO) (*campaign.Campaign, error)); ok {
		return rf(dto)
	}
	if rf, ok := ret.Get(0).(func(*campaign.CreateCampaignDTO) *campaign.Campaign); ok {
		r0 = rf(dto)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*campaign.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func(*campaign.CreateCampaignDTO) error); ok {
		r1 = rf(dto)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepositoryInterface_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockRepositoryInterface_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - dto *campaign.CreateCampaignDTO
func (_e *MockRepositoryInterface_Expecter) Create(dto interface{}) *MockRepositoryInterface_Create_Call {
	return &MockRepositoryInterface_Create_Call{Call: _e.mock.On("Create", dto)}
}

func (_c *MockRepositoryInterface_Create_Call) Run(run func(dto *campaign.CreateCampaignDTO)) *MockRepositoryInterface_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*campaign.CreateCampaignDTO))
	})
	return _c
}

func (_c *MockRepositoryInterface_Create_Call) Return(_a0 *campaign.Campaign, _a1 error) *MockRepositoryInterface_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepositoryInterface_Create_Call) RunAndReturn(run func(*campaign.CreateCampaignDTO) (*campaign.Campaign, error)) *MockRepositoryInterface_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: dto
func (_m *MockRepositoryInterface) Delete(dto campaign.DeleteCampaignDTO) error {
	ret := _m.Called(dto)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(campaign.DeleteCampaignDTO) error); ok {
		r0 = rf(dto)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRepositoryInterface_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockRepositoryInterface_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - dto campaign.DeleteCampaignDTO
func (_e *MockRepositoryInterface_Expecter) Delete(dto interface{}) *MockRepositoryInterface_Delete_Call {
	return &MockRepositoryInterface_Delete_Call{Call: _e.mock.On("Delete", dto)}
}

func (_c *MockRepositoryInterface_Delete_Call) Run(run func(dto campaign.DeleteCampaignDTO)) *MockRepositoryInterface_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(campaign.DeleteCampaignDTO))
	})
	return _c
}

func (_c *MockRepositoryInterface_Delete_Call) Return(_a0 error) *MockRepositoryInterface_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRepositoryInterface_Delete_Call) RunAndReturn(run func(campaign.DeleteCampaignDTO) error) *MockRepositoryInterface_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields:
func (_m *MockRepositoryInterface) GetAll() ([]campaign.Campaign, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
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

// MockRepositoryInterface_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type MockRepositoryInterface_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
func (_e *MockRepositoryInterface_Expecter) GetAll() *MockRepositoryInterface_GetAll_Call {
	return &MockRepositoryInterface_GetAll_Call{Call: _e.mock.On("GetAll")}
}

func (_c *MockRepositoryInterface_GetAll_Call) Run(run func()) *MockRepositoryInterface_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRepositoryInterface_GetAll_Call) Return(_a0 []campaign.Campaign, _a1 error) *MockRepositoryInterface_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepositoryInterface_GetAll_Call) RunAndReturn(run func() ([]campaign.Campaign, error)) *MockRepositoryInterface_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: dto
func (_m *MockRepositoryInterface) GetByID(dto campaign.GetCampaignDTO) (*campaign.Campaign, error) {
	ret := _m.Called(dto)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *campaign.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(campaign.GetCampaignDTO) (*campaign.Campaign, error)); ok {
		return rf(dto)
	}
	if rf, ok := ret.Get(0).(func(campaign.GetCampaignDTO) *campaign.Campaign); ok {
		r0 = rf(dto)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*campaign.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func(campaign.GetCampaignDTO) error); ok {
		r1 = rf(dto)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepositoryInterface_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type MockRepositoryInterface_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - dto campaign.GetCampaignDTO
func (_e *MockRepositoryInterface_Expecter) GetByID(dto interface{}) *MockRepositoryInterface_GetByID_Call {
	return &MockRepositoryInterface_GetByID_Call{Call: _e.mock.On("GetByID", dto)}
}

func (_c *MockRepositoryInterface_GetByID_Call) Run(run func(dto campaign.GetCampaignDTO)) *MockRepositoryInterface_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(campaign.GetCampaignDTO))
	})
	return _c
}

func (_c *MockRepositoryInterface_GetByID_Call) Return(_a0 *campaign.Campaign, _a1 error) *MockRepositoryInterface_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepositoryInterface_GetByID_Call) RunAndReturn(run func(campaign.GetCampaignDTO) (*campaign.Campaign, error)) *MockRepositoryInterface_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetCampaign provides a mock function with given fields: dto
func (_m *MockRepositoryInterface) GetCampaign(dto campaign.GetCampaignDTO) (*campaign.Campaign, error) {
	ret := _m.Called(dto)

	if len(ret) == 0 {
		panic("no return value specified for GetCampaign")
	}

	var r0 *campaign.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(campaign.GetCampaignDTO) (*campaign.Campaign, error)); ok {
		return rf(dto)
	}
	if rf, ok := ret.Get(0).(func(campaign.GetCampaignDTO) *campaign.Campaign); ok {
		r0 = rf(dto)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*campaign.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func(campaign.GetCampaignDTO) error); ok {
		r1 = rf(dto)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepositoryInterface_GetCampaign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCampaign'
type MockRepositoryInterface_GetCampaign_Call struct {
	*mock.Call
}

// GetCampaign is a helper method to define mock.On call
//   - dto campaign.GetCampaignDTO
func (_e *MockRepositoryInterface_Expecter) GetCampaign(dto interface{}) *MockRepositoryInterface_GetCampaign_Call {
	return &MockRepositoryInterface_GetCampaign_Call{Call: _e.mock.On("GetCampaign", dto)}
}

func (_c *MockRepositoryInterface_GetCampaign_Call) Run(run func(dto campaign.GetCampaignDTO)) *MockRepositoryInterface_GetCampaign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(campaign.GetCampaignDTO))
	})
	return _c
}

func (_c *MockRepositoryInterface_GetCampaign_Call) Return(_a0 *campaign.Campaign, _a1 error) *MockRepositoryInterface_GetCampaign_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepositoryInterface_GetCampaign_Call) RunAndReturn(run func(campaign.GetCampaignDTO) (*campaign.Campaign, error)) *MockRepositoryInterface_GetCampaign_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: dto
func (_m *MockRepositoryInterface) Update(dto *campaign.UpdateCampaignDTO) error {
	ret := _m.Called(dto)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*campaign.UpdateCampaignDTO) error); ok {
		r0 = rf(dto)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRepositoryInterface_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockRepositoryInterface_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - dto *campaign.UpdateCampaignDTO
func (_e *MockRepositoryInterface_Expecter) Update(dto interface{}) *MockRepositoryInterface_Update_Call {
	return &MockRepositoryInterface_Update_Call{Call: _e.mock.On("Update", dto)}
}

func (_c *MockRepositoryInterface_Update_Call) Run(run func(dto *campaign.UpdateCampaignDTO)) *MockRepositoryInterface_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*campaign.UpdateCampaignDTO))
	})
	return _c
}

func (_c *MockRepositoryInterface_Update_Call) Return(_a0 error) *MockRepositoryInterface_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRepositoryInterface_Update_Call) RunAndReturn(run func(*campaign.UpdateCampaignDTO) error) *MockRepositoryInterface_Update_Call {
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