// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	cart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"

	mock "github.com/stretchr/testify/mock"
)

// GiftCardBehaviour is an autogenerated mock type for the GiftCardBehaviour type
type GiftCardBehaviour struct {
	mock.Mock
}

type GiftCardBehaviour_Expecter struct {
	mock *mock.Mock
}

func (_m *GiftCardBehaviour) EXPECT() *GiftCardBehaviour_Expecter {
	return &GiftCardBehaviour_Expecter{mock: &_m.Mock}
}

// ApplyGiftCard provides a mock function with given fields: ctx, _a1, giftCardCode
func (_m *GiftCardBehaviour) ApplyGiftCard(ctx context.Context, _a1 *cart.Cart, giftCardCode string) (*cart.Cart, cart.DeferEvents, error) {
	ret := _m.Called(ctx, _a1, giftCardCode)

	if len(ret) == 0 {
		panic("no return value specified for ApplyGiftCard")
	}

	var r0 *cart.Cart
	var r1 cart.DeferEvents
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart, string) (*cart.Cart, cart.DeferEvents, error)); ok {
		return rf(ctx, _a1, giftCardCode)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart, string) *cart.Cart); ok {
		r0 = rf(ctx, _a1, giftCardCode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cart.Cart)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cart.Cart, string) cart.DeferEvents); ok {
		r1 = rf(ctx, _a1, giftCardCode)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(cart.DeferEvents)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, *cart.Cart, string) error); ok {
		r2 = rf(ctx, _a1, giftCardCode)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GiftCardBehaviour_ApplyGiftCard_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ApplyGiftCard'
type GiftCardBehaviour_ApplyGiftCard_Call struct {
	*mock.Call
}

// ApplyGiftCard is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 *cart.Cart
//   - giftCardCode string
func (_e *GiftCardBehaviour_Expecter) ApplyGiftCard(ctx interface{}, _a1 interface{}, giftCardCode interface{}) *GiftCardBehaviour_ApplyGiftCard_Call {
	return &GiftCardBehaviour_ApplyGiftCard_Call{Call: _e.mock.On("ApplyGiftCard", ctx, _a1, giftCardCode)}
}

func (_c *GiftCardBehaviour_ApplyGiftCard_Call) Run(run func(ctx context.Context, _a1 *cart.Cart, giftCardCode string)) *GiftCardBehaviour_ApplyGiftCard_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cart.Cart), args[2].(string))
	})
	return _c
}

func (_c *GiftCardBehaviour_ApplyGiftCard_Call) Return(_a0 *cart.Cart, _a1 cart.DeferEvents, _a2 error) *GiftCardBehaviour_ApplyGiftCard_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *GiftCardBehaviour_ApplyGiftCard_Call) RunAndReturn(run func(context.Context, *cart.Cart, string) (*cart.Cart, cart.DeferEvents, error)) *GiftCardBehaviour_ApplyGiftCard_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveGiftCard provides a mock function with given fields: ctx, _a1, giftCardCode
func (_m *GiftCardBehaviour) RemoveGiftCard(ctx context.Context, _a1 *cart.Cart, giftCardCode string) (*cart.Cart, cart.DeferEvents, error) {
	ret := _m.Called(ctx, _a1, giftCardCode)

	if len(ret) == 0 {
		panic("no return value specified for RemoveGiftCard")
	}

	var r0 *cart.Cart
	var r1 cart.DeferEvents
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart, string) (*cart.Cart, cart.DeferEvents, error)); ok {
		return rf(ctx, _a1, giftCardCode)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart, string) *cart.Cart); ok {
		r0 = rf(ctx, _a1, giftCardCode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cart.Cart)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cart.Cart, string) cart.DeferEvents); ok {
		r1 = rf(ctx, _a1, giftCardCode)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(cart.DeferEvents)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, *cart.Cart, string) error); ok {
		r2 = rf(ctx, _a1, giftCardCode)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GiftCardBehaviour_RemoveGiftCard_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveGiftCard'
type GiftCardBehaviour_RemoveGiftCard_Call struct {
	*mock.Call
}

// RemoveGiftCard is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 *cart.Cart
//   - giftCardCode string
func (_e *GiftCardBehaviour_Expecter) RemoveGiftCard(ctx interface{}, _a1 interface{}, giftCardCode interface{}) *GiftCardBehaviour_RemoveGiftCard_Call {
	return &GiftCardBehaviour_RemoveGiftCard_Call{Call: _e.mock.On("RemoveGiftCard", ctx, _a1, giftCardCode)}
}

func (_c *GiftCardBehaviour_RemoveGiftCard_Call) Run(run func(ctx context.Context, _a1 *cart.Cart, giftCardCode string)) *GiftCardBehaviour_RemoveGiftCard_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cart.Cart), args[2].(string))
	})
	return _c
}

func (_c *GiftCardBehaviour_RemoveGiftCard_Call) Return(_a0 *cart.Cart, _a1 cart.DeferEvents, _a2 error) *GiftCardBehaviour_RemoveGiftCard_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *GiftCardBehaviour_RemoveGiftCard_Call) RunAndReturn(run func(context.Context, *cart.Cart, string) (*cart.Cart, cart.DeferEvents, error)) *GiftCardBehaviour_RemoveGiftCard_Call {
	_c.Call.Return(run)
	return _c
}

// NewGiftCardBehaviour creates a new instance of GiftCardBehaviour. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGiftCardBehaviour(t interface {
	mock.TestingT
	Cleanup(func())
}) *GiftCardBehaviour {
	mock := &GiftCardBehaviour{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
