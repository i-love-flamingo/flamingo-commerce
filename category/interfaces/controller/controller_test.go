package controller_test

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/flamingo-commerce/v3/category/interfaces/controller"
	"flamingo.me/flamingo-commerce/v3/product/application"
	searchApplication "flamingo.me/flamingo-commerce/v3/search/application"
	searchdomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/utils"
)

type mockCategoryService struct {
	mock.Mock
}

func (m *mockCategoryService) Tree(ctx context.Context, activeCategoryCode string) (domain.Tree, error) {
	called := m.Called(ctx, activeCategoryCode)
	return called.Get(0).(domain.Tree), called.Error(1)
}

func (m *mockCategoryService) Get(ctx context.Context, categoryCode string) (domain.Category, error) {
	called := m.Called(ctx, categoryCode)
	return called.Get(0).(domain.Category), called.Error(1)
}

func TestDefaultBaseViewController_Get_SearchService(t *testing.T) {
	type args struct {
		categoryService                domain.CategoryService
		searchServiceFind              controller.SearchServiceFindFunc
		breadcrumbServiceAddBreadcrumb controller.BreadcrumbServiceAddBreadcrumbFunc
		paginationInfoFactoryBuild     controller.PaginationInfoFactoryBuildFunc
	}

	catService := mockCategoryService{}
	catService.On("Tree", mock.Anything, mock.Anything).Return(&domain.TreeData{}, nil)
	catService.On("Get", mock.Anything, mock.Anything).Return(&domain.CategoryData{CategoryName: "test"}, nil)

	emptyBreadcrumbServiceAddBreadcrumb := func(ctx context.Context, tree domain.Tree) {}
	emptyPaginationInfoFactoryBuild := func(activePage int, totalHits int, pageSize int, lastPage int, urlBase url.URL) utils.PaginationInfo {
		return utils.PaginationInfo{}
	}

	tests := []struct {
		name             string
		args             args
		request          controller.ViewRequest
		wantViewData     *controller.ViewData
		wantViewRedirect *controller.ViewRedirect
		wantViewError    *controller.ViewError
	}{
		{
			name: "successful product search results in success response",
			args: args{
				categoryService: &catService,
				searchServiceFind: func(ctx context.Context, searchRequest *searchApplication.SearchRequest) (*application.SearchResult, error) {
					return &application.SearchResult{}, nil
				},
				breadcrumbServiceAddBreadcrumb: emptyBreadcrumbServiceAddBreadcrumb,
				paginationInfoFactoryBuild:     emptyPaginationInfoFactoryBuild,
			},
			request: controller.ViewRequest{
				Code: "test",
				Name: "test",
			},
			wantViewData: &controller.ViewData{
				ProductSearchResult: &application.SearchResult{},
				Category:            &domain.CategoryData{CategoryName: "test"},
				CategoryTree:        &domain.TreeData{},
				SearchMeta:          searchdomain.SearchMeta{},
				PaginationInfo:      utils.PaginationInfo{},
			},
		},
		{
			name: "not found error results in not found error response",
			args: args{
				categoryService: &catService,
				searchServiceFind: func(ctx context.Context, searchRequest *searchApplication.SearchRequest) (*application.SearchResult, error) {
					return nil, searchdomain.ErrNotFound
				},
				breadcrumbServiceAddBreadcrumb: emptyBreadcrumbServiceAddBreadcrumb,
				paginationInfoFactoryBuild:     emptyPaginationInfoFactoryBuild,
			},
			request: controller.ViewRequest{
				Code: "test",
				Name: "test",
			},
			wantViewError: &controller.ViewError{
				NotFound: searchdomain.ErrNotFound,
			},
		},
		{
			name: "error results in internal server error response",
			args: args{
				categoryService: &catService,
				searchServiceFind: func(ctx context.Context, searchRequest *searchApplication.SearchRequest) (*application.SearchResult, error) {
					return nil, errors.New("other")
				},
				breadcrumbServiceAddBreadcrumb: emptyBreadcrumbServiceAddBreadcrumb,
				paginationInfoFactoryBuild:     emptyPaginationInfoFactoryBuild,
			},
			request: controller.ViewRequest{
				Code: "test",
				Name: "test",
			},
			wantViewError: &controller.ViewError{
				Other: errors.New("other"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := controller.NewDefaultBaseViewController(
				tt.args.categoryService,
				tt.args.searchServiceFind,
				tt.args.breadcrumbServiceAddBreadcrumb,
				tt.args.paginationInfoFactoryBuild,
			)

			gotViewData, gotRedirect, gotError := vc.Get(context.Background(), tt.request)

			a := assert.New(t)

			a.Equal(tt.wantViewError, gotError)
			a.Equal(tt.wantViewRedirect, gotRedirect)
			a.Equal(tt.wantViewData, gotViewData)
		})
	}
}
