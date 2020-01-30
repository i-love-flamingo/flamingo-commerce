// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import cart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
import context "context"
import mock "github.com/stretchr/testify/mock"

// CompleteBehaviour is an autogenerated mock type for the CompleteBehaviour type
type CompleteBehaviour struct {
	mock.Mock
}

// Complete provides a mock function with given fields: _a0, _a1
func (_m *CompleteBehaviour) Complete(_a0 context.Context, _a1 cart.Cart) (*cart.Cart, cart.DeferEvents, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *cart.Cart
	if rf, ok := ret.Get(0).(func(context.Context, cart.Cart) *cart.Cart); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cart.Cart)
		}
	}

	var r1 cart.DeferEvents
	if rf, ok := ret.Get(1).(func(context.Context, cart.Cart) cart.DeferEvents); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(cart.DeferEvents)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, cart.Cart) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Restore provides a mock function with given fields: _a0, _a1
func (_m *CompleteBehaviour) Restore(_a0 context.Context, _a1 cart.Cart) (*cart.Cart, cart.DeferEvents, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *cart.Cart
	if rf, ok := ret.Get(0).(func(context.Context, cart.Cart) *cart.Cart); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cart.Cart)
		}
	}

	var r1 cart.DeferEvents
	if rf, ok := ret.Get(1).(func(context.Context, cart.Cart) cart.DeferEvents); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(cart.DeferEvents)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, cart.Cart) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
