package controller_test

import (
	"context"
	"errors"

	"flamingo.me/flamingo-commerce/v3/category/interfaces/controller"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"flamingo.me/flamingo-commerce/v3/category/domain"
	productApplication "flamingo.me/flamingo-commerce/v3/product/application"
	searchApplication "flamingo.me/flamingo-commerce/v3/search/application"
	searchDomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/utils"
)

type (
	mockCategoryService struct {
		mock.Mock
	}
	mockProductSearchService struct {
		mockFunc mockSearchServiceFindFunc
	}
	mockSearchServiceFindFunc func(ctx context.Context, searchRequest *searchApplication.SearchRequest) (*productApplication.SearchResult, error)
)

func (m *mockCategoryService) Tree(ctx context.Context, activeCategoryCode string) (domain.Tree, error) {
	called := m.Called(ctx, activeCategoryCode)
	return called.Get(0).(domain.Tree), called.Error(1)
}

func (m *mockCategoryService) Get(ctx context.Context, categoryCode string) (domain.Category, error) {
	called := m.Called(ctx, categoryCode)
	return called.Get(0).(domain.Category), called.Error(1)
}

func (m *mockProductSearchService) Find(ctx context.Context, searchRequest *searchApplication.SearchRequest) (*productApplication.SearchResult, error) {
	return m.mockFunc(ctx, searchRequest)
}

func TestDefaultCommandHandler_Execute_SearchService(t *testing.T) {
	type args struct {
		categoryService   domain.CategoryService
		searchServiceFind mockSearchServiceFindFunc
	}

	categoryDataFixture := &domain.CategoryData{CategoryName: "test", CategoryCode: "test"}
	catService := mockCategoryService{}
	catService.On("Tree", mock.Anything, mock.Anything).Return(&domain.TreeData{}, nil)
	catService.On("Get", mock.Anything, mock.Anything).Return(categoryDataFixture, nil)

	tests := []struct {
		name             string
		args             args
		request          controller.Request
		wantViewData     *controller.Result
		wantViewRedirect *controller.RedirectResult
		wantViewError    error
	}{
		{
			name: "successful product search results in success response",
			args: args{
				categoryService: &catService,
				searchServiceFind: func(ctx context.Context, searchRequest *searchApplication.SearchRequest) (*productApplication.SearchResult, error) {
					return &productApplication.SearchResult{}, nil
				},
			},
			request: controller.Request{
				Code: "test",
				Name: "test",
			},
			wantViewData: &controller.Result{
				ProductSearchResult: &productApplication.SearchResult{},
				Category:            categoryDataFixture,
				CategoryTree:        &domain.TreeData{},
				SearchMeta:          searchDomain.SearchMeta{},
				PaginationInfo:      utils.PaginationInfo{},
			},
		},
		{
			name: "not found error results in not found error response",
			args: args{
				categoryService: &catService,
				searchServiceFind: func(ctx context.Context, searchRequest *searchApplication.SearchRequest) (*productApplication.SearchResult, error) {
					return nil, searchDomain.ErrNotFound
				},
			},
			request: controller.Request{
				Code: "test",
				Name: "test",
			},
			wantViewError: searchDomain.ErrNotFound,
		},
		{
			name: "error results in internal server error response",
			args: args{
				categoryService: &catService,
				searchServiceFind: func(ctx context.Context, searchRequest *searchApplication.SearchRequest) (*productApplication.SearchResult, error) {
					return nil, errors.New("other")
				},
			},
			request: controller.Request{
				Code: "test",
				Name: "test",
			},
			wantViewError: errors.New("other"),
		},
		{
			name: "redirect if name is wrong",
			args: args{
				categoryService: &catService,
				searchServiceFind: func(ctx context.Context, searchRequest *searchApplication.SearchRequest) (*productApplication.SearchResult, error) {
					return &productApplication.SearchResult{}, nil
				},
			},
			request: controller.Request{
				Code: "test",
				Name: "testt",
			},
			wantViewRedirect: &controller.RedirectResult{Code: "test", Name: "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			commandHandler := controller.QueryHandlerImpl{}
			commandHandler.Inject(tt.args.categoryService, &mockProductSearchService{
				mockFunc: tt.args.searchServiceFind,
			})

			gotViewData, gotRedirect, gotError := commandHandler.Execute(context.Background(), tt.request)

			a := assert.New(t)

			a.Equal(tt.wantViewError, gotError)
			a.Equal(tt.wantViewRedirect, gotRedirect)
			a.Equal(tt.wantViewData, gotViewData)
		})
	}
}
