package controller

import (
	"context"
	"net/url"
	"strconv"

	breadcrumb "flamingo.me/flamingo-commerce/v3/category/application"

	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/flamingo-commerce/v3/product/application"
	searchApplication "flamingo.me/flamingo-commerce/v3/search/application"
	searchdomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/utils"
)

type (
	// ViewController provides web-specific actions for category single view
	ViewController struct {
		base           BaseViewController
		responder      *web.Responder
		router         *web.Router
		template       string
		teaserTemplate string
	}

	// BaseViewController provides the base controller logic that is agnostic to the actual view type
	BaseViewController interface {
		Get(ctx context.Context, req ViewRequest) (*ViewData, *ViewRedirect, *ViewError)
	}

	// SearchServiceFindFunc is a simple abstraction to application.ProductSearchService.Find
	SearchServiceFindFunc func(ctx context.Context, searchRequest *searchApplication.SearchRequest) (*application.SearchResult, error)
	// BreadcrumbServiceAddBreadcrumbFunc is a simple abstraction to breadcrumb.BreadcrumbService.AddBreadcrumb
	BreadcrumbServiceAddBreadcrumbFunc func(ctx context.Context, tree domain.Tree)
	// PaginationInfoFactoryBuildFunc is a simple abstraction to utils.PaginationInfoFactory.Build
	PaginationInfoFactoryBuildFunc func(activePage int, totalHits int, pageSize int, lastPage int, urlBase url.URL) utils.PaginationInfo

	// DefaultBaseViewController is the default implementation of BaseViewController
	DefaultBaseViewController struct {
		categoryService                domain.CategoryService
		searchServiceFind              SearchServiceFindFunc
		breadcrumbServiceAddBreadcrumb BreadcrumbServiceAddBreadcrumbFunc
		paginationInfoFactoryBuild     PaginationInfoFactoryBuildFunc
	}

	// ViewRequest is a request for a category view
	ViewRequest struct {
		Code     string
		Name     string
		QueryAll url.Values
		URL      url.URL
	}

	// ViewData for rendering context
	ViewData struct {
		ProductSearchResult *application.SearchResult
		Category            domain.Category
		CategoryTree        domain.Tree
		SearchMeta          searchdomain.SearchMeta
		PaginationInfo      utils.PaginationInfo
	}

	// ViewRedirect signals that a request for a category view shall be redirected using the defined parameters
	ViewRedirect struct {
		Code string
		Name string
	}

	// ViewError signals a category view error
	ViewError struct {
		NotFound error
		Other    error
	}
)

// NewDefaultBaseViewController creates a new DefaultBaseViewController
func NewDefaultBaseViewController(
	categoryService domain.CategoryService,
	searchServiceFind SearchServiceFindFunc,
	breadcrumbServiceAddBreadcrumb BreadcrumbServiceAddBreadcrumbFunc,
	paginationInfoFactoryBuild PaginationInfoFactoryBuildFunc,
) *DefaultBaseViewController {
	return &DefaultBaseViewController{
		categoryService:                categoryService,
		searchServiceFind:              searchServiceFind,
		breadcrumbServiceAddBreadcrumb: breadcrumbServiceAddBreadcrumb,
		paginationInfoFactoryBuild:     paginationInfoFactoryBuild,
	}
}

// Inject injects dependencies
func (c *DefaultBaseViewController) Inject(
	categoryService domain.CategoryService,
	searchService *application.ProductSearchService,
	paginationInfoFactory *utils.PaginationInfoFactory,
	breadcrumbService *breadcrumb.BreadcrumbService,
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

// Get Action to display a category page for any view
func (c *DefaultBaseViewController) Get(ctx context.Context, req ViewRequest) (*ViewData, *ViewRedirect, *ViewError) {
	treeRoot, err := c.categoryService.Tree(ctx, req.Code)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, nil, &ViewError{NotFound: err}
		}
		return nil, nil, &ViewError{Other: err}
	}

	currentCategory, err := c.categoryService.Get(ctx, req.Code)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, nil, &ViewError{NotFound: err}
		}
		return nil, nil, &ViewError{Other: err}
	}

	// Normalize url if required:
	expectedName := web.URLTitle(currentCategory.Name())
	if expectedName != req.Name {
		return nil, &ViewRedirect{Code: currentCategory.Code(), Name: expectedName}, nil
	}

	searchRequest := &searchApplication.SearchRequest{}
	for k, v := range req.QueryAll {
		switch k {
		case "page":
			page, _ := strconv.ParseInt(v[0], 10, 64)
			searchRequest.SetAdditionalFilter(searchdomain.NewPaginationPageFilter(int(page)))
			break
		default:
			searchRequest.SetAdditionalFilter(searchdomain.NewKeyValueFilter(k, v))
		}
	}
	searchRequest.SetAdditionalFilter(domain.NewCategoryFacet(currentCategory.Code()))

	products, err := c.searchServiceFind(ctx, searchRequest)
	if err != nil {
		if err == searchdomain.ErrNotFound {
			return nil, nil, &ViewError{NotFound: err}
		}
		return nil, nil, &ViewError{Other: err}
	}

	c.breadcrumbServiceAddBreadcrumb(ctx, treeRoot)

	paginationInfo := c.paginationInfoFactoryBuild(products.SearchMeta.Page, products.SearchMeta.NumResults, 30, products.SearchMeta.NumPages, req.URL)

	return &ViewData{
		Category:            currentCategory,
		CategoryTree:        treeRoot,
		ProductSearchResult: products,
		SearchMeta:          products.SearchMeta,
		PaginationInfo:      paginationInfo,
	}, nil, nil
}

// Inject the ViewController controller required dependencies
func (vc *ViewController) Inject(
	base BaseViewController,
	responder *web.Responder,
	router *web.Router,
	config *struct {
		Template       string `inject:"config:commerce.category.view.template"`
		TeaserTemplate string `inject:"config:commerce.category.view.teaserTemplate"`
	},
) *ViewController {
	vc.base = base
	vc.responder = responder
	vc.router = router

	if config != nil {
		vc.template = config.Template
		vc.teaserTemplate = config.TeaserTemplate
	}

	return vc
}

// Get Action to display a category page for the web
func (vc *ViewController) Get(c context.Context, request *web.Request) web.Result {

	viewData, viewRedirect, viewError := vc.base.Get(c, ViewRequest{
		Code:     request.Params["code"],
		Name:     request.Params["name"],
		URL:      *request.Request().URL,
		QueryAll: request.QueryAll(),
	})

	if viewError != nil {
		if viewError.NotFound != nil {
			return vc.responder.NotFound(viewError.NotFound)
		} else {
			return vc.responder.ServerError(viewError.Other)
		}
	}

	if viewRedirect != nil {
		redirectParams := map[string]string{
			"code": viewRedirect.Code,
			"name": viewRedirect.Name,
		}

		u, _ := vc.router.Relative("category.view", redirectParams)
		u.RawQuery = request.QueryAll().Encode()
		return vc.responder.URLRedirect(u).Permanent()
	}

	var template string
	switch viewData.Category.CategoryType() {
	case domain.TypeTeaser:
		template = vc.teaserTemplate
	default:
		template = vc.template
	}

	return vc.responder.Render(template, viewData)
}
