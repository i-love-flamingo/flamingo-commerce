package controller

import (
	"context"
	"errors"
	breadcrumb "flamingo.me/flamingo-commerce/category/application"
	"net/url"

	"flamingo.me/flamingo-commerce/category/domain"
	"flamingo.me/flamingo-commerce/product/application"
	searchApplication "flamingo.me/flamingo-commerce/search/application"
	searchdomain "flamingo.me/flamingo-commerce/search/domain"
	"flamingo.me/flamingo-commerce/search/utils"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
)

type (
	// View demonstrates a product view controller
	View struct {
		responder.ErrorAware
		responder.RenderAware
		responder.RedirectAware
		domain.CategoryService
		SearchService         *application.ProductSearchService
		router                *router.Router
		template              string
		teaserTemplate        string
		logger                flamingo.Logger
		paginationInfoFactory *utils.PaginationInfoFactory
		breadcrumbService     *breadcrumb.BreadcrumbService
	}

	// ViewData for rendering context
	ViewData struct {
		ProductSearchResult *application.SearchResult
		Category            domain.Category
		CategoryTree        domain.Category
		SearchMeta          searchdomain.SearchMeta
		PaginationInfo      utils.PaginationInfo
	}
)

// Inject the View controller required dependencies
func (vc *View) Inject(
	errorAware responder.ErrorAware,
	renderAware responder.RenderAware,
	redirectAware responder.RedirectAware,
	categoryService domain.CategoryService,
	searchService *application.ProductSearchService,
	router *router.Router,
	logger flamingo.Logger,
	paginationInfoFactory *utils.PaginationInfoFactory,
	breadcrumbService *breadcrumb.BreadcrumbService,

	config *struct {
		Template       string `inject:"config:category.view.template"`
		TeaserTemplate string `inject:"config:category.view.teaserTemplate"`
	},
) {
	vc.ErrorAware = errorAware
	vc.RenderAware = renderAware
	vc.RedirectAware = redirectAware
	vc.CategoryService = categoryService
	vc.SearchService = searchService
	vc.router = router
	vc.logger = logger
	vc.paginationInfoFactory = paginationInfoFactory
	vc.template = config.Template
	vc.teaserTemplate = config.TeaserTemplate
	vc.breadcrumbService = breadcrumbService
}

// Get Response for Product matching sku param
func (vc *View) Get(c context.Context, request *web.Request) web.Response {
	categoryRoot, err := vc.CategoryService.Tree(c, request.MustParam1("code"))
	if err == domain.ErrNotFound {
		return vc.ErrorNotFound(c, err)
	} else if err != nil {
		return vc.Error(c, err)
	}

	category := domain.GetActive(categoryRoot)
	if category == nil {
		return vc.ErrorNotFound(c, errors.New("Active Category not found"))
	}

	expectedName := web.URLTitle(category.Name())
	if name, _ := request.Param1("name"); expectedName != name {

		redirectParams := router.P{
			"code": category.Code(),
			"name": expectedName,
		}

		u := vc.router.URL("category.view", redirectParams)
		u.RawQuery = url.Values(request.QueryAll()).Encode()
		return vc.RedirectPermanentURL(u.String())
	}

	queryAll := request.QueryAll()
	filter := make(map[string]interface{}, len(queryAll)+1)
	for k, v := range queryAll {
		filter[k] = v
	}
	filter[string(domain.CategoryKey)] = domain.NewCategoryFacet(category)

	searchRequest := &searchApplication.SearchRequest{
		FilterBy: filter,
	}

	products, err := vc.SearchService.Find(c, searchRequest)
	if err != nil {
		return vc.Error(c, err)
	}

	vc.breadcrumbService.AddBreadcrumb(c, categoryRoot)

	paginationInfo := vc.paginationInfoFactory.Build(products.SearchMeta.Page, products.SearchMeta.NumResults, 30, products.SearchMeta.NumPages, request.Request().URL)

	var template string
	switch category.CategoryType() {
	case domain.TypeTeaser:
		template = vc.teaserTemplate
	default:
		template = vc.template
	}

	return vc.Render(c, template, ViewData{
		Category:            category,
		CategoryTree:        categoryRoot,
		ProductSearchResult: products,
		SearchMeta:          products.SearchMeta,
		PaginationInfo:      paginationInfo,
	})
}
