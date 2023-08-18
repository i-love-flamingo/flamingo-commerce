// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"

	cart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"

	domain "flamingo.me/flamingo-commerce/v3/payment/domain"

	mock "github.com/stretchr/testify/mock"

	placeorder "flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"

	url "net/url"
)

// WebCartPaymentGateway is an autogenerated mock type for the WebCartPaymentGateway type
type WebCartPaymentGateway struct {
	mock.Mock
}

type WebCartPaymentGateway_Expecter struct {
	mock *mock.Mock
}

func (_m *WebCartPaymentGateway) EXPECT() *WebCartPaymentGateway_Expecter {
	return &WebCartPaymentGateway_Expecter{mock: &_m.Mock}
}

// CancelOrderPayment provides a mock function with given fields: ctx, cartPayment
func (_m *WebCartPaymentGateway) CancelOrderPayment(ctx context.Context, cartPayment *placeorder.Payment) error {
	ret := _m.Called(ctx, cartPayment)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *placeorder.Payment) error); ok {
		r0 = rf(ctx, cartPayment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WebCartPaymentGateway_CancelOrderPayment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CancelOrderPayment'
type WebCartPaymentGateway_CancelOrderPayment_Call struct {
	*mock.Call
}

// CancelOrderPayment is a helper method to define mock.On call
//   - ctx context.Context
//   - cartPayment *placeorder.Payment
func (_e *WebCartPaymentGateway_Expecter) CancelOrderPayment(ctx interface{}, cartPayment interface{}) *WebCartPaymentGateway_CancelOrderPayment_Call {
	return &WebCartPaymentGateway_CancelOrderPayment_Call{Call: _e.mock.On("CancelOrderPayment", ctx, cartPayment)}
}

func (_c *WebCartPaymentGateway_CancelOrderPayment_Call) Run(run func(ctx context.Context, cartPayment *placeorder.Payment)) *WebCartPaymentGateway_CancelOrderPayment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*placeorder.Payment))
	})
	return _c
}

func (_c *WebCartPaymentGateway_CancelOrderPayment_Call) Return(_a0 error) *WebCartPaymentGateway_CancelOrderPayment_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *WebCartPaymentGateway_CancelOrderPayment_Call) RunAndReturn(run func(context.Context, *placeorder.Payment) error) *WebCartPaymentGateway_CancelOrderPayment_Call {
	_c.Call.Return(run)
	return _c
}

// ConfirmResult provides a mock function with given fields: ctx, _a1, cartPayment
func (_m *WebCartPaymentGateway) ConfirmResult(ctx context.Context, _a1 *cart.Cart, cartPayment *placeorder.Payment) error {
	ret := _m.Called(ctx, _a1, cartPayment)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart, *placeorder.Payment) error); ok {
		r0 = rf(ctx, _a1, cartPayment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WebCartPaymentGateway_ConfirmResult_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConfirmResult'
type WebCartPaymentGateway_ConfirmResult_Call struct {
	*mock.Call
}

// ConfirmResult is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 *cart.Cart
//   - cartPayment *placeorder.Payment
func (_e *WebCartPaymentGateway_Expecter) ConfirmResult(ctx interface{}, _a1 interface{}, cartPayment interface{}) *WebCartPaymentGateway_ConfirmResult_Call {
	return &WebCartPaymentGateway_ConfirmResult_Call{Call: _e.mock.On("ConfirmResult", ctx, _a1, cartPayment)}
}

func (_c *WebCartPaymentGateway_ConfirmResult_Call) Run(run func(ctx context.Context, _a1 *cart.Cart, cartPayment *placeorder.Payment)) *WebCartPaymentGateway_ConfirmResult_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cart.Cart), args[2].(*placeorder.Payment))
	})
	return _c
}

func (_c *WebCartPaymentGateway_ConfirmResult_Call) Return(_a0 error) *WebCartPaymentGateway_ConfirmResult_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *WebCartPaymentGateway_ConfirmResult_Call) RunAndReturn(run func(context.Context, *cart.Cart, *placeorder.Payment) error) *WebCartPaymentGateway_ConfirmResult_Call {
	_c.Call.Return(run)
	return _c
}

// FlowStatus provides a mock function with given fields: ctx, _a1, correlationID
func (_m *WebCartPaymentGateway) FlowStatus(ctx context.Context, _a1 *cart.Cart, correlationID string) (*domain.FlowStatus, error) {
	ret := _m.Called(ctx, _a1, correlationID)

	var r0 *domain.FlowStatus
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart, string) (*domain.FlowStatus, error)); ok {
		return rf(ctx, _a1, correlationID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart, string) *domain.FlowStatus); ok {
		r0 = rf(ctx, _a1, correlationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.FlowStatus)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cart.Cart, string) error); ok {
		r1 = rf(ctx, _a1, correlationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WebCartPaymentGateway_FlowStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FlowStatus'
type WebCartPaymentGateway_FlowStatus_Call struct {
	*mock.Call
}

// FlowStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 *cart.Cart
//   - correlationID string
func (_e *WebCartPaymentGateway_Expecter) FlowStatus(ctx interface{}, _a1 interface{}, correlationID interface{}) *WebCartPaymentGateway_FlowStatus_Call {
	return &WebCartPaymentGateway_FlowStatus_Call{Call: _e.mock.On("FlowStatus", ctx, _a1, correlationID)}
}

func (_c *WebCartPaymentGateway_FlowStatus_Call) Run(run func(ctx context.Context, _a1 *cart.Cart, correlationID string)) *WebCartPaymentGateway_FlowStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cart.Cart), args[2].(string))
	})
	return _c
}

func (_c *WebCartPaymentGateway_FlowStatus_Call) Return(_a0 *domain.FlowStatus, _a1 error) *WebCartPaymentGateway_FlowStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *WebCartPaymentGateway_FlowStatus_Call) RunAndReturn(run func(context.Context, *cart.Cart, string) (*domain.FlowStatus, error)) *WebCartPaymentGateway_FlowStatus_Call {
	_c.Call.Return(run)
	return _c
}

// Methods provides a mock function with given fields:
func (_m *WebCartPaymentGateway) Methods() []domain.Method {
	ret := _m.Called()

	var r0 []domain.Method
	if rf, ok := ret.Get(0).(func() []domain.Method); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Method)
		}
	}

	return r0
}

// WebCartPaymentGateway_Methods_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Methods'
type WebCartPaymentGateway_Methods_Call struct {
	*mock.Call
}

// Methods is a helper method to define mock.On call
func (_e *WebCartPaymentGateway_Expecter) Methods() *WebCartPaymentGateway_Methods_Call {
	return &WebCartPaymentGateway_Methods_Call{Call: _e.mock.On("Methods")}
}

func (_c *WebCartPaymentGateway_Methods_Call) Run(run func()) *WebCartPaymentGateway_Methods_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *WebCartPaymentGateway_Methods_Call) Return(_a0 []domain.Method) *WebCartPaymentGateway_Methods_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *WebCartPaymentGateway_Methods_Call) RunAndReturn(run func() []domain.Method) *WebCartPaymentGateway_Methods_Call {
	_c.Call.Return(run)
	return _c
}

// OrderPaymentFromFlow provides a mock function with given fields: ctx, _a1, correlationID
func (_m *WebCartPaymentGateway) OrderPaymentFromFlow(ctx context.Context, _a1 *cart.Cart, correlationID string) (*placeorder.Payment, error) {
	ret := _m.Called(ctx, _a1, correlationID)

	var r0 *placeorder.Payment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart, string) (*placeorder.Payment, error)); ok {
		return rf(ctx, _a1, correlationID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart, string) *placeorder.Payment); ok {
		r0 = rf(ctx, _a1, correlationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*placeorder.Payment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cart.Cart, string) error); ok {
		r1 = rf(ctx, _a1, correlationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WebCartPaymentGateway_OrderPaymentFromFlow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OrderPaymentFromFlow'
type WebCartPaymentGateway_OrderPaymentFromFlow_Call struct {
	*mock.Call
}

// OrderPaymentFromFlow is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 *cart.Cart
//   - correlationID string
func (_e *WebCartPaymentGateway_Expecter) OrderPaymentFromFlow(ctx interface{}, _a1 interface{}, correlationID interface{}) *WebCartPaymentGateway_OrderPaymentFromFlow_Call {
	return &WebCartPaymentGateway_OrderPaymentFromFlow_Call{Call: _e.mock.On("OrderPaymentFromFlow", ctx, _a1, correlationID)}
}

func (_c *WebCartPaymentGateway_OrderPaymentFromFlow_Call) Run(run func(ctx context.Context, _a1 *cart.Cart, correlationID string)) *WebCartPaymentGateway_OrderPaymentFromFlow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cart.Cart), args[2].(string))
	})
	return _c
}

func (_c *WebCartPaymentGateway_OrderPaymentFromFlow_Call) Return(_a0 *placeorder.Payment, _a1 error) *WebCartPaymentGateway_OrderPaymentFromFlow_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *WebCartPaymentGateway_OrderPaymentFromFlow_Call) RunAndReturn(run func(context.Context, *cart.Cart, string) (*placeorder.Payment, error)) *WebCartPaymentGateway_OrderPaymentFromFlow_Call {
	_c.Call.Return(run)
	return _c
}

// StartFlow provides a mock function with given fields: ctx, _a1, correlationID, returnURL
func (_m *WebCartPaymentGateway) StartFlow(ctx context.Context, _a1 *cart.Cart, correlationID string, returnURL *url.URL) (*domain.FlowResult, error) {
	ret := _m.Called(ctx, _a1, correlationID, returnURL)

	var r0 *domain.FlowResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart, string, *url.URL) (*domain.FlowResult, error)); ok {
		return rf(ctx, _a1, correlationID, returnURL)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *cart.Cart, string, *url.URL) *domain.FlowResult); ok {
		r0 = rf(ctx, _a1, correlationID, returnURL)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.FlowResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *cart.Cart, string, *url.URL) error); ok {
		r1 = rf(ctx, _a1, correlationID, returnURL)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WebCartPaymentGateway_StartFlow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StartFlow'
type WebCartPaymentGateway_StartFlow_Call struct {
	*mock.Call
}

// StartFlow is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 *cart.Cart
//   - correlationID string
//   - returnURL *url.URL
func (_e *WebCartPaymentGateway_Expecter) StartFlow(ctx interface{}, _a1 interface{}, correlationID interface{}, returnURL interface{}) *WebCartPaymentGateway_StartFlow_Call {
	return &WebCartPaymentGateway_StartFlow_Call{Call: _e.mock.On("StartFlow", ctx, _a1, correlationID, returnURL)}
}

func (_c *WebCartPaymentGateway_StartFlow_Call) Run(run func(ctx context.Context, _a1 *cart.Cart, correlationID string, returnURL *url.URL)) *WebCartPaymentGateway_StartFlow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*cart.Cart), args[2].(string), args[3].(*url.URL))
	})
	return _c
}

func (_c *WebCartPaymentGateway_StartFlow_Call) Return(_a0 *domain.FlowResult, _a1 error) *WebCartPaymentGateway_StartFlow_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *WebCartPaymentGateway_StartFlow_Call) RunAndReturn(run func(context.Context, *cart.Cart, string, *url.URL) (*domain.FlowResult, error)) *WebCartPaymentGateway_StartFlow_Call {
	_c.Call.Return(run)
	return _c
}

// NewWebCartPaymentGateway creates a new instance of WebCartPaymentGateway. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWebCartPaymentGateway(t interface {
	mock.TestingT
	Cleanup(func())
}) *WebCartPaymentGateway {
	mock := &WebCartPaymentGateway{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
