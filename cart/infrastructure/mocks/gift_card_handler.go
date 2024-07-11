// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	cart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"

	mock "github.com/stretchr/testify/mock"
)

// GiftCardHandler is an autogenerated mock type for the GiftCardHandler type
type GiftCardHandler struct {
	mock.Mock
}

type GiftCardHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *GiftCardHandler) EXPECT() *GiftCardHandler_Expecter {
	return &GiftCardHandler_Expecter{mock: &_m.Mock}
}

// ApplyGiftCard provides a mock function with given fields: ctx, _a1, giftCardCode
func (_m *GiftCardHandler) ApplyGiftCard(ctx context.Context, _a1 *cart.Cart, giftCardCode string) (*cart.Cart, error) {
	ret := _m.Called(ctx, _a1, giftCardCode)

	if len(ret) == 0 {
		panic("no return value specified for ApplyGiftCard")
	}

	var r0 *cart.Cart
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart, string) (*cart.Cart, error)); ok {
		return rf(ctx, _a1, giftCardCode)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart, string) *cart.Cart); ok {
		r0 = rf(ctx, _a1, giftCardCode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cart.Cart)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cart.Cart, string) error); ok {
		r1 = rf(ctx, _a1, giftCardCode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GiftCardHandler_ApplyGiftCard_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ApplyGiftCard'
type GiftCardHandler_ApplyGiftCard_Call struct {
	*mock.Call
}

// ApplyGiftCard is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 *cart.Cart
//   - giftCardCode string
func (_e *GiftCardHandler_Expecter) ApplyGiftCard(ctx interface{}, _a1 interface{}, giftCardCode interface{}) *GiftCardHandler_ApplyGiftCard_Call {
	return &GiftCardHandler_ApplyGiftCard_Call{Call: _e.mock.On("ApplyGiftCard", ctx, _a1, giftCardCode)}
}

func (_c *GiftCardHandler_ApplyGiftCard_Call) Run(run func(ctx context.Context, _a1 *cart.Cart, giftCardCode string)) *GiftCardHandler_ApplyGiftCard_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cart.Cart), args[2].(string))
	})
	return _c
}

func (_c *GiftCardHandler_ApplyGiftCard_Call) Return(_a0 *cart.Cart, _a1 error) *GiftCardHandler_ApplyGiftCard_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *GiftCardHandler_ApplyGiftCard_Call) RunAndReturn(run func(context.Context, *cart.Cart, string) (*cart.Cart, error)) *GiftCardHandler_ApplyGiftCard_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveGiftCard provides a mock function with given fields: ctx, _a1, giftCardCode
func (_m *GiftCardHandler) RemoveGiftCard(ctx context.Context, _a1 *cart.Cart, giftCardCode string) (*cart.Cart, error) {
	ret := _m.Called(ctx, _a1, giftCardCode)

	if len(ret) == 0 {
		panic("no return value specified for RemoveGiftCard")
	}

	var r0 *cart.Cart
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart, string) (*cart.Cart, error)); ok {
		return rf(ctx, _a1, giftCardCode)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart, string) *cart.Cart); ok {
		r0 = rf(ctx, _a1, giftCardCode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cart.Cart)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cart.Cart, string) error); ok {
		r1 = rf(ctx, _a1, giftCardCode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GiftCardHandler_RemoveGiftCard_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveGiftCard'
type GiftCardHandler_RemoveGiftCard_Call struct {
	*mock.Call
}

// RemoveGiftCard is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 *cart.Cart
//   - giftCardCode string
func (_e *GiftCardHandler_Expecter) RemoveGiftCard(ctx interface{}, _a1 interface{}, giftCardCode interface{}) *GiftCardHandler_RemoveGiftCard_Call {
	return &GiftCardHandler_RemoveGiftCard_Call{Call: _e.mock.On("RemoveGiftCard", ctx, _a1, giftCardCode)}
}

func (_c *GiftCardHandler_RemoveGiftCard_Call) Run(run func(ctx context.Context, _a1 *cart.Cart, giftCardCode string)) *GiftCardHandler_RemoveGiftCard_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cart.Cart), args[2].(string))
	})
	return _c
}

func (_c *GiftCardHandler_RemoveGiftCard_Call) Return(_a0 *cart.Cart, _a1 error) *GiftCardHandler_RemoveGiftCard_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *GiftCardHandler_RemoveGiftCard_Call) RunAndReturn(run func(context.Context, *cart.Cart, string) (*cart.Cart, error)) *GiftCardHandler_RemoveGiftCard_Call {
	_c.Call.Return(run)
	return _c
}

// NewGiftCardHandler creates a new instance of GiftCardHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGiftCardHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *GiftCardHandler {
	mock := &GiftCardHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
