// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	cart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"

	mock "github.com/stretchr/testify/mock"
)

// CompleteBehaviour is an autogenerated mock type for the CompleteBehaviour type
type CompleteBehaviour struct {
	mock.Mock
}

type CompleteBehaviour_Expecter struct {
	mock *mock.Mock
}

func (_m *CompleteBehaviour) EXPECT() *CompleteBehaviour_Expecter {
	return &CompleteBehaviour_Expecter{mock: &_m.Mock}
}

// Complete provides a mock function with given fields: _a0, _a1
func (_m *CompleteBehaviour) Complete(_a0 context.Context, _a1 *cart.Cart) (*cart.Cart, cart.DeferEvents, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Complete")
	}

	var r0 *cart.Cart
	var r1 cart.DeferEvents
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart) (*cart.Cart, cart.DeferEvents, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart) *cart.Cart); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cart.Cart)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cart.Cart) cart.DeferEvents); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(cart.DeferEvents)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, *cart.Cart) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CompleteBehaviour_Complete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Complete'
type CompleteBehaviour_Complete_Call struct {
	*mock.Call
}

// Complete is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *cart.Cart
func (_e *CompleteBehaviour_Expecter) Complete(_a0 interface{}, _a1 interface{}) *CompleteBehaviour_Complete_Call {
	return &CompleteBehaviour_Complete_Call{Call: _e.mock.On("Complete", _a0, _a1)}
}

func (_c *CompleteBehaviour_Complete_Call) Run(run func(_a0 context.Context, _a1 *cart.Cart)) *CompleteBehaviour_Complete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cart.Cart))
	})
	return _c
}

func (_c *CompleteBehaviour_Complete_Call) Return(_a0 *cart.Cart, _a1 cart.DeferEvents, _a2 error) *CompleteBehaviour_Complete_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *CompleteBehaviour_Complete_Call) RunAndReturn(run func(context.Context, *cart.Cart) (*cart.Cart, cart.DeferEvents, error)) *CompleteBehaviour_Complete_Call {
	_c.Call.Return(run)
	return _c
}

// Restore provides a mock function with given fields: _a0, _a1
func (_m *CompleteBehaviour) Restore(_a0 context.Context, _a1 *cart.Cart) (*cart.Cart, cart.DeferEvents, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Restore")
	}

	var r0 *cart.Cart
	var r1 cart.DeferEvents
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart) (*cart.Cart, cart.DeferEvents, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart) *cart.Cart); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cart.Cart)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cart.Cart) cart.DeferEvents); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(cart.DeferEvents)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, *cart.Cart) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CompleteBehaviour_Restore_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Restore'
type CompleteBehaviour_Restore_Call struct {
	*mock.Call
}

// Restore is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *cart.Cart
func (_e *CompleteBehaviour_Expecter) Restore(_a0 interface{}, _a1 interface{}) *CompleteBehaviour_Restore_Call {
	return &CompleteBehaviour_Restore_Call{Call: _e.mock.On("Restore", _a0, _a1)}
}

func (_c *CompleteBehaviour_Restore_Call) Run(run func(_a0 context.Context, _a1 *cart.Cart)) *CompleteBehaviour_Restore_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cart.Cart))
	})
	return _c
}

func (_c *CompleteBehaviour_Restore_Call) Return(_a0 *cart.Cart, _a1 cart.DeferEvents, _a2 error) *CompleteBehaviour_Restore_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *CompleteBehaviour_Restore_Call) RunAndReturn(run func(context.Context, *cart.Cart) (*cart.Cart, cart.DeferEvents, error)) *CompleteBehaviour_Restore_Call {
	_c.Call.Return(run)
	return _c
}

// NewCompleteBehaviour creates a new instance of CompleteBehaviour. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCompleteBehaviour(t interface {
	mock.TestingT
	Cleanup(func())
}) *CompleteBehaviour {
	mock := &CompleteBehaviour{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
