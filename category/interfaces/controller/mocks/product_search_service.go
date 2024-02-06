// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	application "flamingo.me/flamingo-commerce/v3/search/application"

	mock "github.com/stretchr/testify/mock"

	productapplication "flamingo.me/flamingo-commerce/v3/product/application"
)

// ProductSearchService is an autogenerated mock type for the ProductSearchService type
type ProductSearchService struct {
	mock.Mock
}

type ProductSearchService_Expecter struct {
	mock *mock.Mock
}

func (_m *ProductSearchService) EXPECT() *ProductSearchService_Expecter {
	return &ProductSearchService_Expecter{mock: &_m.Mock}
}

// Find provides a mock function with given fields: ctx, searchRequest
func (_m *ProductSearchService) Find(ctx context.Context, searchRequest *application.SearchRequest) (*productapplication.SearchResult, error) {
	ret := _m.Called(ctx, searchRequest)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 *productapplication.SearchResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *application.SearchRequest) (*productapplication.SearchResult, error)); ok {
		return rf(ctx, searchRequest)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *application.SearchRequest) *productapplication.SearchResult); ok {
		r0 = rf(ctx, searchRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*productapplication.SearchResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *application.SearchRequest) error); ok {
		r1 = rf(ctx, searchRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductSearchService_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type ProductSearchService_Find_Call struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - ctx context.Context
//   - searchRequest *application.SearchRequest
func (_e *ProductSearchService_Expecter) Find(ctx interface{}, searchRequest interface{}) *ProductSearchService_Find_Call {
	return &ProductSearchService_Find_Call{Call: _e.mock.On("Find", ctx, searchRequest)}
}

func (_c *ProductSearchService_Find_Call) Run(run func(ctx context.Context, searchRequest *application.SearchRequest)) *ProductSearchService_Find_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*application.SearchRequest))
	})
	return _c
}

func (_c *ProductSearchService_Find_Call) Return(_a0 *productapplication.SearchResult, _a1 error) *ProductSearchService_Find_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductSearchService_Find_Call) RunAndReturn(run func(context.Context, *application.SearchRequest) (*productapplication.SearchResult, error)) *ProductSearchService_Find_Call {
	_c.Call.Return(run)
	return _c
}

// NewProductSearchService creates a new instance of ProductSearchService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProductSearchService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProductSearchService {
	mock := &ProductSearchService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
