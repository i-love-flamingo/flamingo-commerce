package application

import (
	"context"
	"net/url"
	"strconv"

	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/flamingo-commerce/v3/product/application"
	searchApplication "flamingo.me/flamingo-commerce/v3/search/application"
	searchDomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/utils"
)

type (
	// CommandHandler provides the base command logic that is agnostic to the actual view type
	CommandHandler interface {
		Execute(ctx context.Context, req CategoryRequest) (*CommandResult, *CommandRedirect, *CommandError)
	}

	// SearchServiceFindFunc is a simple abstraction to application.ProductSearchService.Find
	SearchServiceFindFunc func(ctx context.Context, searchRequest *searchApplication.SearchRequest) (*application.SearchResult, error)
	// BreadcrumbServiceAddBreadcrumbFunc is a simple abstraction to breadcrumb.BreadcrumbService.AddBreadcrumb
	BreadcrumbServiceAddBreadcrumbFunc func(ctx context.Context, tree domain.Tree)
	// PaginationInfoFactoryBuildFunc is a simple abstraction to utils.PaginationInfoFactory.Build
	PaginationInfoFactoryBuildFunc func(activePage int, totalHits int, pageSize int, lastPage int, urlBase url.URL) utils.PaginationInfo

	// DefaultCommandHandler is the default implementation of CommandHandler
	DefaultCommandHandler struct {
		categoryService                domain.CategoryService
		searchServiceFind              SearchServiceFindFunc
		breadcrumbServiceAddBreadcrumb BreadcrumbServiceAddBreadcrumbFunc
		paginationInfoFactoryBuild     PaginationInfoFactoryBuildFunc
	}

	// CategoryRequest is a request for a category view
	CategoryRequest struct {
		Code     string
		Name     string
		QueryAll url.Values
		URL      url.URL
	}

	// CommandResult for rendering context
	CommandResult struct {
		ProductSearchResult *application.SearchResult
		Category            domain.Category
		CategoryTree        domain.Tree
		SearchMeta          searchDomain.SearchMeta
		PaginationInfo      utils.PaginationInfo
	}

	// CommandRedirect signals that a request for a category view shall be redirected using the defined parameters
	CommandRedirect struct {
		Code string
		Name string
	}

	// CommandError signals a category view error
	CommandError struct {
		NotFound error
		Other    error
	}
)

var _ CommandHandler = (*DefaultCommandHandler)(nil)

// NewDefaultCommandHandler creates a new DefaultCommandHandler
func NewDefaultCommandHandler(
	categoryService domain.CategoryService,
	searchServiceFind SearchServiceFindFunc,
	breadcrumbServiceAddBreadcrumb BreadcrumbServiceAddBreadcrumbFunc,
	paginationInfoFactoryBuild PaginationInfoFactoryBuildFunc,
) *DefaultCommandHandler {
	return &DefaultCommandHandler{
		categoryService:                categoryService,
		searchServiceFind:              searchServiceFind,
		breadcrumbServiceAddBreadcrumb: breadcrumbServiceAddBreadcrumb,
		paginationInfoFactoryBuild:     paginationInfoFactoryBuild,
	}
}

// Inject injects dependencies
func (c *DefaultCommandHandler) Inject(
	categoryService domain.CategoryService,
	searchService *application.ProductSearchService,
	paginationInfoFactory *utils.PaginationInfoFactory,
	breadcrumbService *BreadcrumbService,
) {
	c.categoryService = categoryService

	c.searchServiceFind = func(ctx context.Context, searchRequest *searchApplication.SearchRequest) (result *application.SearchResult, e error) {
		return searchService.Find(ctx, searchRequest)
	}

	c.paginationInfoFactoryBuild = func(activePage int, totalHits int, pageSize int, lastPage int, urlBase url.URL) utils.PaginationInfo {
		return paginationInfoFactory.Build(activePage, totalHits, pageSize, lastPage, &urlBase)
	}

	c.breadcrumbServiceAddBreadcrumb = func(ctx context.Context, tree domain.Tree) {
		breadcrumbService.AddBreadcrumb(ctx, tree)
	}
}

// Execute Action to display a category page for any view
func (c *DefaultCommandHandler) Execute(ctx context.Context, req CategoryRequest) (*CommandResult, *CommandRedirect, *CommandError) {
	treeRoot, err := c.categoryService.Tree(ctx, req.Code)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, nil, &CommandError{NotFound: err}
		}
		return nil, nil, &CommandError{Other: err}
	}

	currentCategory, err := c.categoryService.Get(ctx, req.Code)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, nil, &CommandError{NotFound: err}
		}
		return nil, nil, &CommandError{Other: err}
	}

	// Normalize url if required:
	expectedName := web.URLTitle(currentCategory.Name())
	if expectedName != req.Name {
		return nil, &CommandRedirect{Code: currentCategory.Code(), Name: expectedName}, nil
	}

	searchRequest := &searchApplication.SearchRequest{}
	for k, v := range req.QueryAll {
		switch k {
		case "page":
			page, _ := strconv.ParseInt(v[0], 10, 64)
			searchRequest.SetAdditionalFilter(searchDomain.NewPaginationPageFilter(int(page)))
			break
		default:
			searchRequest.SetAdditionalFilter(searchDomain.NewKeyValueFilter(k, v))
		}
	}
	searchRequest.SetAdditionalFilter(domain.NewCategoryFacet(currentCategory.Code()))

	products, err := c.searchServiceFind(ctx, searchRequest)
	if err != nil {
		if err == searchDomain.ErrNotFound {
			return nil, nil, &CommandError{NotFound: err}
		}
		return nil, nil, &CommandError{Other: err}
	}

	c.breadcrumbServiceAddBreadcrumb(ctx, treeRoot)

	paginationInfo := c.paginationInfoFactoryBuild(products.SearchMeta.Page, products.SearchMeta.NumResults, 30, products.SearchMeta.NumPages, req.URL)

	return &CommandResult{
		Category:            currentCategory,
		CategoryTree:        treeRoot,
		ProductSearchResult: products,
		SearchMeta:          products.SearchMeta,
		PaginationInfo:      paginationInfo,
	}, nil, nil
}
