package controller_test

import (
	"context"
	"errors"
	"testing"

	"flamingo.me/flamingo-commerce/v3/category/domain/mocks"
	"flamingo.me/flamingo-commerce/v3/category/interfaces/controller"
	categoryMocks "flamingo.me/flamingo-commerce/v3/category/interfaces/controller/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"flamingo.me/flamingo-commerce/v3/category/domain"
	productApplication "flamingo.me/flamingo-commerce/v3/product/application"
	searchDomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/utils"
)

func TestDefaultCommandHandler_Execute_SearchService(t *testing.T) {
	type args struct {
		categoryService     domain.CategoryService
		searchServiceResult *productApplication.SearchResult
		searchServiceError  error
	}

	categoryDataFixture := &domain.CategoryData{CategoryName: "test", CategoryCode: "test"}
	catService := mocks.NewCategoryService(t)
	catService.EXPECT().Tree(mock.Anything, mock.Anything).Return(&domain.TreeData{}, nil)
	catService.EXPECT().Get(mock.Anything, mock.Anything).Return(categoryDataFixture, nil)

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
				categoryService:     catService,
				searchServiceResult: &productApplication.SearchResult{},
				searchServiceError:  nil,
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
				categoryService:     catService,
				searchServiceResult: nil,
				searchServiceError:  searchDomain.ErrNotFound,
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
				categoryService:     catService,
				searchServiceResult: nil,
				searchServiceError:  errors.New("other"),
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
				categoryService:     catService,
				searchServiceResult: nil,
				searchServiceError:  nil,
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

			productSearchService := new(categoryMocks.ProductSearchService)
			productSearchService.EXPECT().Find(mock.Anything, mock.Anything).Return(tt.args.searchServiceResult, tt.args.searchServiceError)

			commandHandler := controller.QueryHandlerImpl{}
			commandHandler.Inject(tt.args.categoryService, productSearchService)

			gotViewData, gotRedirect, gotError := commandHandler.Execute(context.Background(), tt.request)

			a := assert.New(t)

			a.Equal(tt.wantViewError, gotError)
			a.Equal(tt.wantViewRedirect, gotRedirect)
			a.Equal(tt.wantViewData, gotViewData)
		})
	}
}
