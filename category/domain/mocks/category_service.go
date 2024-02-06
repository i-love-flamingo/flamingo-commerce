// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "flamingo.me/flamingo-commerce/v3/category/domain"
	mock "github.com/stretchr/testify/mock"
)

// CategoryService is an autogenerated mock type for the CategoryService type
type CategoryService struct {
	mock.Mock
}

type CategoryService_Expecter struct {
	mock *mock.Mock
}

func (_m *CategoryService) EXPECT() *CategoryService_Expecter {
	return &CategoryService_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: ctx, categoryCode
func (_m *CategoryService) Get(ctx context.Context, categoryCode string) (domain.Category, error) {
	ret := _m.Called(ctx, categoryCode)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 domain.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.Category, error)); ok {
		return rf(ctx, categoryCode)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Category); ok {
		r0 = rf(ctx, categoryCode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, categoryCode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CategoryService_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type CategoryService_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - categoryCode string
func (_e *CategoryService_Expecter) Get(ctx interface{}, categoryCode interface{}) *CategoryService_Get_Call {
	return &CategoryService_Get_Call{Call: _e.mock.On("Get", ctx, categoryCode)}
}

func (_c *CategoryService_Get_Call) Run(run func(ctx context.Context, categoryCode string)) *CategoryService_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *CategoryService_Get_Call) Return(_a0 domain.Category, _a1 error) *CategoryService_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CategoryService_Get_Call) RunAndReturn(run func(context.Context, string) (domain.Category, error)) *CategoryService_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Tree provides a mock function with given fields: ctx, activeCategoryCode
func (_m *CategoryService) Tree(ctx context.Context, activeCategoryCode string) (domain.Tree, error) {
	ret := _m.Called(ctx, activeCategoryCode)

	if len(ret) == 0 {
		panic("no return value specified for Tree")
	}

	var r0 domain.Tree
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.Tree, error)); ok {
		return rf(ctx, activeCategoryCode)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Tree); ok {
		r0 = rf(ctx, activeCategoryCode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.Tree)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, activeCategoryCode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CategoryService_Tree_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Tree'
type CategoryService_Tree_Call struct {
	*mock.Call
}

// Tree is a helper method to define mock.On call
//   - ctx context.Context
//   - activeCategoryCode string
func (_e *CategoryService_Expecter) Tree(ctx interface{}, activeCategoryCode interface{}) *CategoryService_Tree_Call {
	return &CategoryService_Tree_Call{Call: _e.mock.On("Tree", ctx, activeCategoryCode)}
}

func (_c *CategoryService_Tree_Call) Run(run func(ctx context.Context, activeCategoryCode string)) *CategoryService_Tree_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *CategoryService_Tree_Call) Return(_a0 domain.Tree, _a1 error) *CategoryService_Tree_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CategoryService_Tree_Call) RunAndReturn(run func(context.Context, string) (domain.Tree, error)) *CategoryService_Tree_Call {
	_c.Call.Return(run)
	return _c
}

// NewCategoryService creates a new instance of CategoryService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCategoryService(t interface {
	mock.TestingT
	Cleanup(func())
}) *CategoryService {
	mock := &CategoryService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
