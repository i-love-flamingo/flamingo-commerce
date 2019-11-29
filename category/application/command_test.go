package application_test

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"flamingo.me/flamingo-commerce/v3/category/application"
	"flamingo.me/flamingo-commerce/v3/category/domain"
	productApplication "flamingo.me/flamingo-commerce/v3/product/application"
	searchApplication "flamingo.me/flamingo-commerce/v3/search/application"
	searchDomain "flamingo.me/flamingo-commerce/v3/search/domain"
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

func TestDefaultCommandHandler_Execute_SearchService(t *testing.T) {
	type args struct {
		categoryService                domain.CategoryService
		searchServiceFind              application.SearchServiceFindFunc
		breadcrumbServiceAddBreadcrumb application.BreadcrumbServiceAddBreadcrumbFunc
		paginationInfoFactoryBuild     application.PaginationInfoFactoryBuildFunc
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
		request          application.CategoryRequest
		wantViewData     *application.CommandResult
		wantViewRedirect *application.CommandRedirect
		wantViewError    *application.CommandError
	}{
		{
			name: "successful product search results in success response",
			args: args{
				categoryService: &catService,
				searchServiceFind: func(ctx context.Context, searchRequest *searchApplication.SearchRequest) (*productApplication.SearchResult, error) {
					return &productApplication.SearchResult{}, nil
				},
				breadcrumbServiceAddBreadcrumb: emptyBreadcrumbServiceAddBreadcrumb,
				paginationInfoFactoryBuild:     emptyPaginationInfoFactoryBuild,
			},
			request: application.CategoryRequest{
				Code: "test",
				Name: "test",
			},
			wantViewData: &application.CommandResult{
				ProductSearchResult: &productApplication.SearchResult{},
				Category:            &domain.CategoryData{CategoryName: "test"},
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
				breadcrumbServiceAddBreadcrumb: emptyBreadcrumbServiceAddBreadcrumb,
				paginationInfoFactoryBuild:     emptyPaginationInfoFactoryBuild,
			},
			request: application.CategoryRequest{
				Code: "test",
				Name: "test",
			},
			wantViewError: &application.CommandError{
				NotFound: searchDomain.ErrNotFound,
			},
		},
		{
			name: "error results in internal server error response",
			args: args{
				categoryService: &catService,
				searchServiceFind: func(ctx context.Context, searchRequest *searchApplication.SearchRequest) (*productApplication.SearchResult, error) {
					return nil, errors.New("other")
				},
				breadcrumbServiceAddBreadcrumb: emptyBreadcrumbServiceAddBreadcrumb,
				paginationInfoFactoryBuild:     emptyPaginationInfoFactoryBuild,
			},
			request: application.CategoryRequest{
				Code: "test",
				Name: "test",
			},
			wantViewError: &application.CommandError{
				Other: errors.New("other"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := application.NewDefaultCommandHandler(
				tt.args.categoryService,
				tt.args.searchServiceFind,
				tt.args.breadcrumbServiceAddBreadcrumb,
				tt.args.paginationInfoFactoryBuild,
			)

			gotViewData, gotRedirect, gotError := vc.Execute(context.Background(), tt.request)

			a := assert.New(t)

			a.Equal(tt.wantViewError, gotError)
			a.Equal(tt.wantViewRedirect, gotRedirect)
			a.Equal(tt.wantViewData, gotViewData)
		})
	}
}
