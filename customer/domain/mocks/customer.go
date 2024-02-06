// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	domain "flamingo.me/flamingo-commerce/v3/customer/domain"
	mock "github.com/stretchr/testify/mock"
)

// Customer is an autogenerated mock type for the Customer type
type Customer struct {
	mock.Mock
}

type Customer_Expecter struct {
	mock *mock.Mock
}

func (_m *Customer) EXPECT() *Customer_Expecter {
	return &Customer_Expecter{mock: &_m.Mock}
}

// GetAddresses provides a mock function with given fields:
func (_m *Customer) GetAddresses() []domain.Address {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAddresses")
	}

	var r0 []domain.Address
	if rf, ok := ret.Get(0).(func() []domain.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Address)
		}
	}

	return r0
}

// Customer_GetAddresses_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAddresses'
type Customer_GetAddresses_Call struct {
	*mock.Call
}

// GetAddresses is a helper method to define mock.On call
func (_e *Customer_Expecter) GetAddresses() *Customer_GetAddresses_Call {
	return &Customer_GetAddresses_Call{Call: _e.mock.On("GetAddresses")}
}

func (_c *Customer_GetAddresses_Call) Run(run func()) *Customer_GetAddresses_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Customer_GetAddresses_Call) Return(_a0 []domain.Address) *Customer_GetAddresses_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Customer_GetAddresses_Call) RunAndReturn(run func() []domain.Address) *Customer_GetAddresses_Call {
	_c.Call.Return(run)
	return _c
}

// GetDefaultBillingAddress provides a mock function with given fields:
func (_m *Customer) GetDefaultBillingAddress() *domain.Address {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetDefaultBillingAddress")
	}

	var r0 *domain.Address
	if rf, ok := ret.Get(0).(func() *domain.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Address)
		}
	}

	return r0
}

// Customer_GetDefaultBillingAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDefaultBillingAddress'
type Customer_GetDefaultBillingAddress_Call struct {
	*mock.Call
}

// GetDefaultBillingAddress is a helper method to define mock.On call
func (_e *Customer_Expecter) GetDefaultBillingAddress() *Customer_GetDefaultBillingAddress_Call {
	return &Customer_GetDefaultBillingAddress_Call{Call: _e.mock.On("GetDefaultBillingAddress")}
}

func (_c *Customer_GetDefaultBillingAddress_Call) Run(run func()) *Customer_GetDefaultBillingAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Customer_GetDefaultBillingAddress_Call) Return(_a0 *domain.Address) *Customer_GetDefaultBillingAddress_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Customer_GetDefaultBillingAddress_Call) RunAndReturn(run func() *domain.Address) *Customer_GetDefaultBillingAddress_Call {
	_c.Call.Return(run)
	return _c
}

// GetDefaultShippingAddress provides a mock function with given fields:
func (_m *Customer) GetDefaultShippingAddress() *domain.Address {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetDefaultShippingAddress")
	}

	var r0 *domain.Address
	if rf, ok := ret.Get(0).(func() *domain.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Address)
		}
	}

	return r0
}

// Customer_GetDefaultShippingAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDefaultShippingAddress'
type Customer_GetDefaultShippingAddress_Call struct {
	*mock.Call
}

// GetDefaultShippingAddress is a helper method to define mock.On call
func (_e *Customer_Expecter) GetDefaultShippingAddress() *Customer_GetDefaultShippingAddress_Call {
	return &Customer_GetDefaultShippingAddress_Call{Call: _e.mock.On("GetDefaultShippingAddress")}
}

func (_c *Customer_GetDefaultShippingAddress_Call) Run(run func()) *Customer_GetDefaultShippingAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Customer_GetDefaultShippingAddress_Call) Return(_a0 *domain.Address) *Customer_GetDefaultShippingAddress_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Customer_GetDefaultShippingAddress_Call) RunAndReturn(run func() *domain.Address) *Customer_GetDefaultShippingAddress_Call {
	_c.Call.Return(run)
	return _c
}

// GetID provides a mock function with given fields:
func (_m *Customer) GetID() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetID")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Customer_GetID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetID'
type Customer_GetID_Call struct {
	*mock.Call
}

// GetID is a helper method to define mock.On call
func (_e *Customer_Expecter) GetID() *Customer_GetID_Call {
	return &Customer_GetID_Call{Call: _e.mock.On("GetID")}
}

func (_c *Customer_GetID_Call) Run(run func()) *Customer_GetID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Customer_GetID_Call) Return(_a0 string) *Customer_GetID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Customer_GetID_Call) RunAndReturn(run func() string) *Customer_GetID_Call {
	_c.Call.Return(run)
	return _c
}

// GetPersonalData provides a mock function with given fields:
func (_m *Customer) GetPersonalData() domain.PersonData {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetPersonalData")
	}

	var r0 domain.PersonData
	if rf, ok := ret.Get(0).(func() domain.PersonData); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(domain.PersonData)
	}

	return r0
}

// Customer_GetPersonalData_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPersonalData'
type Customer_GetPersonalData_Call struct {
	*mock.Call
}

// GetPersonalData is a helper method to define mock.On call
func (_e *Customer_Expecter) GetPersonalData() *Customer_GetPersonalData_Call {
	return &Customer_GetPersonalData_Call{Call: _e.mock.On("GetPersonalData")}
}

func (_c *Customer_GetPersonalData_Call) Run(run func()) *Customer_GetPersonalData_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Customer_GetPersonalData_Call) Return(_a0 domain.PersonData) *Customer_GetPersonalData_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Customer_GetPersonalData_Call) RunAndReturn(run func() domain.PersonData) *Customer_GetPersonalData_Call {
	_c.Call.Return(run)
	return _c
}

// NewCustomer creates a new instance of Customer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCustomer(t interface {
	mock.TestingT
	Cleanup(func())
}) *Customer {
	mock := &Customer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
