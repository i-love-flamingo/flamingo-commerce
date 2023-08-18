package controller

import (
	"context"
	"net/url"
	"strconv"

	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/category/domain"
	productApplication "flamingo.me/flamingo-commerce/v3/product/application"

	searchApplication "flamingo.me/flamingo-commerce/v3/search/application"
	searchDomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/utils"
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.4 --name ProductSearchService --case snake

type (
	// QueryHandler provides the base command logic that is agnostic to the actual view type
	QueryHandler interface {
		Execute(ctx context.Context, req Request) (*Result, *RedirectResult, error)
	}

	// ProductSearchService interface that describes the expected dependency. (Is fulfilled by the product package)
	ProductSearchService interface {
		Find(ctx context.Context, searchRequest *searchApplication.SearchRequest) (*productApplication.SearchResult, error)
	}

	// QueryHandlerImpl is the default implementation of QueryHandler
	QueryHandlerImpl struct {
		categoryService      domain.CategoryService
		productSearchService ProductSearchService
	}

	// Request is a request for a category view
	Request struct {
		Code     string
		Name     string
		QueryAll url.Values
		URL      url.URL
	}

	// Result for found category
	Result struct {
		ProductSearchResult *productApplication.SearchResult
		Category            domain.Category
		CategoryTree        domain.Tree
		SearchMeta          searchDomain.SearchMeta
		PaginationInfo      utils.PaginationInfo
	}

	// RedirectResult signals that a request for a category view shall be redirected using the defined parameters
	RedirectResult struct {
		Code string
		Name string
	}
)

var _ QueryHandler = (*QueryHandlerImpl)(nil)

// Inject injects dependencies
func (c *QueryHandlerImpl) Inject(
	categoryService domain.CategoryService,
	searchService ProductSearchService,
) {
	c.categoryService = categoryService
	c.productSearchService = searchService
}

// Execute Action to display a category page for any view
// error might be domain.ErrNotFound to indicate that the category was not found
func (c *QueryHandlerImpl) Execute(ctx context.Context, req Request) (*Result, *RedirectResult, error) {
	treeRoot, err := c.categoryService.Tree(ctx, req.Code)
	if err != nil {
		return nil, nil, err
	}

	currentCategory, err := c.categoryService.Get(ctx, req.Code)
	if err != nil {
		return nil, nil, err
	}

	// Normalize url if required:
	expectedName := web.URLTitle(currentCategory.Name())
	if expectedName != req.Name {
		return nil, &RedirectResult{Code: currentCategory.Code(), Name: expectedName}, nil
	}

	searchRequest := &searchApplication.SearchRequest{}
	for k, v := range req.QueryAll {
		switch k {
		case "page":
			page, _ := strconv.ParseInt(v[0], 10, 64)
			searchRequest.SetAdditionalFilter(searchDomain.NewPaginationPageFilter(int(page)))

		default:
			searchRequest.SetAdditionalFilter(searchDomain.NewKeyValueFilter(k, v))
		}
	}
	searchRequest.SetAdditionalFilter(domain.NewCategoryFacet(currentCategory.Code()))

	products, err := c.productSearchService.Find(ctx, searchRequest)
	if err != nil {
		return nil, nil, err
	}

	return &Result{
		Category:            currentCategory,
		CategoryTree:        treeRoot,
		ProductSearchResult: products,
		SearchMeta:          products.SearchMeta,
		PaginationInfo:      products.PaginationInfo,
	}, nil, nil
}
