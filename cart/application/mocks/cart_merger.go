// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"

	cart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"

	mock "github.com/stretchr/testify/mock"

	web "flamingo.me/flamingo/v3/framework/web"
)

// CartMerger is an autogenerated mock type for the CartMerger type
type CartMerger struct {
	mock.Mock
}

type CartMerger_Expecter struct {
	mock *mock.Mock
}

func (_m *CartMerger) EXPECT() *CartMerger_Expecter {
	return &CartMerger_Expecter{mock: &_m.Mock}
}

// Merge provides a mock function with given fields: ctx, session, guestCart, customerCart
func (_m *CartMerger) Merge(ctx context.Context, session *web.Session, guestCart cart.Cart, customerCart cart.Cart) {
	_m.Called(ctx, session, guestCart, customerCart)
}

// CartMerger_Merge_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Merge'
type CartMerger_Merge_Call struct {
	*mock.Call
}

// Merge is a helper method to define mock.On call
//   - ctx context.Context
//   - session *web.Session
//   - guestCart cart.Cart
//   - customerCart cart.Cart
func (_e *CartMerger_Expecter) Merge(ctx interface{}, session interface{}, guestCart interface{}, customerCart interface{}) *CartMerger_Merge_Call {
	return &CartMerger_Merge_Call{Call: _e.mock.On("Merge", ctx, session, guestCart, customerCart)}
}

func (_c *CartMerger_Merge_Call) Run(run func(ctx context.Context, session *web.Session, guestCart cart.Cart, customerCart cart.Cart)) *CartMerger_Merge_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*web.Session), args[2].(cart.Cart), args[3].(cart.Cart))
	})
	return _c
}

func (_c *CartMerger_Merge_Call) Return() *CartMerger_Merge_Call {
	_c.Call.Return()
	return _c
}

func (_c *CartMerger_Merge_Call) RunAndReturn(run func(context.Context, *web.Session, cart.Cart, cart.Cart)) *CartMerger_Merge_Call {
	_c.Call.Return(run)
	return _c
}

// NewCartMerger creates a new instance of CartMerger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCartMerger(t interface {
	mock.TestingT
	Cleanup(func())
}) *CartMerger {
	mock := &CartMerger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
